package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) onReady(r *gateway.ReadyEvent) {
	b.Logger.Info("connected to Discord gateway", zap.Any("user_id", r.User.ID))
}

func (b *Bot) onReconnect(r *gateway.ReconnectEvent) {
	b.Logger.Info("reconnected to Discord gateway")
}

func (b *Bot) onGuildCreate(g *gateway.GuildCreateEvent) {
	ctx := context.Background()
	log := b.Logger.With(zap.String("guild_id", g.ID.String()))

	exists, err := models.Guilds(qm.Where("guild_id = ?", g.ID.String())).Exists(ctx, b.DB)
	if err != nil {
		log.Error("error querying guild", zap.Error(err))
		return
	}

	if exists {
		return
	}

	guild := &models.Guild{
		GuildID: g.ID.String(),
	}

	if err := guild.Insert(ctx, b.DB, boil.Infer()); err != nil {
		log.Error("error inserting guild into database", zap.Error(err))
		return
	}

	log.Info("created record for guild")
}

func (b *Bot) onMessage(m *gateway.MessageCreateEvent) {
	if m.Author.Bot {
		return
	}

	ctx := context.Background()
	ctx = ctxlog.WithLogger(ctx, b.Logger)
	ctx = ctxlog.With(ctx, zap.String("guild_id", m.GuildID.String()))

	if m.GuildID == holeGuildID {
		b.onHoleMessage(ctx, m)
		return
	}

	if m.GuildID == scroteGuildID {
		b.onScroteMessage(ctx, m)
		return
	}
}
