package bot

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/google/go-github/v56/github"
	"github.com/holedaemon/bot2/internal/api/jerkcity"
	"github.com/holedaemon/bot2/internal/api/steam"
	"github.com/holedaemon/bot2/internal/api/topster"
	"github.com/holedaemon/bot2/internal/pkg/imagecache"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Bot is a Discord bot account.
type Bot struct {
	debug bool

	topsterAddr string
	siteAddress string
	steamAPIKey string

	state           *state.State
	loggingWebhook  *webhook.Client
	feedbackWebhook *webhook.Client
	logger          *zap.Logger

	jerkcity *jerkcity.Client
	topster  *topster.Client
	steam    *steam.Client
	github   *github.Client
	db       *sql.DB

	admins map[discord.UserID]struct{}

	imageCache     *imagecache.Cache
	lastGameChange time.Time
	when           *when.Parser
}

// New creates a new Bot.
func New(token string, opts ...Option) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("bot: token is blank")
	}

	b := &Bot{
		jerkcity:   jerkcity.New(),
		imageCache: imagecache.New(),
		github:     github.NewClient(nil),
		admins:     make(map[discord.UserID]struct{}),
		state:      state.New("Bot " + token),
	}

	for _, o := range opts {
		o(b)
	}

	if b.logger == nil {
		l := ctxlog.New(b.debug)
		b.logger = l
	}

	// required options
	if b.siteAddress == "" {
		return nil, fmt.Errorf("bot: site address is blank")
	} else {
		if !strings.HasPrefix(b.siteAddress, "http") {
			b.siteAddress = "https://" + b.siteAddress
		}

		if !strings.HasSuffix(b.siteAddress, "/") {
			b.siteAddress = strings.TrimSuffix(b.siteAddress, "/")
		}
	}

	if b.db == nil {
		return nil, fmt.Errorf("bot: db is nil")
	}

	// optional options
	if b.loggingWebhook == nil {
		b.logger.Warn("webhook logs have been disabled")
	} else {
		b.logger = b.logger.WithOptions(
			zap.Hooks(b.webhookHook),
		)
	}

	if b.feedbackWebhook == nil {
		b.logger.Warn("feedback webhook has been disabled")
	}

	if len(b.admins) == 0 {
		b.logger.Warn("no admins have been configured; admin-only commands are unusable")
	}

	if b.topsterAddr == "" {
		b.logger.Warn("topster address has not been set; topster command is unusable")
	} else {
		if !strings.HasPrefix(b.topsterAddr, "http") {
			b.topsterAddr = "https://" + b.topsterAddr
		}

		topsterClient := &http.Client{
			Timeout: time.Second * 20,
		}

		tp, err := topster.New(b.topsterAddr, topster.WithHTTPClient(topsterClient))
		if err != nil {
			return nil, fmt.Errorf("%w: creating topster client", err)
		}

		b.topster = tp
	}

	if b.steamAPIKey == "" {
		b.logger.Warn("steam API key not set; role updater disabled")
	} else {
		stm, err := steam.New(b.steamAPIKey)
		if err != nil {
			return nil, fmt.Errorf("%w: creating steam client", err)
		}

		b.steam = stm
	}

	b.when = when.New(nil)
	b.when.Add(en.All...)
	b.when.Add(common.All...)

	b.state.AddHandler(b.onReady)
	b.state.AddHandler(b.onGuildCreate)
	b.state.AddHandler(b.onGuildUpdate)
	b.state.AddHandler(b.onGuildRoleDelete)
	b.state.AddHandler(b.onGuildMemberRemove)
	b.state.AddHandler(b.onMessage)
	b.state.AddHandler(b.onReconnect)
	b.state.AddHandler(b.onMessageReactionAdd)
	b.state.AddHandler(b.onMessageEdit)

	r := b.router()
	b.state.AddInteractionHandler(r)
	b.state.AddIntents(gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentGuildMessageReactions | gateway.IntentDirectMessages | gateway.IntentGuildMembers)

	if b.debug {
		commands = commands.Scoped(testGuildID)
	}

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

	app, err := b.state.CurrentApplication()
	if err != nil {
		return nil, fmt.Errorf("bot: getting current application: %w", err)
	}

	for scope, cmd := range cmds {
		if scope == 0 {
			if _, err := b.state.BulkOverwriteCommands(app.ID, cmd); err != nil {
				return nil, fmt.Errorf("bot: overwriting global commands: %w", err)
			}
		} else {
			if _, err := b.state.BulkOverwriteGuildCommands(app.ID, scope, cmd); err != nil {
				return nil, fmt.Errorf("bot: overwriting guild commands (%d): %w", scope, err)
			}
		}
	}

	return b, nil
}

// Start opens a connection to Discord and enables
// the internal image cache's automatic deletion.
func (b *Bot) Start(ctx context.Context) error {
	go b.imageCache.Start()
	defer b.imageCache.Stop()

	go b.roleUpdater(ctx)

	return b.state.Connect(ctx)
}

func (b *Bot) isAdmin(sf discord.UserID) bool {
	if _, ok := b.admins[sf]; ok {
		return true
	}

	return false
}

func (b *Bot) webhookHook(entry zapcore.Entry) error {
	if b.loggingWebhook == nil {
		return nil
	}

	if entry.Level < zapcore.ErrorLevel {
		return nil
	}

	data := webhook.ExecuteData{
		Username:  "BOT/2 Logs",
		AvatarURL: "https://holedaemon.net/images/yousuck.jpg",
		Content:   entry.Message,
	}

	return b.loggingWebhook.Execute(data)
}
