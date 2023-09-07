package web

// Most of this is borrowed from HortBot's own web package.

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func requestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxlog.WithLogger(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func (s *Server) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				ctx := ctxlog.WithOptions(r.Context(), zap.AddStacktrace(zap.ErrorLevel))
				ctxlog.Error(ctx, "PANIC", zap.Any("val", rvr))

				s.errorPage(w, r, http.StatusInternalServerError, "")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (s *Server) guildCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := s.sessionManager.GetString(ctx, sessionDiscordID)
		if id == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		gid := chi.URLParam(r, "id")
		if gid == "" {
			s.errorPage(w, r, http.StatusBadRequest, "")
			return
		}

		exists, err := models.Guilds(qm.Where("guild_id = ?", gid)).Exists(ctx, s.DB)
		if err != nil {
			ctxlog.Error(ctx, "error checking if guild exists", zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
			return
		}

		if !exists {
			s.errorPage(w, r, http.StatusNotFound, "That guild ain't real, guy")
			return
		}

		guilds, err := s.fetchGuilds(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/login", http.StatusOK)
				return
			}

			ctxlog.Error(ctx, "error fetching guilds", zap.Error(err))
			s.errorPage(w, r, http.StatusInternalServerError, "")
			return
		}

		if !stringInSlice(gid, guilds) {
			s.errorPage(w, r, http.StatusForbidden, "CAN'T LET YOU DO THAT!!!!!")
		}

		next.ServeHTTP(w, r)
	})
}
