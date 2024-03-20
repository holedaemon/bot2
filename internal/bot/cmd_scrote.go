package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/db/models"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

const maxTimeoutSeconds = 604800

func (b *Bot) cmdEgoraptorToggle(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	toggled, err := data.Options.Find("toggled").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error converting argument into bool", zap.Error(err))
		return respondError("Error converting argument to boolean")
	}

	ego, err := modelsx.FetchSetting(ctx, b.db, data.Event.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error querying for egoraptor mention", zap.Error(err))
		return dbError
	}

	ego.TimeoutOnMention = toggled

	if err := modelsx.UpsertSetting(ctx, b.db, ego); err != nil {
		ctxlog.Error(ctx, "error upserting egoraptor mention", zap.Error(err))
		return dbError
	}

	if toggled {
		return respond("I will now time users out when they mention egoraptor eating pussy")
	}

	return respond("I will no longer time users out when they mention egoraptor eating pussy")
}

func (b *Bot) cmdEgoraptorSetTimeout(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	seconds, err := data.Options.Find("seconds").IntValue()
	if err != nil {
		ctxlog.Error(ctx, "error converting argument to int64", zap.Error(err))
		return respondError("Error converting argument to int64")
	}

	if seconds <= 0 {
		return respondError("The number of seconds must be a positive number greater than zero")
	}

	if seconds > maxTimeoutSeconds {
		return respondError("The number of seconds cannot exceed a week")
	}

	ego, err := modelsx.FetchSetting(ctx, b.db, data.Event.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error querying for egoraptor mention", zap.Error(err))
		return dbError
	}

	ego.TimeoutLength = int(seconds)

	if err := modelsx.UpsertSetting(ctx, b.db, ego); err != nil {
		ctxlog.Error(ctx, "error upserting egoraptor mention", zap.Error(err))
		return dbError
	}

	if seconds == 1 {
		return respond("Timeout length has been set to 1 second")
	}

	return respondf("Timeout length has been set to %d seconds", seconds)
}

func (b *Bot) cmdMetticToggle(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	toggled, err := data.Options.Find("toggled").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing value as bool", zap.Error(err))
		return respondError("Unable to parse value as a boolean!!")
	}

	updater, err := modelsx.FetchRoleUpdater(ctx, b.db, scroteGuildID.String())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctxlog.Error(ctx, "error fetching rolre updater", zap.Error(err))
			return dbError
		}

		updater = &models.RoleUpdater{
			GuildID:       scroteGuildID.String(),
			LastTimestamp: time.Now(),
		}
	}

	updater.DoUpdates = toggled
	if err := modelsx.UpsertRoleUpdater(ctx, b.db, updater); err != nil {
		ctxlog.Error(ctx, "error upserting role updater", zap.Error(err))
		return dbError
	}

	if toggled {
		return respond("I will now occasionally update Mettic's role")
	}

	return respond("I will no longer update Mettic's role")
}

var leaderboardMap = map[int]string{
	0: "🥇",
	1: "🥈",
	2: "🥉",
}

var leaderboardEmbedColor = discord.Color(14360865)

func (b *Bot) cmdLeaderboard(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	mentions, err := modelsx.FetchTopThreeMentions(ctx, b.db, data.Event.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error fetching mentions from database", zap.Error(err))
		return dbError
	}

	if len(mentions) == 0 {
		return respond("There aren't any users to shame...")
	}

	e := discord.Embed{
		Title: "Egoraptor Leaderboard",
		Color: leaderboardEmbedColor,
	}

	fields := make([]discord.EmbedField, 0, 3)

	for i, m := range mentions {
		sf, err := discord.ParseSnowflake(m.UserID)
		if err != nil {
			ctxlog.Error(ctx, "error parsing snowflake", zap.Error(err))
			return respondError("Error building leaderboard...")
		}

		nick, err := b.state.MemberDisplayName(data.Event.GuildID, discord.UserID(sf))
		if err != nil {
			ctxlog.Error(ctx, "error fetching user nickname", zap.Error(err))
			return respondError("Error fetching user nickname...")
		}

		name := fmt.Sprintf("%s **%s**", leaderboardMap[i], nick)
		var value string
		if m.Count == 1 {
			value = "**1** mention"
		} else {
			value = fmt.Sprintf("**%d** mentions", m.Count)
		}

		fields = append(fields, discord.EmbedField{
			Name:  name,
			Value: value,
		})
	}

	e.Fields = fields
	return respondEmbeds(e)
}
