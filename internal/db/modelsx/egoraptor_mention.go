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

const DefaultEgoraptorTimeout = 60

var mentionUpdate = boil.Whitelist(
	models.EgoraptorMentionColumns.UpdatedAt,
	models.EgoraptorMentionColumns.LastTimestamp,
	models.EgoraptorMentionColumns.LastUser,
	models.EgoraptorMentionColumns.Count,
	models.EgoraptorMentionColumns.TimeoutLength,
	models.EgoraptorMentionColumns.TimeoutOnMention,
)

var azLoc *time.Location

func init() {
	loc, err := time.LoadLocation("America/Phoenix")
	if err != nil {
		panic(err)
	}

	azLoc = loc
}

// UpsertMention inserts a mention, or updates it should it have changed.
func UpsertMention(ctx context.Context, exec boil.ContextExecutor, m *models.EgoraptorMention) error {
	return m.Upsert(ctx, exec, true, []string{"guild_id"}, mentionUpdate, boil.Infer())
}

// FetchMention fetches a mention from the database, or returns a new one with defaults set.
func FetchMention(ctx context.Context, exec boil.ContextExecutor, id discord.GuildID) (*models.EgoraptorMention, error) {
	data, err := models.EgoraptorMentions(qm.Where("guild_id = ?", id.String())).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data = &models.EgoraptorMention{
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
