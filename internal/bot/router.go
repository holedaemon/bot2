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
			ctx = ctxlog.WithLogger(ctx, b.logger)
			ctx = ctxlog.With(ctx, zap.String("guild_id", ie.GuildID.String()))
			return next.HandleInteraction(ctx, ie)
		})
	}
}

func (b *Bot) router() *cmdroute.Router {
	r := cmdroute.NewRouter()
	r.Use(b.useRouter())
	r.Use(cmdroute.Deferrable(b.state, cmdroute.DeferOpts{}))

	r.AddFunc("info", b.cmdInfo)
	r.AddFunc("help", b.cmdInfo)
	r.AddFunc("ping", b.cmdPing)
	r.AddFunc("is-admin", b.cmdIsAdmin)
	r.AddFunc("game", b.cmdGame)
	r.AddFunc("panic", b.cmdPanic)
	r.AddFunc("topster", b.cmdTopster)
	r.AddFunc("tag", b.cmdTag)

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
		r.AddFunc("fix", b.cmdRoleFix)
	})

	r.Sub("tags", func(r *cmdroute.Router) {
		r.AddFunc("create", b.cmdTagsCreate)
		r.AddFunc("update", b.cmdTagsUpdate)
		r.AddFunc("rename", b.cmdTagsRename)
		r.AddFunc("delete", b.cmdTagsDelete)
		r.AddFunc("list", b.cmdTagsList)
	})

	r.Sub("settings", func(r *cmdroute.Router) {
		r.Sub("quotes", func(r *cmdroute.Router) {
			r.AddFunc("toggle", b.cmdSettingsQuotesToggle)
			r.AddFunc("set-min-required", b.cmdSettingsQuotesSetMinRequired)
		})
	})

	r.Sub("profile", func(r *cmdroute.Router) {
		r.AddFunc("init", b.cmdProfileInit)
		r.AddFunc("delete", b.cmdProfileDelete)
		r.AddFunc("get", b.cmdProfileGet)
		r.AddFunc("set-timezone", b.cmdProfileSetTimezone)
	})

	r.AddFunc("q", b.cmdQ)
	r.Sub("quote", func(r *cmdroute.Router) {
		r.AddFunc("delete", b.cmdQuoteDelete)
	})

	r.Sub("egoraptor", func(r *cmdroute.Router) {
		r.AddFunc("toggle", b.cmdEgoraptorToggle)
		r.AddFunc("settimeout", b.cmdEgoraptorSetTimeout)
	})

	return r
}
