package web

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/holedaemon/bot2/internal/web/templates"
	"github.com/patrickmn/go-cache"
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

var (
	errNoGuilds      = errors.New("web: no guilds in the database")
	errTokenNotFound = errors.New("web: token not found")
)

func (s *Server) fetchGuilds(ctx context.Context, id string) (*cachedGuilds, error) {
	gc, found := s.guildCache.Get(id)
	if found {
		cache, ok := gc.(*cachedGuilds)
		if !ok {
			return nil, nil
		}

		return cache, nil
	}

	dbGuilds, err := models.Guilds().All(ctx, s.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNoGuilds
		}

		return nil, err
	}

	tok, err := s.fetchToken(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTokenNotFound
		}

		return nil, err
	}

	cli := state.NewAPIOnlyState("Bearer "+tok.AccessToken, handler.New())
	apiGuilds, err := cli.GuildsBefore(0, 0)
	if err != nil {
		return nil, err
	}

	guilds := newGuildCache()
	for _, ag := range apiGuilds {
		for _, dg := range dbGuilds {
			if strings.EqualFold(ag.ID.String(), dg.GuildID) {
				guilds.Add(&templates.Guild{
					ID:   ag.ID.String(),
					Name: ag.Name,
				})
			}
		}
	}

	s.guildCache.Set(id, guilds, cache.DefaultExpiration)
	return s.fetchGuilds(ctx, id)
}
