package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

const (
	updateInterval = (time.Hour * 24) * 7
	roleFmt        = "%d HOURS IN XIV"
)

func (b *Bot) roleUpdater(ctx context.Context) {
	t := time.NewTicker(updateInterval)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			settings, err := models.RoleUpdateSettings(qm.Where("guild_id = ?", scroteGuildID.String())).One(ctx, b.db)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}

				ctxlog.Error(ctx, "error fetching role update settings")
				continue
			}

			if !settings.DoUpdates {
				continue
			}

			if settings.LastTimestamp.Add(updateInterval).After(time.Now()) {
				continue
			}

			games, err := b.steam.GetOwnedGames(ctx, settings.SteamUserID)
			if err != nil {
				ctxlog.Error(ctx, "error fetching mettic's games")
				continue
			}

			for _, g := range games.Games {
				if g.AppID == settings.SteamAppID {
					hours := g.PlaytimeForever / 60
					sf, err := discord.ParseSnowflake(settings.RoleID)
					if err != nil {
						ctxlog.Error(ctx, "error parsing mettic's role as a snowflake", zap.Error(err))
						continue
					}

					roleID := discord.RoleID(sf)

					if _, err := b.state.ModifyRole(scroteGuildID, roleID, api.ModifyRoleData{
						Name: option.NewNullableString(fmt.Sprintf(roleFmt, hours)),
					}); err != nil {
						ctxlog.Error(ctx, "error updating mettic's role", zap.Error(err))
						continue
					}
				}
			}

			settings.LastTimestamp = time.Now()
			if err := settings.Update(
				ctx,
				b.db,
				boil.Whitelist(
					models.RoleUpdateSettingColumns.LastTimestamp,
					models.RoleUpdateSettingColumns.UpdatedAt,
				),
			); err != nil {
				ctxlog.Error(ctx, "error updating role updater settings in database", zap.Error(err))
			}
		}
	}
}
