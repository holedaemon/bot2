package bot

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/holedaemon/bot2/internal/api/jerkcity"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

var jerkcityEmbedColor = discord.Color(2123412)

func makeEpisodeEmbed(e *jerkcity.Episode) discord.Embed {
	return discord.Embed{
		Title:     e.Title,
		Color:     jerkcityEmbedColor,
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
		ctxlog.Error(ctx, "error reading value as int", zap.Error(err))
		return respondError("Something about the episode number you gave is wrong. Fix it.")
	}

	episode, err := b.jc.FetchEpisode(ctx, int(number))
	if err != nil {
		ctxlog.Error(ctx, "error fetching episode", zap.Error(err), zap.Int64("number", number))
		return respondError("Sorry, something went wrong and I couldn't get that episode")
	}

	return respondEmbeds(makeEpisodeEmbed(episode))
}

func (b *Bot) cmdJerkcityLatest(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	high, err := b.jc.FetchMeta(ctx)
	if err != nil {
		ctxlog.Error(ctx, "error fetching meta payload", zap.Error(err))
		return respondError("Couldn't get the latest episode number. Sorry?")
	}

	episode, err := b.jc.FetchEpisode(ctx, high)
	if err != nil {
		ctxlog.Error(ctx, "error fetching latest episode", zap.Error(err))
		return respondError("Couldn't get the latest episode. Sorry I guess")
	}

	return respondEmbeds(makeEpisodeEmbed(episode))
}

func (b *Bot) cmdJerkcityRandom(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	episode, err := b.jc.FetchQuote(ctx)
	if err != nil {
		ctxlog.Error(ctx, "error fetching episode", zap.Error(err))
		return respondError("Sorry, something went wrong and I wasn't able to get an episode")
	}

	embed := makeEpisodeEmbed(episode)
	embed.Fields = append(embed.Fields, discord.EmbedField{
		Name:  "Quote",
		Value: episode.Quote,
	})

	return respondEmbeds(embed)
}

func (b *Bot) cmdJerkcitySearch(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	so := data.Options.Find("query")
	search := so.String()

	results, err := b.jc.FetchSearch(ctx, search)
	if err != nil {
		ctxlog.Error(ctx, "error searching jerkcity API", zap.Error(err))
		return respondError("Something happened and I wasn't able to perform the search")
	}

	embed := discord.Embed{
		Title: "Search Results",
		Color: jerkcityEmbedColor,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Took %s", time.Second*time.Duration(results.Search.Runtime)),
		},
	}

	if len(results.Episodes) == 0 {
		embed.Description = "No results..."
	} else {
		var sb strings.Builder

		for i := 0; i < 10; i++ {
			e := results.Episodes[i]
			sb.WriteString(
				fmt.Sprintf("[%d - %s](https://bonequest.com/%d)\n", e.Episode, e.Title, e.Episode),
			)
		}

		more := len(results.Episodes[9:])
		sb.WriteString(
			fmt.Sprintf(
				"[and %d more](https://bonequest.com/search/?q=%s)",
				more,
				url.QueryEscape(search),
			),
		)

		embed.Description = sb.String()
	}

	if results.Search.Sums != nil {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name: "Sums",
			Value: fmt.Sprintf(
				"Dates: %d\nEpisodes: %d\nTags: %d\nTitles: %d\nWords: %d",
				results.Search.Sums.Dates,
				results.Search.Sums.Episodes,
				results.Search.Sums.Tags,
				results.Search.Sums.Titles,
				results.Search.Sums.Words,
			),
		})
	}

	return respondEmbeds(embed)
}
