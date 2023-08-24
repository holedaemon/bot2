package web

import (
	"net/http"

	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/gofrs/uuid"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/patrickmn/go-cache"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

const (
	sessionName      = "bot2-session-v1"
	sessionDiscordID = "discord_id"
	sessionUsername  = "discord_username"
)

func (s *Server) authDiscord(w http.ResponseWriter, r *http.Request) {
	state := uuid.Must(uuid.NewV4()).String()

	s.stateCache.Set(state, struct{}{}, cache.DefaultExpiration)
	url := s.OAuth2.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (s *Server) authDiscordCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	st := r.FormValue("state")
	if st == "" {
		s.errorPage(w, r, http.StatusBadRequest, "")
		return
	}

	if _, ok := s.stateCache.Get(st); !ok {
		s.errorPage(w, r, http.StatusBadRequest, "Unexpected state")
		return
	}

	tok, err := s.OAuth2.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		ctxlog.Error(ctx, "error exchanging code", zap.Error(err))
		s.errorPage(w, r, http.StatusBadRequest, "")
		return
	}

	cli := state.NewAPIOnlyState("Bearer "+tok.AccessToken, handler.New())
	user, err := cli.Me()
	if err != nil {
		ctxlog.Error(ctx, "error fetching user from Discord", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
		return
	}

	dt := &models.DiscordToken{
		UserID:       user.ID.String(),
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
		TokenType:    tok.TokenType,
	}

	if err := modelsx.UpsertDiscordToken(ctx, s.DB, dt); err != nil {
		ctxlog.Error(ctx, "error upserting token", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
		return
	}

	if err := s.sessionManager.RenewToken(ctx); err != nil {
		ctxlog.Error(ctx, "error renewing token", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
		return
	}

	s.sessionManager.Put(ctx, sessionDiscordID, user.ID.String())
	s.sessionManager.Put(ctx, sessionUsername, user.Username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := s.sessionManager.Destroy(ctx); err != nil {
		ctxlog.Error(ctx, "error destroying session", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
	}
}
