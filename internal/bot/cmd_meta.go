package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/holedaemon/bot2/internal/version"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) cmdInfo(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	user, err := b.state.Me()
	if err != nil {
		ctxlog.Error(ctx, "error getting user", zap.Error(err))
		return respondError("Sorry, I'm having an identity crisis")
	}

	embed := discord.Embed{
		Title:       fmt.Sprintf("Hi, I'm %s 👋", user.Username),
		Description: "The winning submission of the 2016 Desert Bowl College STEM fair, submitted by James McCormick. James unfortunately passed away in an accident shortly after; this project is left running to honor his memory.",
		Color:       discord.Color(4289797),
		Author: &discord.EmbedAuthor{
			Name: user.Username,
			Icon: user.AvatarURL(),
		},
		Fields: []discord.EmbedField{
			{
				Name:   "Version",
				Value:  version.Version(),
				Inline: true,
			},
			{
				Name:   "Uptime",
				Value:  "∞",
				Inline: true,
			},
			{
				Name:   "Wow",
				Value:  "[Website](https://bot.holedaemon.net)",
				Inline: true,
			},
		},
	}

	return respondEmbeds(embed)
}

func (b *Bot) cmdPing(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return respond("Who up riding they pig !?")
}

func (b *Bot) cmdIsAdmin(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.isAdmin(data.Event.SenderID()) {
		return respond("You are an admin")
	}

	return respond("HA! BIIIIIIIIIIIIIIIIIIIIIIIIIITCH")
}

func (b *Bot) cmdGame(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.lastGameChange.Add(time.Hour).After(time.Now()) && !b.isAdmin(data.Event.SenderID()) {
		return respond("The game can only be changed once an hour")
	}

	newGame := data.Options.Find("new-game").String()
	if newGame == "" {
		return respondError("You gotta gimme something to work with here!!!")
	}

	if len(newGame) > 128 {
		return respondError("Game can only be 128 characters in length")
	}

	if err := b.state.Gateway().Send(ctx, &gateway.UpdatePresenceCommand{
		Activities: []discord.Activity{{
			Name: newGame,
			Type: discord.GameActivity,
		}},
	}); err != nil {
		ctxlog.Error(ctx, "error changing presence", zap.Error(err))

		return respondError("That shit broked")
	}

	b.lastGameChange = time.Now()

	return respond("The game has been changed 👉👌")
}

func (b *Bot) cmdPanic(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if !b.isAdmin(data.Event.SenderID()) {
		return respondError("I don't think so, weather man")
	}

	ctxlog.Error(ctx, "why should we panic when we can just log?")
	return respond("Nothing to see here!!")
}

func (b *Bot) cmdChangelog(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	release, _, err := b.github.Repositories.GetLatestRelease(ctx, repoOwner, repoName)
	if err != nil {
		ctxlog.Error(ctx, "error fetching github release", zap.Error(err))
		return respondError("Unable to fetch GitHub release...")
	}

	var sb strings.Builder

	if release.Name != nil {
		sb.WriteString("# " + release.GetName() + "\n")
	}

	if release.Body != nil {
		sb.WriteString(release.GetBody())
	} else {
		sb.WriteString("Release doesn't have a body...")
	}

	return respond(sb.String())
}

var feedbackEmbedColor = discord.Color(1815851)

func makeFeedbackEmbed(data cmdroute.CommandData) discord.Embed {
	sender := data.Event.Sender()

	embed := discord.Embed{
		Title:     "Feedback Submitted",
		Color:     feedbackEmbedColor,
		Timestamp: discord.NewTimestamp(time.Now()),
	}

	if sender != nil {
		embed.Author = &discord.EmbedAuthor{
			Name: sender.Username,
			Icon: sender.AvatarURL(),
		}
	}

	embed.Fields = append(embed.Fields,
		discord.EmbedField{
			Name:  "Feedback",
			Value: data.Options.Find("feedback").String(),
		},
		discord.EmbedField{
			Name:  "Channel ID",
			Value: data.Event.ChannelID.String(),
		},
	)

	return embed
}

func (b *Bot) cmdFeedback(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.feedbackWebhook == nil {
		return respondSilent("Submitting feedback has been disabled...!")
	}

	if data.Options.Find("feedback").String() == "" {
		return respondError("You gotta provide some feedback to use this command...")
	}

	hookData := webhook.ExecuteData{
		Embeds: []discord.Embed{
			makeFeedbackEmbed(data),
		},
	}

	if err := b.feedbackWebhook.Execute(hookData); err != nil {
		ctxlog.Error(ctx, "error firing feedback hook", zap.Error(err))
		return respondError("There was an error submitting your feedback. Try again later!")
	}

	return respondSilent("Your feedback has been recorded")
}
