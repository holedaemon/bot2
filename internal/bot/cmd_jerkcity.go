package bot

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/holedaemon/bot2/internal/api/jerkcity"
	"go.uber.org/zap"
)

var embedColor = discord.Color(2123412)

func makeEpisodeEmbed(e *jerkcity.Episode) discord.Embed {
	return discord.Embed{
		Title:     e.Title,
		Color:     embedColor,
		Timestamp: discord.NewTimestamp(e.Time()),
		Image: &discord.EmbedImage{
			URL: "https://bonequest.com/" + e.Image,
		},
		URL: fmt.Sprintf("https://bonequest.com/%d", e.Episode),
	}
}

func (b *Bot) cmdJerkcityEpisode(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	opt := data.Options.Find("number")

	number, err := opt.IntValue()
	if err != nil {
		b.l.Error("error reading value as int", zap.Error(err))
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Something about the episode number you gave is wrong. Fix it."),
		}
	}

	episode, err := b.jc.FetchEpisode(ctx, int(number))
	if err != nil {
		b.l.Error("error fetching episode", zap.Error(err), zap.Int64("number", number))
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Sorry, something went wrong and I couldn't get that episode."),
		}
	}

	return &api.InteractionResponseData{
		Embeds: &[]discord.Embed{
			makeEpisodeEmbed(episode),
		},
	}
}

func (b *Bot) cmdJerkcityRandom(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	episode, err := b.jc.FetchQuote(ctx)
	if err != nil {
		b.l.Error("error fetching episode", zap.Error(err))
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Sorry, something went wrong and I wasn't able to get an episode."),
		}
	}

	embed := makeEpisodeEmbed(episode)
	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:  "Quote",
		Value: episode.Quote,
	})

	return &api.InteractionResponseData{
		Embeds: &[]discord.Embed{
			embed,
		},
	}
}
