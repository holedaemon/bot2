package bot

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

	guild, err := modelsx.FetchGuild(ctx, b.DB, g.ID.String())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error querying guild", zap.Error(err))
		}
	}

	if guild != nil {
		if guild.QuotesRequiredReactions.Valid {
			return
		}

		guild.QuotesRequiredReactions = null.IntFrom(1)

		if err := guild.Update(ctx, b.DB, boil.Whitelist(
			models.GuildColumns.UpdatedAt,
			models.GuildColumns.QuotesRequiredReactions,
		)); err != nil {
			log.Error("error fixing guild", zap.Error(err))
			return
		}

		log.Info("fixed null records for guild")
		return
	}

	guild = &models.Guild{
		GuildID:                 g.ID.String(),
		GuildName:               g.Name,
		QuotesRequiredReactions: null.IntFrom(1),
	}

	if err := guild.Insert(ctx, b.DB, boil.Infer()); err != nil {
		log.Error("error inserting guild into database", zap.Error(err))
		return
	}

	log.Info("created record for guild")
}

func (b *Bot) onGuildUpdate(g *gateway.GuildUpdateEvent) {
	ctx := context.Background()
	log := b.Logger.With(zap.String("guild_id", g.ID.String()))

	log.Info("guild settings have changed")

	guild, err := modelsx.FetchGuild(ctx, b.DB, g.ID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			guild = &models.Guild{
				GuildID: g.ID.String(),
			}
		} else {
			log.Error("error fetching guild", zap.Error(err))
			return
		}
	}

	if strings.EqualFold(g.Name, guild.GuildName) {
		return
	}

	guild.GuildName = g.Name

	if err := modelsx.UpsertGuild(ctx, b.DB, guild); err != nil {
		log.Error("error upserting guild record", zap.Error(err))
	}
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
