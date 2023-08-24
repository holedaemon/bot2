package web

import (
	"context"
	"time"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
)

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
