package bot

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) onHoleMessage(ctx context.Context, m *gateway.MessageCreateEvent) {
	if roleInSlice(holeFortniteRoleID, m.MentionRoleIDs) {
		idx := rand.Intn(4) + 1
		str := strconv.FormatInt(int64(idx), 10)
		image := fakeCDN + "/fortnite-" + str + ".gif"
		_, err := b.State.SendMessage(m.ChannelID, image)
		if err != nil {
			ctxlog.Error(ctx, "error sending message", zap.Error(err))
			return
		}
	}
}
