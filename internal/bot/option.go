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
		b.debug = debug
	}
}

// WithLogger sets a Bot's logger.
func WithLogger(l *zap.Logger) Option {
	return func(b *Bot) {
		b.logger = l
	}
}

// WithAdmin adds an admin to the Bot.
func WithAdmin(sf discord.UserID) Option {
	return func(b *Bot) {
		if b.admins == nil {
			b.admins = make(map[discord.UserID]struct{})
		}

		b.admins[sf] = struct{}{}
	}
}

// WithDB sets a Bot's DB.
func WithDB(db *sql.DB) Option {
	return func(b *Bot) {
		b.db = db
	}
}

// WithTopsterAddr sets the Bot's address for a Topster microservice.
func WithTopsterAddr(t string) Option {
	return func(b *Bot) {
		b.topsterAddr = t
	}
}

// WithWebhook sets the Bot's error log webhook.
func WithWebhook(hook *webhook.Client) Option {
	return func(b *Bot) {
		b.webhook = hook
	}
}

// WithSiteAddr sets a Bot's site address.
func WithSiteAddr(site string) Option {
	return func(b *Bot) {
		b.siteAddress = site
	}
}
