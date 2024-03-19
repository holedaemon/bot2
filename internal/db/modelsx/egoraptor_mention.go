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

// UpsertSetting inserts a mention, or updates it should it have changed.
func UpsertSetting(ctx context.Context, exec boil.ContextExecutor, m *models.EgoraptorSetting) error {
	return m.Upsert(ctx, exec, true, []string{"guild_id"}, mentionUpdate, boil.Infer())
}

// FetchSetting fetches a mention from the database, or returns a new one with defaults set.
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
