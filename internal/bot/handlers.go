package bot

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"go.uber.org/zap"
)

func (b *Bot) onReady(r *gateway.ReadyEvent) {
	b.l.Info("connected to Discord gateway", zap.Any("user_id", r.User.ID))
}
