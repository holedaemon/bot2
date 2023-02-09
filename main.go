package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/bot"
	"go.uber.org/zap"
)

func main() {
	rd := os.Getenv("BOT2_DEBUG")
	debug := rd != ""

	var (
		logger *zap.Logger
		err    error
	)

	if debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalf("error creating new development logger: %s\n", err.Error())
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("error creating new production logger: %s\n", err.Error())
		}
	}

	rawAdmins := os.Getenv("BOT2_ADMINS")
	admins := make(map[discord.UserID]struct{})
	if rawAdmins != "" {
		adminsSplit := strings.Split(rawAdmins, ",")
		for _, a := range adminsSplit {
			sf, err := discord.ParseSnowflake(a)
			if err != nil {
				logger.Fatal("error parsing admin ID into snowflake", zap.Error(err), zap.String("id", a))
			}

			admins[discord.UserID(sf)] = struct{}{}
		}
	}

	token := os.Getenv("BOT2_TOKEN")
	if token == "" {
		logger.Fatal("$BOT2_TOKEN is blank")
	}

	b, err := bot.New(token, bot.WithLogger(logger), bot.WithAdminMap(admins))
	if err != nil {
		logger.Fatal("error creating bot", zap.Error(err))
	}

	ctx := context.Background()
	if err := b.Start(ctx); err != nil {
		logger.Fatal("error starting bot", zap.Error(err))
	}
}
