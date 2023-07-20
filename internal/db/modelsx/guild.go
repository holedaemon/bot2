package modelsx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// FetchGuild fetches a guild from the database.
func FetchGuild(ctx context.Context, exec boil.ContextExecutor, id discord.GuildID) (*models.Guild, error) {
	guild, err := models.Guilds(qm.Where("guild_id = ?", id.String())).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return guild, nil
}
