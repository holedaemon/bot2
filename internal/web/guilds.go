package web

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/web/templates"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (s *Server) guilds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := s.sessionManager.GetString(ctx, sessionDiscordID)
	if id == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	dbGuilds, err := models.Guilds().All(ctx, s.DB)
	if err != nil {
		ctxlog.Error(ctx, "error fetching guilds from database", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
		return
	}

	tok, err := s.fetchToken(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctxlog.Error(ctx, "error fetching token from database", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "")
	}

	cli := state.NewAPIOnlyState("Bearer "+tok.AccessToken, handler.New())
	apiGuilds, err := cli.GuildsAfter(0, 0)
	if err != nil {
		ctxlog.Error(ctx, "error fetching user guilds from Discord", zap.Error(err))
		s.errorPage(w, r, http.StatusInternalServerError, "Error fetching user guilds from Discord")
		return
	}

	guilds := make([]*templates.Guild, 0)

	for _, ag := range apiGuilds {
		for _, dg := range dbGuilds {
			if strings.EqualFold(ag.ID.String(), dg.GuildID) {
				guilds = append(guilds, &templates.Guild{
					ID:        ag.ID.String(),
					Name:      ag.Name,
					AvatarURL: ag.IconURLWithType(discord.JPEGImage),
				})
			}
		}
	}

	templates.WritePageTemplate(w, &templates.GuildsPage{
		BasePage: s.basePage(r),
		Guilds:   guilds,
	})
}
