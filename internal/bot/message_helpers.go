package bot

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

func (b *Bot) reply(m discord.Message, content string) error {
	if content == "" {
		panic("bot: blank content given to Reply")
	}

	_, err := b.state.SendMessageReply(m.ChannelID, content, m.ID)
	return err
}

func (b *Bot) sendImage(ctx context.Context, chID discord.ChannelID, content string, image string) error {
	cachedImage, err := b.imageCache.Get(ctx, image)
	if err != nil {
		return err
	}

	rawName := path.Base(image)
	files := make([]sendpart.File, 0)
	files = append(files, sendpart.File{
		Name:   rawName,
		Reader: cachedImage,
	})

	_, err = b.state.SendMessageComplex(
		chID,
		api.SendMessageData{
			Content: content,
			Files:   files,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func respondError(msg string) *api.InteractionResponseData {
	if msg == "" {
		panic("bot: blank string given to respondError")
	}

	msg = "❌ " + msg

	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
}

func respondErrorf(msg string, args ...interface{}) *api.InteractionResponseData {
	if msg == "" {
		panic("bot: blank string given to respondErrorf")
	}

	msg = fmt.Sprintf(msg, args...)

	return respondError(msg)
}

func respond(msg string) *api.InteractionResponseData {
	if msg == "" {
		panic("bot: blank string given to respond")
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
	}
}

func respondf(msg string, args ...interface{}) *api.InteractionResponseData {
	if msg == "" {
		panic("bot: blank string given to respondf")
	}

	msg = fmt.Sprintf(msg, args...)

	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
	}
}

func respondEmbeds(embeds ...discord.Embed) *api.InteractionResponseData {
	if len(embeds) == 0 {
		panic("bot: no embeds were given to respondEmbeds")
	}

	return &api.InteractionResponseData{
		Embeds: &embeds,
	}
}

func respondImage(name string, image io.Reader) *api.InteractionResponseData {
	if image == nil {
		panic("bot: nil reader given to respondImage")
	}

	files := make([]sendpart.File, 0)
	files = append(files, sendpart.File{
		Name:   name,
		Reader: image,
	})

	return &api.InteractionResponseData{
		Files: files,
	}
}
