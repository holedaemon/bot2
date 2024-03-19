package modelsx

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var azLoc *time.Location

func init() {
	loc, err := time.LoadLocation("America/Phoenix")
	if err != nil {
		panic(err)
	}

	azLoc = loc
}

// DefaultEgoraptorTimeout is the default number of seconds to timeout
// users.
const DefaultEgoraptorTimeout = 60

var mentionUpdate = boil.Whitelist(
	models.EgoraptorSettingColumns.UpdatedAt,
	models.EgoraptorSettingColumns.LastTimestamp,
	models.EgoraptorSettingColumns.LastUser,
	models.EgoraptorSettingColumns.Count,
	models.EgoraptorSettingColumns.TimeoutLength,
	models.EgoraptorSettingColumns.TimeoutOnMention,
)

// UpsertSetting inserts a setting, or updates it should it have changed.
func UpsertSetting(ctx context.Context, exec boil.ContextExecutor, m *models.EgoraptorSetting) error {
	return m.Upsert(ctx, exec, true, []string{"guild_id"}, mentionUpdate, boil.Infer())
}

// FetchSetting fetches settings from the database, or returns a new one with defaults set.
func FetchSetting(ctx context.Context, exec boil.ContextExecutor, id discord.GuildID) (*models.EgoraptorSetting, error) {
	data, err := models.EgoraptorSettings(qm.Where("guild_id = ?", id.String())).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data = &models.EgoraptorSetting{
				GuildID:          id.String(),
				LastTimestamp:    time.Now().In(azLoc),
				TimeoutLength:    DefaultEgoraptorTimeout,
				TimeoutOnMention: false,
			}
		} else {
			return nil, err
		}
	}

	return data, nil
}

// BumpMention bumps the count for an Egoraptor mention.
func BumpMention(ctx context.Context, exec boil.ContextExecutor, guildID discord.GuildID, userID discord.UserID) error {
	data, err := models.EgoraptorMentions(qm.Where("guild_id = ?", guildID.String()), qm.Where("user_id = ?", userID.String())).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data = &models.EgoraptorMention{
				GuildID: guildID.String(),
				UserID:  userID.String(),
				Count:   1,
			}

			return data.Insert(ctx, exec, boil.Infer())
		}

		return err
	}

	data.Count++

	if err := data.Update(ctx, exec, boil.Infer()); err != nil {
		return err
	}

	return nil
}

// FetchTopThreeMentions fetches the top three Egoraptor mentions for a guild.
func FetchTopThreeMentions(ctx context.Context, exec boil.ContextExecutor, id discord.GuildID) ([]*models.EgoraptorMention, error) {
	return models.EgoraptorMentions(qm.Where("guild_id = ?", id.String()), qm.OrderBy("count DESC"), qm.Limit(3)).All(ctx, exec)
}
