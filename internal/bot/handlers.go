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
	b.logger.Info("connected to Discord gateway", zap.Any("user_id", r.User.ID))
}

func (b *Bot) onReconnect(r *gateway.ReconnectEvent) {
	b.logger.Info("reconnected to Discord gateway")
}

func (b *Bot) onGuildCreate(g *gateway.GuildCreateEvent) {
	ctx := context.Background()
	log := b.logger.With(zap.String("guild_id", g.ID.String()))

	guild, err := modelsx.FetchGuild(ctx, b.db, g.ID.String())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error querying guild", zap.Error(err))
		}
	}

	if guild != nil {
		whitelist := make([]string, 0)

		if !guild.QuotesRequiredReactions.Valid {
			guild.QuotesRequiredReactions = null.IntFrom(1)
			whitelist = append(whitelist, models.GuildColumns.QuotesRequiredReactions)
		}

		if !strings.EqualFold(guild.GuildName, g.Name) {
			guild.GuildName = g.Name
			whitelist = append(whitelist, models.GuildColumns.GuildName)
		}

		if len(whitelist) > 0 {
			whitelist = append(whitelist, models.GuildColumns.UpdatedAt)

			if err := guild.Update(ctx, b.db, boil.Whitelist(whitelist...)); err != nil {
				log.Error("error fixing guild", zap.Error(err))
				return
			}

			log.Info("updated record for guild")
		}

		return
	}

	guild = &models.Guild{
		GuildID:                 g.ID.String(),
		GuildName:               g.Name,
		QuotesRequiredReactions: null.IntFrom(1),
	}

	if err := guild.Insert(ctx, b.db, boil.Infer()); err != nil {
		log.Error("error inserting guild into database", zap.Error(err))
		return
	}

	log.Info("created record for guild")
}

func (b *Bot) onGuildUpdate(g *gateway.GuildUpdateEvent) {
	ctx := context.Background()
	log := b.logger.With(zap.String("guild_id", g.ID.String()))

	log.Info("guild settings have changed")

	guild, err := modelsx.FetchGuild(ctx, b.db, g.ID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctxlog.Error(ctx, "received update for guild we aren't tracking...")
			return
		}

		log.Error("error fetching guild", zap.Error(err))
		return
	}

	if strings.EqualFold(g.Name, guild.GuildName) {
		return
	}

	guild.GuildName = g.Name
	if err := guild.Update(ctx, b.db, boil.Whitelist(models.GuildColumns.GuildName, models.GuildColumns.UpdatedAt)); err != nil {
		ctxlog.Error(ctx, "error updating guild", zap.Error(err))
		return
	}
}

func (b *Bot) onMessage(m *gateway.MessageCreateEvent) {
	if m.Author.Bot {
		return
	}

	ctx := context.Background()
	ctx = ctxlog.WithLogger(ctx, b.logger)
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

func (b *Bot) onGuildMemberRemove(e *gateway.GuildMemberRemoveEvent) {
	ctx := context.Background()
	ctx = ctxlog.WithLogger(ctx, b.logger)
	ctx = ctxlog.With(ctx, zap.String("guild_Id", e.GuildID.String()))

	ctxlog.Info(ctx, "a user has left a server", zap.String("username", e.User.Username), zap.String("user_id", e.User.ID.String()))

	if e.GuildID == scroteGuildID {
		b.onScroteDeparture(ctx, e)
		return
	}
}
