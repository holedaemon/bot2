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
	{
		Scoped: testGuildID,
		Data: api.CreateCommandData{
			Name:        "is-admin",
			Description: "Debugging command to see if one is an admin",
		},
	},
	{
		Scoped: testGuildID,
		Data: api.CreateCommandData{
			Name:        "game",
			Description: "Change the bot's game presence",
			Options: discord.CommandOptions{
				discord.NewStringOption("new-game", "The new game to change to", true),
			},
		},
	},
}
