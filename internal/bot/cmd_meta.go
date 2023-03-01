package bot

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) cmdPing(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content: option.NewNullableString("Who up riding they pig !?"),
	}
}

func (b *Bot) cmdIsAdmin(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.IsAdmin(data.Event.SenderID()) {
		return respond("You are an admin")
	}

	return respond("HA! BIIIIIIIIIIIIIIIIIIIIIIIIIITCH")
}

func (b *Bot) cmdGame(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.lastGameChange.Add(time.Hour).After(time.Now()) && !b.IsAdmin(data.Event.SenderID()) {
		return respond("The game can only be changed once an hour")
	}

	newGame := data.Options.Find("new-game").String()
	if newGame == "" {
		return respondError("You gotta gimme something to work with here!!!")
	}

	if err := b.State.Gateway().Send(ctx, &gateway.UpdatePresenceCommand{
		Activities: []discord.Activity{{
			Name: newGame,
			Type: discord.GameActivity,
		}},
	}); err != nil {
		ctxlog.Error(ctx, "error changing presence", zap.Error(err))

		return respondError("That shit broked")
	}

	b.lastGameChange = time.Now()

	return respond("The game has been changed. 👉👌")
}
