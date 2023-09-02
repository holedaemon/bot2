package modelsx

import (
	"context"

	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var guildUpdate = boil.Whitelist(
	models.GuildColumns.CreatedAt,
	models.GuildColumns.GuildName,
	models.GuildColumns.DoQuotes,
)

// UpsertGuild inserts a guild, or updates it should it have changed.
func UpsertGuild(ctx context.Context, exec boil.ContextExecutor, g *models.Guild) error {
	return g.Upsert(ctx, exec, true, []string{"guild_id"}, guildUpdate, boil.Infer())
}

// FetchGuild fetches a guild from the database.
func FetchGuild(ctx context.Context, exec boil.ContextExecutor, id string) (*models.Guild, error) {
	return models.Guilds(qm.Where("guild_id = ?", id)).One(ctx, exec)
}
