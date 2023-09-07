package web

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/jellydator/ttlcache/v3"
)

func (s *Server) newAPIClient(ctx context.Context, id string) (*state.State, error) {
	tok, err := s.fetchToken(ctx, id)
	if err != nil {
		return nil, err
	}

	return state.NewAPIOnlyState("Bearer "+tok.AccessToken, handler.New()), nil
}

func (s *Server) notAuthorized(w http.ResponseWriter, r *http.Request, header bool) {
	if header {
		w.Header().Add("WWW-Authenticate", `Basic realm="dilf"`)
	}

	s.errorPage(w, r, http.StatusUnauthorized, "You need special pants to open this door")
}

// fetchToken retrieves a Discord token from the database, and refreshes it if necessary.
func (s *Server) fetchToken(ctx context.Context, id string) (*models.DiscordToken, error) {
	dt, err := modelsx.FetchToken(ctx, s.DB, id)
	if err != nil {
		return nil, err
	}

	if dt.Expiry.Before(time.Now()) {
		tok := modelsx.ModelToToken(dt)
		src := s.OAuth2.TokenSource(ctx, tok)
		newTok, err := src.Token()
		if err != nil {
			return nil, err
		}

		if tok.AccessToken != newTok.AccessToken {
			dt.AccessToken = newTok.AccessToken
			dt.Expiry = newTok.Expiry
			dt.TokenType = newTok.TokenType
			dt.RefreshToken = newTok.RefreshToken

			err = modelsx.UpsertDiscordToken(ctx, s.DB, dt)
			if err != nil {
				return nil, err
			}
		}
	}

	return dt, nil
}

func (s *Server) fetchGuilds(ctx context.Context, userID string) ([]string, error) {
	item := s.userCache.Get(userID)
	if item != nil {
		return item.Value(), nil
	}

	cli, err := s.newAPIClient(ctx, userID)
	if err != nil {
		return nil, err
	}

	guilds, err := cli.GuildsBefore(0, 0)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(guilds))

	for _, g := range guilds {
		ids = append(ids, g.ID.String())
	}

	item = s.userCache.Set(userID, ids, ttlcache.DefaultTTL)
	return item.Value(), nil
}

func stringInSlice(want string, in []string) bool {
	for _, s := range in {
		if strings.EqualFold(s, want) {
			return true
		}
	}

	return false
}
