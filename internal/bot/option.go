package bot

import "go.uber.org/zap"

// Option configures a Bot.
type Option func(*Bot)

// WithLogger sets a Bot's logger.
func WithLogger(l *zap.Logger) Option {
	return func(b *Bot) {
		b.l = l
	}
}
