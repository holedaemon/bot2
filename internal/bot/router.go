package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) useRouter() cmdroute.Middleware {
	return func(next cmdroute.InteractionHandler) cmdroute.InteractionHandler {
		return cmdroute.InteractionHandlerFunc(func(ctx context.Context, ie *discord.InteractionEvent) *api.InteractionResponse {
			ctx = ctxlog.WithLogger(ctx, b.l)
			ctx = ctxlog.With(ctx, zap.String("guild_id", ie.GuildID.String()))
			return next.HandleInteraction(ctx, ie)
		})
	}
}

func (b *Bot) router() *cmdroute.Router {
	r := cmdroute.NewRouter()
	r.Use(b.useRouter())

	r.AddFunc("ping", b.cmdPing)
	r.AddFunc("is-admin", b.cmdIsAdmin)
	r.AddFunc("game", b.cmdGame)

	r.Sub("jerkcity", func(r *cmdroute.Router) {
		r.AddFunc("latest", b.cmdJerkcityLatest)
		r.AddFunc("episode", b.cmdJerkcityEpisode)
		r.AddFunc("quote", b.cmdJerkcityRandom)
		r.AddFunc("search", b.cmdJerkcitySearch)
	})

	r.Sub("role", func(r *cmdroute.Router) {
		r.AddFunc("create", b.cmdRoleCreate)
		r.AddFunc("add", b.cmdRoleAdd)
	})

	return r
}
