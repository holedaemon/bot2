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
			ctx = ctxlog.WithLogger(ctx, b.Logger)
			ctx = ctxlog.With(ctx, zap.String("guild_id", ie.GuildID.String()))
			return next.HandleInteraction(ctx, ie)
		})
	}
}

func recoverer() cmdroute.Middleware {
	return func(next cmdroute.InteractionHandler) cmdroute.InteractionHandler {
		return cmdroute.InteractionHandlerFunc(func(ctx context.Context, ie *discord.InteractionEvent) *api.InteractionResponse {
			defer func() {
				if rvr := recover(); rvr != nil {
					ctx := ctxlog.WithOptions(ctx, zap.AddStacktrace(zap.ErrorLevel))
					ctxlog.Error(ctx, "PANIC", zap.Any("recover", rvr))
				}
			}()

			return next.HandleInteraction(ctx, ie)
		})
	}
}

func (b *Bot) router() *cmdroute.Router {
	r := cmdroute.NewRouter()
	r.Use(b.useRouter())
	r.Use(recoverer())

	r.AddFunc("info", b.cmdInfo)
	r.AddFunc("help", b.cmdInfo)
	r.AddFunc("ping", b.cmdPing)
	r.AddFunc("is-admin", b.cmdIsAdmin)
	r.AddFunc("game", b.cmdGame)
	r.AddFunc("panic", b.cmdPanic)

	r.Sub("jerkcity", func(r *cmdroute.Router) {
		r.AddFunc("latest", b.cmdJerkcityLatest)
		r.AddFunc("episode", b.cmdJerkcityEpisode)
		r.AddFunc("quote", b.cmdJerkcityRandom)
		r.AddFunc("search", b.cmdJerkcitySearch)
	})

	r.Sub("role", func(r *cmdroute.Router) {
		r.AddFunc("add", b.cmdRoleAdd)
		r.AddFunc("remove", b.cmdRoleRemove)
		r.AddFunc("list", b.cmdRoleList)
	})

	r.Sub("roles", func(r *cmdroute.Router) {
		r.AddFunc("create", b.cmdRoleCreate)
		r.AddFunc("delete", b.cmdRoleDelete)
		r.AddFunc("rename", b.cmdRoleRename)
		r.AddFunc("setcolor", b.cmdRoleSetColor)
		r.AddFunc("import", b.cmdRoleImport)
		r.AddFunc("relinquish", b.cmdRoleRelinquish)
	})

	r.Sub("egoraptor", func(r *cmdroute.Router) {
		r.AddFunc("toggle", b.cmdEgoraptorToggle)
		r.AddFunc("settimeout", b.cmdEgoraptorSetTimeout)
	})

	r.AddFunc("q", b.cmdQ)

	return r
}
