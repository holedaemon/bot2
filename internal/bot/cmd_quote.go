package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/volatiletech/null/v9"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

// TODO: better error handling?
func (b *Bot) onMessageReactionAdd(ev *gateway.MessageReactionAddEvent) {
	if !ev.GuildID.IsValid() {
		return
	}

	if ev.Member.User.Bot {
		return
	}

	if !ev.Emoji.IsUnicode() {
		return
	}

	if ev.Emoji.Name != "💬" {
		return
	}

	ctx := context.Background()
	log := b.Logger.With(zap.String("guild_id", ev.GuildID.String()))

	tx, err := b.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Error("error creating transaction", zap.Error(err))
		return
	}

	defer func() {
		if err := tx.Commit(); err != nil {
			log.Error("error committing transaction", zap.Error(err))
			return
		}
	}()

	exists, err := models.Quotes(qm.Where("message_id = ?", ev.MessageID.String())).Exists(ctx, tx)
	if err != nil {
		log.Error("error checking if quote exists", zap.Error(err))
		return
	}

	if exists {
		return
	}

	msg, err := b.State.Message(ev.ChannelID, ev.MessageID)
	if err != nil {
		log.Error("error retrieving quoted message", zap.Error(err))
		return
	}

	var row struct {
		MaxNum null.Int
	}

	err = models.Quotes(
		qm.Where("guild_id = ?", ev.GuildID.String()),
		qm.Select("max("+models.QuoteColumns.Num+") as max_num"),
	).Bind(ctx, tx, &row)
	if err != nil {
		log.Error("error getting latest quote number from database", zap.Error(err))
		return
	}

	nextNum := row.MaxNum.Int + 1

	quote := &models.Quote{
		Quote:          msg.Content,
		GuildID:        ev.GuildID.String(),
		ChannelID:      ev.ChannelID.String(),
		MessageID:      ev.MessageID.String(),
		QuoterID:       ev.UserID.String(),
		QuotedID:       msg.Author.ID.String(),
		QuotedUsername: msg.Author.Tag(),
		Num:            nextNum,
	}

	if err := quote.Insert(ctx, tx, boil.Infer()); err != nil {
		log.Error("error inserting quote", zap.Error(err))
		return
	}

	if err := b.State.React(msg.ChannelID, msg.ID, discord.APIEmoji("💬")); err != nil {
		log.Error("error adding reaction", zap.Error(err))
	}

	response := fmt.Sprintf(
		"New quote added by %s as #%d %s",
		ev.Member.User.Tag(),
		nextNum,
		jumpLink(ev.GuildID, ev.ChannelID, ev.MessageID),
	)

	if _, err := b.State.SendMessage(ev.ChannelID, response); err != nil {
		log.Error("error sending quote message", zap.Error(err))
	}
}

func (b *Bot) onMessageEdit(e *gateway.MessageUpdateEvent) {
	ctx := context.Background()
	log := b.Logger.With(zap.String("guild_id", e.GuildID.String()), zap.String("message_id", e.ID.String()))

	tx, err := b.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Error("error creating transaction", zap.Error(err))
		return
	}

	defer func() {
		if err := tx.Commit(); err != nil {
			log.Error("error committing transaction", zap.Error(err))
		}
	}()

	quote, err := models.Quotes(qm.Where("message_id = ?", e.ID.String())).One(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}

		log.Error("error querying for quote", zap.Error(err))
		return
	}

	quote.Quote = e.Content
	if err := quote.Update(ctx, tx, boil.Infer()); err != nil {
		log.Error("error updating quote", zap.Error(err))
		if err := tx.Rollback(); err != nil {
			log.Error("error rolling back transaction", zap.Error(err))
			return
		}
	}
}

func makeQuoteEmbed(quote *models.Quote) discord.Embed {
	e := discord.Embed{
		Color:       1974050,
		Description: fmt.Sprintf("%s\n\\-\t<@%s> [(Jump)](%s)", quote.Quote, quote.QuotedID, jumpLinkString(quote.GuildID, quote.ChannelID, quote.MessageID)),
		Timestamp:   discord.Timestamp(quote.CreatedAt),
	}

	e.Author = &discord.EmbedAuthor{
		Name: fmt.Sprintf("#%d", quote.Num),
	}
	return e
}

func (b *Bot) cmdQ(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	idx, err := data.Options.Find("index").IntValue()
	if err != nil {
		ctxlog.Error(ctx, "error retrieving quote index", zap.Error(err))
		return respondError("The index provided is malformed")
	}

	if idx > 0 {
		quote, err := models.Quotes(
			qm.Where("guild_id = ?", data.Event.GuildID.String()),
			qm.Where("num = ?", idx),
		).One(ctx, b.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return respond("No quote at the requested index")
			}

			ctxlog.Error(ctx, "error fetching quote", zap.Error(err))
			return respondError("Error fetching quote!")
		}

		embed := makeQuoteEmbed(quote)
		return respondEmbeds(embed)
	}

	user, err := data.Options.Find("user").SnowflakeValue()
	if err != nil {
		ctxlog.Error(ctx, "error retrieving quote mention", zap.Error(err))
		return respondError("The mention you provided is malformed")
	}

	if user.IsValid() {
		quote, err := models.Quotes(
			qm.Where("guild_id = ?", data.Event.GuildID.String()),
			qm.Where("quoted_id = ?", user.String()),
			qm.OrderBy("random()"),
		).One(ctx, b.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return respond("There aren't any quotes by that user. They should try being funny")
			}

			ctxlog.Error(ctx, "error fetching quote", zap.Error(err))
			return dbError
		}

		embed := makeQuoteEmbed(quote)
		return respondEmbeds(embed)
	}

	quote, err := models.Quotes(
		qm.Where("guild_id = ?", data.Event.GuildID.String()),
		qm.OrderBy("random()"),
	).One(ctx, b.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respond("No quotes found. Sad")
		}

		ctxlog.Error(ctx, "error fetching random quote", zap.Error(err))
		return dbError
	}

	embed := makeQuoteEmbed(quote)
	return respondEmbeds(embed)
}

func (b *Bot) cmdQuoteDelete(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	index, err := data.Options.Find("index").IntValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing int value", zap.Error(err))
		return respondError("Error parsing index value, oops")
	}

	quote, err := models.Quotes(
		qm.Where("num = ?", index),
		qm.Where("guild_id = ?", data.Event.GuildID.String()),
	).One(ctx, b.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respond("There isn't a quote at that index, doofus")
		}

		ctxlog.Error(ctx, "error querying quote", zap.Error(err))
		return dbError
	}

	if err := quote.Delete(ctx, b.DB); err != nil {
		ctxlog.Error(ctx, "error deleting quote", zap.Error(err))
		return dbError
	}

	return respondf("Quote #%d has been deleted", index)
}
