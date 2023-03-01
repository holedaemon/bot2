package bot

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/holedaemon/bot2/internal/api/jerkcity"
	"go.uber.org/zap"
)

// Bot is a Discord bot account.
type Bot struct {
	s *state.State
	l *zap.Logger

	jc *jerkcity.Client
	db *sql.DB

	admins map[discord.UserID]struct{}

	lastGameChange time.Time
}

// New creates a new Bot.
func New(token string, opts ...Option) (*Bot, error) {
	b := &Bot{}

	for _, o := range opts {
		o(b)
	}

	if b.l == nil {
		l, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("bot: creating logger: %w", err)
		}

		b.l = l
	}

	if b.admins == nil {
		b.admins = make(map[discord.UserID]struct{})
	}

	b.jc = jerkcity.New()

	b.s = state.New("Bot " + token)
	b.s.AddHandler(b.onReady)
	b.s.AddHandler(b.onGuildRoleDelete)
	b.s.AddHandler(b.onMessage)

	r := b.router()
	b.s.AddInteractionHandler(r)
	b.s.AddIntents(gateway.IntentGuilds | gateway.IntentGuildMessages)

	cmds := make(map[discord.GuildID][]api.CreateCommandData)
	for _, c := range commands {
		if c.Scoped.IsNull() {
			if _, ok := cmds[0]; !ok {
				cmds[0] = make([]api.CreateCommandData, 0)
			}

			cmds[0] = append(cmds[0], c.Data)
		} else {
			if _, ok := cmds[c.Scoped]; !ok {
				cmds[c.Scoped] = make([]api.CreateCommandData, 0)
			}

			cmds[c.Scoped] = append(cmds[c.Scoped], c.Data)
		}
	}

	app, err := b.s.CurrentApplication()
	if err != nil {
		return nil, fmt.Errorf("bot: getting current application: %w", err)
	}

	for scope, cmd := range cmds {
		if scope == 0 {
			if _, err := b.s.BulkOverwriteCommands(app.ID, cmd); err != nil {
				return nil, fmt.Errorf("bot: overwriting global commands: %w", err)
			}
		} else {
			if _, err := b.s.BulkOverwriteGuildCommands(app.ID, scope, cmd); err != nil {
				return nil, fmt.Errorf("bot: overwriting guild commands (%d): %w", scope, err)
			}
		}
	}

	return b, nil
}

// IsAdmin checks if the given UserID is a bot admin.
func (b *Bot) IsAdmin(sf discord.UserID) bool {
	if _, ok := b.admins[sf]; ok {
		return true
	}

	return false
}

// Start opens a connection to Discord.
func (b *Bot) Start(ctx context.Context) error {
	return b.s.Connect(ctx)
}
