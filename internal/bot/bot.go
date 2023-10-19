package bot

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/holedaemon/bot2/internal/api/jerkcity"
	"github.com/holedaemon/lastfm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Bot is a Discord bot account.
type Bot struct {
	debug bool

	topsterAddr string
	siteAddress string

	state   *state.State
	webhook *webhook.Client
	logger  *zap.Logger

	lastfm   *lastfm.Client
	jerkcity *jerkcity.Client
	db       *sql.DB

	admins map[discord.UserID]struct{}

	imageCache     *ImageCache
	lastGameChange time.Time
}

// New creates a new Bot.
func New(token string, opts ...Option) (*Bot, error) {
	b := &Bot{}

	for _, o := range opts {
		o(b)
	}

	if b.logger == nil {
		l, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("bot: creating logger: %w", err)
		}

		b.logger = l
	}

	if b.webhook != nil {
		b.logger = b.logger.WithOptions(
			zap.Hooks(b.webhookHook),
		)
	}

	if b.admins == nil {
		b.admins = make(map[discord.UserID]struct{})
	}

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

	b.jerkcity = jerkcity.New()
	b.imageCache = NewImageCache()

	b.state = state.New("Bot " + token)
	b.state.AddHandler(b.onReady)
	b.state.AddHandler(b.onGuildCreate)
	b.state.AddHandler(b.onGuildUpdate)
	b.state.AddHandler(b.onGuildRoleDelete)
	b.state.AddHandler(b.onMessage)
	b.state.AddHandler(b.onReconnect)
	b.state.AddHandler(b.onMessageReactionAdd)
	b.state.AddHandler(b.onMessageEdit)

	r := b.router()
	b.state.AddInteractionHandler(r)
	b.state.AddIntents(gateway.IntentGuilds | gateway.IntentGuildMessages | gateway.IntentGuildMessageReactions)

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

// IsAdmin checks if the given UserID is a bot admin.
func (b *Bot) IsAdmin(sf discord.UserID) bool {
	if _, ok := b.admins[sf]; ok {
		return true
	}

	return false
}

// Start opens a connection to Discord.
func (b *Bot) Start(ctx context.Context) error {
	return b.state.Connect(ctx)
}

func (b *Bot) webhookHook(entry zapcore.Entry) error {
	if b.webhook == nil {
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

	return b.webhook.Execute(data)
}

func (b *Bot) Reply(m discord.Message, content string) error {
	if content == "" {
		panic("bot: blank content given to Reply")
	}

	_, err := b.state.SendMessageReply(m.ChannelID, content, m.ID)
	return err
}

func (b *Bot) Replyf(m discord.Message, content string, args ...interface{}) error {
	msg := fmt.Sprintf(content, args...)
	return b.Reply(m, msg)
}

func (b *Bot) SendImage(chID discord.ChannelID, content string, image string) error {
	cachedImage := b.imageCache.Get(image)
	if cachedImage == nil {
		err := b.imageCache.Download(image)
		if err != nil {
			return err
		}

		cachedImage = b.imageCache.Get(image)
	}

	rawName := path.Base(image)

	files := make([]sendpart.File, 0)
	files = append(files, sendpart.File{
		Name:   rawName,
		Reader: cachedImage,
	})

	_, err := b.state.SendMessageComplex(
		chID,
		api.SendMessageData{
			Content: content,
			Files:   files,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
