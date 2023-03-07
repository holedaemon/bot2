package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
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

	ego, err := b.loadEgoraptorData()
	if err != nil {
		ctxlog.Error(ctx, "error loading egoraptor data", zap.Error(err))
		return respondError("Error loading egoraptor data")
	}

	ego.SetTimeout(toggled)

	err = b.writeEgoraptorData()
	if err != nil {
		ctxlog.Error(ctx, "error writing egoraptor data", zap.Error(err))
		return respond("Error saving settings")
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

	ego, err := b.loadEgoraptorData()
	if err != nil {
		ctxlog.Error(ctx, "error loading egoraptor data", zap.Error(err))
		return respondError("Error loading egoraptor data")
	}

	ego.SetTimeoutLength(seconds)

	err = b.writeEgoraptorData()
	if err != nil {
		ctxlog.Error(ctx, "error writing egoraptor data", zap.Error(err))
		return respond("Error saving settings")
	}

	if seconds == 1 {
		return respond("Timeout length has been set to 1 second")
	}

	return respondf("Timeout length has been set to %d seconds", seconds)
}
