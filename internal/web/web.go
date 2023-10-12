package web

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/holedaemon/bot2/internal/pkg/pgstore"
	"github.com/holedaemon/bot2/internal/web/templates"
	"github.com/jellydator/ttlcache/v3"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

//go:embed static
var static embed.FS

var assetsDir fs.FS

func init() {
	var err error
	assetsDir, err = fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
}

// Server is the HTTP server responsible for serving BOT/2's website.
type Server struct {
	Debug  bool
	Addr   string
	DB     *sql.DB
	OAuth2 *oauth2.Config
	Admins map[string]string

	sessionManager *scs.SessionManager

	stateCache *ttlcache.Cache[string, bool]
	userCache  *ttlcache.Cache[string, []string]
}

// New creates a new Server.
func New(opts ...Option) (*Server, error) {
	srv := &Server{
		stateCache: ttlcache.New[string, bool](),
		userCache: ttlcache.New[string, []string](
			ttlcache.WithTTL[string, []string](5 * time.Minute),
		),
	}

	for _, o := range opts {
		o(srv)
	}

	if srv.DB == nil {
		return nil, fmt.Errorf("web: missing db")
	}

	if srv.OAuth2 == nil {
		return nil, fmt.Errorf("web: missing oauth2 config")
	}

	if srv.Admins == nil {
		srv.Admins = make(map[string]string)
	}

	sm := scs.New()
	sm.Cookie.Name = sessionName
	sm.Cookie.Secure = !srv.Debug
	srv.sessionManager = sm

	return srv, nil
}

// Run starts a Server.
func (s *Server) Run(ctx context.Context) error {
	r := chi.NewMux()

	r.Use(s.recoverer)

	logger := ctxlog.FromContext(ctx)
	r.Use(requestLogger(logger))

	r.Get("/", s.index)
	r.Get("/about", s.about)
	r.Get("/login", s.authDiscord)
	r.Get("/logout", s.logout)
	r.Get("/auth/discord", s.authDiscord)
	r.Get("/auth/discord/callback", s.authDiscordCallback)
	r.Get("/guilds", s.guilds)
	r.Get("/docs", s.docs)

	r.Route("/admin", s.routeAdmin)

	r.Route("/guild/{id}", func(r chi.Router) {
		r.Use(s.guildCheck)

		r.Get("/", s.guild)
		r.Get("/quotes", s.guildQuotes)
		r.Get("/roles", s.guildRoles)
		r.Get("/tags", s.guildTags)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.errorPage(w, r, http.StatusNotFound, "Whatever you're looking for ain't here")
	})

	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(assetsDir))))
	r.Handle("/favicon.ico", http.RedirectHandler("/static/favicon.ico", http.StatusFound))

	go s.stateCache.Start()
	defer s.stateCache.Stop()

	store := pgstore.New(s.DB)
	store.Start(ctx)
	s.sessionManager.Store = store

	srv := &http.Server{
		Addr:        s.Addr,
		Handler:     s.sessionManager.LoadAndSave(r),
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			ctxlog.Error(ctx, "error shutting down server", zap.Error(err))
			return
		}
	}()

	return srv.ListenAndServe()
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.IndexPage{
		BasePage: s.basePage(r),
	})
}

func (s *Server) about(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.AboutPage{
		BasePage: s.basePage(r),
	})
}

func (s *Server) docs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, err := assetsDir.Open("docs.json")
	if err != nil {
		ctxlog.Error(ctx, "error reading docs.json", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
		return
	}

	var commands []*templates.CommandGroup
	if err := json.NewDecoder(file).Decode(&commands); err != nil {
		ctxlog.Error(ctx, "error unmarshaling docs.json", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
		return
	}

	templates.WritePageTemplate(w, &templates.DocsPage{
		BasePage: s.basePage(r),
		Commands: commands,
	})
}

func (s *Server) guilds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := s.sessionManager.GetString(ctx, sessionDiscordID)
	if id == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var guilds []*models.Guild

	dbGuilds, err := models.Guilds().All(ctx, s.DB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctxlog.Error(ctx, "error fetching guilds", zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
			return
		}
	}

	userGuilds, err := s.fetchGuilds(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/login", http.StatusOK)
			return
		}

		ctxlog.Error(ctx, "error fetching guilds", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
		return
	}

	for _, ug := range userGuilds {
		for _, dg := range dbGuilds {
			if strings.EqualFold(ug, dg.GuildID) {
				guilds = append(guilds, dg)
			}
		}
	}

	templates.WritePageTemplate(w, &templates.GuildsPage{
		BasePage: s.basePage(r),
		Guilds:   guilds,
	})
}

func (s *Server) guildPage(r *http.Request, g *models.Guild) templates.GuildPage {
	return templates.GuildPage{
		BasePage: s.basePage(r),
		Guild:    g,
	}
}

func (s *Server) guild(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	guild, err := modelsx.FetchGuild(ctx, s.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.errorPage(w, r, http.StatusNotFound, "")
		} else {
			ctxlog.Error(ctx, "error fetching guild", zap.String("guild_id", id), zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
		}
		return
	}

	templates.WritePageTemplate(w, &templates.GuildPage{
		BasePage: s.basePage(r),
		Guild:    guild,
	})
}

func (s *Server) guildQuotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	guild, err := modelsx.FetchGuild(ctx, s.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.errorPage(w, r, http.StatusNotFound, "")
		} else {
			ctxlog.Error(ctx, "error fetching guild", zap.String("guild_id", id), zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
		}
		return
	}

	quotes, err := models.Quotes(qm.Where("guild_id = ?", id), qm.OrderBy("num")).All(ctx, s.DB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctxlog.Error(ctx, "error fetching quotes", zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
			return
		}
	}

	templates.WritePageTemplate(w, &templates.GuildQuotesPage{
		GuildPage: s.guildPage(r, guild),
		Quotes:    quotes,
	})
}

func (s *Server) guildRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	guild, err := modelsx.FetchGuild(ctx, s.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.errorPage(w, r, http.StatusNotFound, "")
		} else {
			ctxlog.Error(ctx, "error fetching guild", zap.String("guild_id", id), zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
		}
		return
	}

	roles, err := models.Roles(qm.Where("guild_id = ?", id)).All(ctx, s.DB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctxlog.Error(ctx, "error fetching roles", zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
			return
		}
	}

	templates.WritePageTemplate(w, &templates.GuildRolesPage{
		GuildPage: s.guildPage(r, guild),
		Roles:     roles,
	})
}

func (s *Server) guildTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	guild, err := modelsx.FetchGuild(ctx, s.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.errorPage(w, r, http.StatusNotFound, "")
		} else {
			ctxlog.Error(ctx, "error fetching guild", zap.String("guild_id", id), zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
		}
		return
	}

	tags, err := models.Tags(qm.Where("guild_id = ?", id)).All(ctx, s.DB)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctxlog.Error(ctx, "error fetching tags", zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
			return
		}
	}

	templates.WritePageTemplate(w, &templates.GuildTagsPage{
		GuildPage: s.guildPage(r, guild),
		Tags:      tags,
	})
}
