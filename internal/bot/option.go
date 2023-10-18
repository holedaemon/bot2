package bot

import (
	"database/sql"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"go.uber.org/zap"
)

// Option configures a Bot.
type Option func(*Bot)

// WithDebug toggles debug mode on a Bot.
func WithDebug(debug bool) Option {
	return func(b *Bot) {
		b.Debug = debug
	}
}

// WithLogger sets a Bot's logger.
func WithLogger(l *zap.Logger) Option {
	return func(b *Bot) {
		b.Logger = l
	}
}

// WithAdminMap sets a Bot's admin map.
func WithAdminMap(m map[discord.UserID]struct{}) Option {
	return func(b *Bot) {
		b.Admins = m
	}
}

// WithDB sets a Bot's DB.
func WithDB(db *sql.DB) Option {
	return func(b *Bot) {
		b.DB = db
	}
}

// WithTopsterAddr sets the Bot's address for a Topster microservice.
func WithTopsterAddr(t string) Option {
	return func(b *Bot) {
		b.TopsterAddr = t
	}
}

// WithWebhook sets the Bot's error log webhook.
func WithWebhook(hook *webhook.Client) Option {
	return func(b *Bot) {
		b.Webhook = hook
	}
}

// WithSiteAddr sets a Bot's site address.
func WithSiteAddr(site string) Option {
	return func(b *Bot) {
		b.siteAddress = site
	}
}
