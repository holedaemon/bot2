package bot

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) cmdStfu(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	var name string
	me, err := b.state.Me()
	if err != nil {
		ctxlog.Error(ctx, "error fetching username", zap.Error(err))
		name = "BOT"
	} else {
		name = me.Username
	}

	selfDestruct := rand.Intn(10000)
	msg := fmt.Sprintf("[%s-GPT] %s has detected an argument in this channel. Please diverge into a thread as to not annoy those that aren't participating. Failure to comply will result in explosion. This message will self-destruct in %d seconds.", name, name, selfDestruct)

	_, err = b.state.SendMessage(data.Event.ChannelID, msg)
	if err != nil {
		ctxlog.Error(ctx, "error sending message", zap.Error(err))
		return respondError("Error sending message!")
	}

	return respondSilent("Told everyone to stfu.")
}
