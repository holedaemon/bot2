package bot

import (
	"context"
	"path"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
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
