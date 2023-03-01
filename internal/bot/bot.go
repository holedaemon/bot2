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
	State  *state.State
	Logger *zap.Logger

	Jerkcity *jerkcity.Client
	DB       *sql.DB

	Admins map[discord.UserID]struct{}

	lastGameChange time.Time
}

// New creates a new Bot.
func New(token string, opts ...Option) (*Bot, error) {
	b := &Bot{}

	for _, o := range opts {
		o(b)
	}

	if b.Logger == nil {
		l, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("bot: creating logger: %w", err)
		}

		b.Logger = l
	}

	if b.Admins == nil {
		b.Admins = make(map[discord.UserID]struct{})
	}

	b.Jerkcity = jerkcity.New()

	b.State = state.New("Bot " + token)
	b.State.AddHandler(b.onReady)
	b.State.AddHandler(b.onGuildRoleDelete)
	b.State.AddHandler(b.onMessage)

	r := b.router()
	b.State.AddInteractionHandler(r)
	b.State.AddIntents(gateway.IntentGuilds | gateway.IntentGuildMessages)

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

	app, err := b.State.CurrentApplication()
	if err != nil {
		return nil, fmt.Errorf("bot: getting current application: %w", err)
	}

	for scope, cmd := range cmds {
		if scope == 0 {
			if _, err := b.State.BulkOverwriteCommands(app.ID, cmd); err != nil {
				return nil, fmt.Errorf("bot: overwriting global commands: %w", err)
			}
		} else {
			if _, err := b.State.BulkOverwriteGuildCommands(app.ID, scope, cmd); err != nil {
				return nil, fmt.Errorf("bot: overwriting guild commands (%d): %w", scope, err)
			}
		}
	}

	return b, nil
}

// IsAdmin checks if the given UserID is a bot admin.
func (b *Bot) IsAdmin(sf discord.UserID) bool {
	if _, ok := b.Admins[sf]; ok {
		return true
	}

	return false
}

// Start opens a connection to Discord.
func (b *Bot) Start(ctx context.Context) error {
	return b.State.Connect(ctx)
}
