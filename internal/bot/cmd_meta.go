package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func (b *Bot) cmdPing(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content: option.NewNullableString("Who up riding they pig !?"),
	}
}

func (b *Bot) cmdIsAdmin(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.IsAdmin(data.Event.SenderID()) {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("You are an admin"),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString("HA! BIIIIIIIIIIIIIIIIIIIIIIIIIITCH"),
	}
}
