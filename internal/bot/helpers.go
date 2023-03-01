package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var (
	dbError = &api.InteractionResponseData{
		Content: option.NewNullableString("A database error has occurred xD"),
		Flags:   discord.EphemeralMessage,
	}
)

func respondError(msg string) *api.InteractionResponseData {
	if msg == "" {
		panic("bot: blank string given to respondError")
	}

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

	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
		Flags:   discord.EphemeralMessage,
	}
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

func roleInSlice(id discord.RoleID, list []discord.RoleID) bool {
	for _, r := range list {
		if r == id {
			return true
		}
	}

	return false
}
