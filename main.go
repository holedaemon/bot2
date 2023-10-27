package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	// Per the docs, this should be imported in the main package rather than a library.
	// We use this for the profile timezone function
	_ "time/tzdata"

	"github.com/caarlos0/env/v7"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/bot"
	"github.com/holedaemon/bot2/internal/db/dbx"
	"github.com/holedaemon/bot2/internal/pkg/discordx"
	"github.com/holedaemon/bot2/internal/web"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type BotOptions struct {
	Debug             bool          `env:"BOT2_DEBUG" envDefault:"false"`
	Admins            []string      `env:"BOT2_ADMINS"`
	Token             string        `env:"BOT2_TOKEN"`
	WebhookURL        string        `env:"BOT2_WEBHOOK_URL"`
	DSN               string        `env:"BOT2_DSN"`
	TopsterAddr       string        `env:"BOT2_TOPSTER_ADDR"`
	SiteAddr          string        `env:"BOT2_SITE_ADDR"`
	SteamAPIKey       string        `env:"BOT2_STEAM_API_KEY"`
	DBMaxAttempts     int           `env:"BOT2_DB_MAX_ATTEMPTS" envDefault:"10"`
	DBTimeoutDuration time.Duration `env:"BOT2_DB_TIMEOUT_DURATION" envDefault:"20s"`
}

type WebOptions struct {
	Debug  bool              `env:"BOT2_WEB_DEBUG" envDefault:"false"`
	Addr   string            `env:"BOT2_WEB_ADDR" envDefault:":8080"`
	Admins map[string]string `env:"BOT2_WEB_ADMINS"`

	DSN               string        `env:"BOT2_WEB_DSN"`
	DBMaxAttempts     int           `env:"BOT2_WEB_DB_MAX_ATTEMPTS" envDefault:"10"`
	DBTimeoutDuration time.Duration `env:"BOT2_WEB_DB_TIMEOUT_DURATION" envDefault:"20s"`

	OAuth2ClientID     string   `env:"BOT2_WEB_OAUTH2_CLIENT_ID"`
	OAuth2ClientSecret string   `env:"BOT2_WEB_OAUTH2_CLIENT_SECRET"`
	OAuth2Scopes       []string `env:"BOT2_WEB_OAUTH2_SCOPES"`
	OAuth2RedirectURL  string   `env:"BOT2_WEB_OAUTH2_REDIRECT_URL"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "missing required argument: bot, web")
		return
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "bot":
		runBot()
	case "web":
		runWeb()
	default:
		fmt.Fprintf(os.Stderr, "invalid argument; accepted arguments: bot, web")
	}
}

func runBot() {
	opts := &BotOptions{}
	eo := env.Options{
		RequiredIfNoDef: true,
	}

	if err := env.Parse(opts, eo); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env variables into struct: %s\n", err.Error())
		return
	}

	logger := ctxlog.New(opts.Debug)
	ctx := ctxlog.WithLogger(context.Background(), logger)

	admins := make([]discord.UserID, 0)
	for _, a := range opts.Admins {
		sf, err := discord.ParseSnowflake(a)
		if err != nil {
			logger.Fatal("error parsing admin snowflake", zap.Error(err))
		}

		admins = append(admins, discord.UserID(sf))
	}

	var (
		db  *sql.DB
		err error
	)

	connected := false
	for i := 0; i < opts.DBMaxAttempts && !connected; i++ {
		db, err = sql.Open(dbx.Driver, opts.DSN)
		if err != nil {
			logger.Error("unable to connect to database", zap.Error(err), zap.Int("attempt", i))
			time.Sleep(opts.DBTimeoutDuration)
			continue
		}

		if err = db.PingContext(ctx); err != nil {
			logger.Error("unable to ping database", zap.Error(err), zap.Int("attempt", i))
			time.Sleep(opts.DBTimeoutDuration)
			continue
		}

		connected = true
	}

	if !connected {
		logger.Fatal("max database attempts reached", zap.Int("attempts", opts.DBMaxAttempts))
	}

	hook, err := webhook.NewFromURL(opts.WebhookURL)
	if err != nil {
		logger.Fatal("error creating webhook", zap.Error(err))
	}

	b, err := bot.New(
		opts.Token,
		bot.WithDB(db),
		bot.WithWebhook(hook),
		bot.WithDebug(opts.Debug),
		bot.WithTopsterAddr(opts.TopsterAddr),
		bot.WithSiteAddr(opts.SiteAddr),
		bot.WithSteamAPIKey(opts.SteamAPIKey),
		bot.WithAdmins(admins...),
	)
	if err != nil {
		logger.Fatal("error creating bot", zap.Error(err))
	}

	if err := b.Start(ctx); err != nil {
		logger.Fatal("error starting bot", zap.Error(err))
	}
}

func runWeb() {
	opts := &WebOptions{}
	eo := env.Options{
		RequiredIfNoDef: true,
	}

	if err := env.Parse(opts, eo); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env variables into struct: %s\n", err.Error())
		return
	}

	logger := ctxlog.New(opts.Debug)
	ctx := ctxlog.WithLogger(context.Background(), logger)

	oa := &oauth2.Config{
		ClientID:     opts.OAuth2ClientID,
		ClientSecret: opts.OAuth2ClientSecret,
		Endpoint:     discordx.Endpoint,
		RedirectURL:  opts.OAuth2RedirectURL,
		Scopes:       opts.OAuth2Scopes,
	}

	var (
		db  *sql.DB
		err error
	)

	connected := false
	for i := 0; i < opts.DBMaxAttempts && !connected; i++ {
		db, err = sql.Open(dbx.Driver, opts.DSN)
		if err != nil {
			logger.Error("unable to connect to database", zap.Error(err), zap.Int("attempt", i))
			time.Sleep(opts.DBTimeoutDuration)
			continue
		}

		if err = db.PingContext(ctx); err != nil {
			logger.Error("unable to ping database", zap.Error(err), zap.Int("attempt", i))
			time.Sleep(opts.DBTimeoutDuration)
			continue
		}

		connected = true
	}

	if !connected {
		logger.Fatal("max database attempts reached", zap.Int("attempts", opts.DBMaxAttempts))
	}

	s, err := web.New(
		web.WithDebug(opts.Debug),
		web.WithAddr(opts.Addr),
		web.WithDB(db),
		web.WithOAuth2(oa),
		web.WithAdmins(opts.Admins),
	)

	if err != nil {
		logger.Fatal("error creating bot", zap.Error(err))
	}

	if err := s.Run(ctx); err != nil {
		logger.Fatal("error starting bot", zap.Error(err))
	}
}
