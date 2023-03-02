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
			Name:                     "panic",
			Description:              "Test the recoverer middleware",
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageGuild),
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
			Name:        "roles",
			Description: "Moderator interface for vanity roles",
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"create",
					"Create a new vanity role",
					discord.NewStringOption("name", "Name of new role", true),
					discord.NewStringOption("color", "Color of the new role", false),
					discord.NewBooleanOption("hoisted", "Whether the role is hoisted or not", false),
				),
				discord.NewSubcommandOption(
					"delete",
					"Delete a vanity role",
					discord.NewRoleOption("role", "The role to delete", true),
				),
				discord.NewSubcommandOption(
					"rename",
					"Rename a vanity role",
					discord.NewRoleOption("role", "The role to rename", true),
					discord.NewStringOption("new-name", "The new name", true),
				),
				discord.NewSubcommandOption(
					"setcolor",
					"Change a role's color",
					discord.NewRoleOption("role", "The role to alter", true),
					discord.NewStringOption("color", "The new color", true),
				),
				discord.NewSubcommandOption(
					"import",
					"Import existing roles",
					discord.NewRoleOption("first", "First role to import", true),
					discord.NewRoleOption("second", "Second role to import", false),
					discord.NewRoleOption("third", "Third role to import", false),
					discord.NewRoleOption("fourth", "Fourth role to import", false),
					discord.NewRoleOption("fifth", "Fifth role to import", false),
				),
				discord.NewSubcommandOption(
					"relinquish",
					"Relinquish roles from the database",
					discord.NewRoleOption("first", "First role to import", true),
					discord.NewRoleOption("second", "Second role to import", false),
					discord.NewRoleOption("third", "Third role to import", false),
					discord.NewRoleOption("fourth", "Fourth role to import", false),
					discord.NewRoleOption("fifth", "Fifth role to import", false),
				),
			},
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageRoles),
		},
	},
	{
		Scoped: testGuildID,
		Data: api.CreateCommandData{
			Name:        "role",
			Description: "Interface for vanity roles",
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"add",
					"Add a vanity role to yourself",
					discord.NewRoleOption("role", "The role to add", true),
				),
				discord.NewSubcommandOption(
					"remove",
					"Remove a vanity role from yourself",
					discord.NewStringOption("role", "The role to remove", true),
				),
				discord.NewSubcommandOption(
					"list",
					"List all the vanity roles available",
				),
			},
		},
	},
}
