package web

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/holedaemon/bot2/internal/pkg/pgstore"
	"github.com/holedaemon/bot2/internal/web/templates"
	"github.com/patrickmn/go-cache"
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

	sessionManager *scs.SessionManager
	stateCache     *cache.Cache
	guildCache     *cache.Cache
}

// New creates a new Server.
func New(opts ...Option) (*Server, error) {
	srv := &Server{
		stateCache: cache.New(time.Hour, time.Hour*24),
		guildCache: cache.New(time.Minute*2, time.Minute*10),
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

	r.Route("/guild/{id}", func(r chi.Router) {
		r.Use(s.guildCheck)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.errorPage(w, r, http.StatusNotFound, "Whatever you're looking for ain't here")
	})

	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(assetsDir))))
	r.Handle("/favicon.ico", http.RedirectHandler("/static/favicon.ico", http.StatusFound))

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

func (s *Server) guilds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := s.sessionManager.GetString(ctx, sessionDiscordID)
	if id == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cache, err := s.fetchGuilds(ctx, id)
	if err != nil {
		if errors.Is(err, errTokenNotFound) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !errors.Is(err, errNoGuilds) {
			ctxlog.Error(ctx, "error fetching guilds", zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
		}
	}

	templates.WritePageTemplate(w, &templates.GuildsPage{
		BasePage: s.basePage(r),
		Guilds:   cache.ToSlice(),
	})
}
