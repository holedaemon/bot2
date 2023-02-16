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
	{
		Scoped: testGuildID,
		Data: api.CreateCommandData{
			Name:        "jerkcity",
			Description: "Various commands relating to Jerkcity",
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"latest",
					"Fetch the latest episode",
				),
				discord.NewSubcommandOption(
					"episode",
					"Fetch an episode by its number",
					discord.NewIntegerOption("number", "Number of the episode", true),
				),
				discord.NewSubcommandOption(
					"quote",
					"Get a random episode with a cute little quote :)",
				),
				discord.NewSubcommandOption(
					"search",
					"Search the Jerkcity API for a specific query",
					discord.NewStringOption("query", "Your query", true),
				),
			},
		},
	},
	{
		Scoped: testGuildID,
		Data: api.CreateCommandData{
			Name:        "role",
			Description: "Interface for self-assignable vanity roles (pronouns, etc)",
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"create",
					"Create a new self-assignable role",
					discord.NewStringOption("name", "Name of new role", true),
				),
			},
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageRoles),
		},
	},
}
