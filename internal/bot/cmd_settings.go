package bot

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) cmdSettingsQuotesToggle(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	val, err := data.Options.Find("toggled").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing bool value", zap.Error(err))
		return respondError("Error parsing bool value, oops")
	}

	guild, err := modelsx.FetchGuild(ctx, b.db, data.Event.GuildID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctxlog.Warn(ctx, "guild does not have record present in database")
			return respondError("A database record for this guild hasn't been created. This shouldn't happen. Uh oh!")
		}

		ctxlog.Error(ctx, "error querying guild in database", zap.Error(err))
		return dbError
	}

	guild.DoQuotes = val
	if err := guild.Update(ctx, b.db, boil.Whitelist(models.GuildColumns.DoQuotes, models.GuildColumns.UpdatedAt)); err != nil {
		ctxlog.Error(ctx, "error updating guild record", zap.Error(err))
		return dbError
	}

	if val {
		return respond("Quotes have been enabled")
	}

	return respond("Quotes have been disabled")
}

func (b *Bot) cmdSettingsQuotesSetMinRequired(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	val, err := data.Options.Find("minimum").IntValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing int value", zap.Error(err))
		return respondError("Error parsing the given value into an int!!!!")
	}

	if val <= 0 {
		return respondError("Your value must be a number greater than 0")
	}

	if val > 10 {
		return respondError("Your value must be less than or equal to 10")
	}

	guild, err := modelsx.FetchGuild(ctx, b.db, data.Event.GuildID.String())
	if err != nil {
		ctxlog.Error(ctx, "error fetching guild from database", zap.Error(err))
		return dbError
	}

	newVal := int(val)

	guild.QuotesRequiredReactions = null.IntFrom(newVal)
	if err := guild.Update(ctx, b.db, boil.Whitelist(models.GuildColumns.UpdatedAt, models.GuildColumns.QuotesRequiredReactions)); err != nil {
		ctxlog.Error(ctx, "error updating guild record in database", zap.Error(err))
		return dbError
	}

	return respondf("The required number of reactions to add a quote has been set to %d", newVal)
}
