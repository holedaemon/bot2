package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/bot"
	"github.com/holedaemon/bot2/internal/db/dbx"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

type Options struct {
	Debug             bool          `env:"BOT2_DEBUG" envDefault:"false"`
	Admins            string        `env:"BOT2_ADMINS"`
	Token             string        `env:"BOT2_TOKEN"`
	DSN               string        `env:"BOT2_DSN"`
	DBMaxAttempts     int           `env:"BOT2_DB_MAX_ATTEMPTS" envDefault:"10"`
	DBTimeoutDuration time.Duration `env:"BOT2_DB_TIMEOUT_DURATION" envDefault:"20s"`
}

func main() {
	opts := &Options{}
	eo := env.Options{
		RequiredIfNoDef: true,
	}

	if err := env.Parse(opts, eo); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing env variables into struct: %s\n", err.Error())
		return
	}

	logger := ctxlog.New(opts.Debug)
	ctx := ctxlog.WithLogger(context.Background(), logger)

	rawAdmins := strings.Split(opts.Admins, ",")
	admins := make(map[discord.UserID]struct{})
	for _, a := range rawAdmins {
		sf, err := discord.ParseSnowflake(a)
		if err != nil {
			logger.Fatal("error parsing admin snowflake", zap.Error(err))
		}

		if _, ok := admins[discord.UserID(sf)]; ok {
			continue
		}

		admins[discord.UserID(sf)] = struct{}{}
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

	b, err := bot.New(
		opts.Token,
		bot.WithAdminMap(admins),
		bot.WithDB(db),
		bot.WithDebug(opts.Debug),
	)
	if err != nil {
		logger.Fatal("error creating bot", zap.Error(err))
	}

	if err := b.Start(ctx); err != nil {
		logger.Fatal("error starting bot", zap.Error(err))
	}
}
