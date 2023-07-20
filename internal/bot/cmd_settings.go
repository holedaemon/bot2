package bot

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) cmdSettingsQuotes(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	val, err := data.Options.Find("toggled").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing bool value", zap.Error(err))
		return respondError("Error parsing bool value, oops")
	}

	guild, err := models.Guilds(qm.Where("guild_id = ?", data.Event.GuildID.String())).One(ctx, b.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctxlog.Warn(ctx, "guild does not have record present in database")
			return respondError("A database record for this guild hasn't been created. This shouldn't happen. Uh oh!")
		}

		ctxlog.Error(ctx, "error querying guild in database", zap.Error(err))
		return dbError
	}

	guild.DoQuotes = val
	if err := guild.Update(ctx, b.DB, boil.Infer()); err != nil {
		ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
		return dbError
	}

	if val {
		return respond("Quotes have been enabled")
	}

	return respond("Quotes have been disabled")
}
