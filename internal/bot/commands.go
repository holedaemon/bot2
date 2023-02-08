package bot

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Command struct {
	Scoped discord.GuildID
	Data   api.CreateCommandData
}

var commands = []Command{
	{
		Scoped: testGuildID,
		Data: api.CreateCommandData{
			Name:        "ping",
			Description: "The bot may have a little ping (as a treat)",
		},
	},
}
