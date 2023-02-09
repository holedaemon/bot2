package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"go.uber.org/zap"
)

const testGuildID = discord.GuildID(779875531712757800)

type Bot struct {
	s *state.State
	r *cmdroute.Router
	l *zap.Logger
}

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

	b.s = state.New("Bot " + token)
	b.s.AddHandler(b.onReady)

	b.r = cmdroute.NewRouter()
	b.r.AddFunc("ping", b.cmdPing)

	b.s.AddInteractionHandler(b.r)
	b.s.AddIntents(gateway.IntentGuilds)

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

func (b *Bot) Start(ctx context.Context) error {
	return b.s.Connect(ctx)
}
