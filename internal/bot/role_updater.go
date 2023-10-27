package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

const (
	updateInterval = (time.Hour * 24) * 7
	// updateInterval = time.Minute * 2
	roleFmt = "%d HOURS IN XIV"
)

func (b *Bot) roleUpdater(ctx context.Context) {
	t := time.NewTicker(updateInterval)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			updater, err := modelsx.FetchRoleUpdater(ctx, b.db, scroteGuildID.String())
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}

				ctxlog.Error(ctx, "error fetching role update settings")
				continue
			}

			if !updater.DoUpdates {
				ctxlog.Debug(ctx, "role updates are disabled, skipping...")
				continue
			}

			games, err := b.steam.GetOwnedGames(ctx, metticSteamID)
			if err != nil {
				ctxlog.Error(ctx, "error fetching mettic's games")
				continue
			}

			for _, g := range games.Games {
				if g.AppID == xivAppID {
					hours := g.PlaytimeForever / 60

					if _, err := b.state.ModifyRole(scroteGuildID, metticRoleID, api.ModifyRoleData{
						Name: option.NewNullableString(fmt.Sprintf(roleFmt, hours)),
					}); err != nil {
						ctxlog.Error(ctx, "error updating mettic's role", zap.Error(err))
						continue
					}
				}
			}

			ctxlog.Debug(ctx, "updated mettic's role")
		}
	}
}
