package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

const maxTimeoutSeconds = 604800

func (b *Bot) cmdEgoraptorToggle(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	toggled, err := data.Options.Find("toggled").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error converting argument into bool", zap.Error(err))
		return respondError("Error converting argument to boolean")
	}

	ego, err := modelsx.FetchMention(ctx, b.db, data.Event.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error querying for egoraptor mention", zap.Error(err))
		return dbError
	}

	ego.TimeoutOnMention = toggled

	if err := modelsx.UpsertMention(ctx, b.db, ego); err != nil {
		ctxlog.Error(ctx, "error upserting egoraptor mention", zap.Error(err))
		return dbError
	}

	if toggled {
		return respond("I will now time users out when they mention egoraptor eating pussy")
	}

	return respond("I will no longer time users out when they mention egoraptor eating pussy")
}

func (b *Bot) cmdEgoraptorSetTimeout(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	seconds, err := data.Options.Find("seconds").IntValue()
	if err != nil {
		ctxlog.Error(ctx, "error converting argument to int64", zap.Error(err))
		return respondError("Error converting argument to int64")
	}

	if seconds <= 0 {
		return respondError("The number of seconds must be a positive number greater than zero")
	}

	if seconds > maxTimeoutSeconds {
		return respondError("The number of seconds cannot exceed a week")
	}

	ego, err := modelsx.FetchMention(ctx, b.db, data.Event.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error querying for egoraptor mention", zap.Error(err))
		return dbError
	}

	ego.TimeoutLength = int(seconds)

	if err := modelsx.UpsertMention(ctx, b.db, ego); err != nil {
		ctxlog.Error(ctx, "error upserting egoraptor mention", zap.Error(err))
		return dbError
	}

	if seconds == 1 {
		return respond("Timeout length has been set to 1 second")
	}

	return respondf("Timeout length has been set to %d seconds", seconds)
}
