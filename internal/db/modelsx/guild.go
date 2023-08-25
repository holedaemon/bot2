package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// FetchGuild fetches a guild from the database.
func FetchGuild(ctx context.Context, exec boil.ContextExecutor, id string) (*models.Guild, error) {
	return models.Guilds(qm.Where("guild_id = ?", id)).One(ctx, exec)
}
