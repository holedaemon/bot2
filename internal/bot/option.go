package bot

import (
	"database/sql"

	"github.com/diamondburned/arikawa/v3/discord"
	"go.uber.org/zap"
)

// Option configures a Bot.
type Option func(*Bot)

// WithLogger sets a Bot's logger.
func WithLogger(l *zap.Logger) Option {
	return func(b *Bot) {
		b.l = l
	}
}

// WithAdminMap sets a Bot's admin map.
func WithAdminMap(m map[discord.UserID]struct{}) Option {
	return func(b *Bot) {
		b.admins = m
	}
}

// WithDB sets a Bot's DB.
func WithDB(db *sql.DB) Option {
	return func(b *Bot) {
		b.db = db
	}
}
