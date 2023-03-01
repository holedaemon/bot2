package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) onReady(r *gateway.ReadyEvent) {
	b.l.Info("connected to Discord gateway", zap.Any("user_id", r.User.ID))
}

func (b *Bot) onMessage(m *gateway.MessageCreateEvent) {
	if m.Author.Bot {
		return
	}

	ctx := context.Background()
	ctx = ctxlog.WithLogger(ctx, b.l)
	ctx = ctxlog.With(ctx, zap.String("guild_id", m.GuildID.String()))

	if m.GuildID == holeGuildID {
		b.onHoleMessage(ctx, m)
	}
}
