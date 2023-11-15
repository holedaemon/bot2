package bot

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

type command struct {
	Scoped discord.GuildID
	Data   api.CreateCommandData
}

type commandList []command

func (cl commandList) Scoped(to discord.GuildID) commandList {
	newCmds := make(commandList, 0)

	for _, c := range cl {
		newCmds = append(newCmds, command{
			Data:   c.Data,
			Scoped: to,
		})
	}

	return newCmds
}

var commands = commandList{
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "ping",
			Description:    "The bot may have a little ping (as a treat)",
			NoDMPermission: false,
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "info",
			Description:    "Displays info about the bot",
			NoDMPermission: false,
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "help",
			Description:    "Displays info about the bot",
			NoDMPermission: false,
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "is-admin",
			Description:    "Debugging command to see if one is an admin",
			NoDMPermission: false,
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:        "game",
			Description: "Change the bot's game presence",
			Options: discord.CommandOptions{
				discord.NewStringOption("new-game", "The new game to change to", true),
			},
			NoDMPermission: false,
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:        "changelog",
			Description: "Get the bot's changelog from GitHub",
		},
	},
	{
		Scoped: testGuildID,
		Data: api.CreateCommandData{
			Name:                     "panic",
			Description:              "Test the recoverer middleware",
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageGuild),
			NoDMPermission:           true,
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "jerkcity",
			Description:    "Various commands relating to Jerkcity",
			NoDMPermission: false,
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
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:                     "roles",
			Description:              "Moderator interface for vanity roles",
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageRoles),
			NoDMPermission:           true,
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
				discord.NewSubcommandOption(
					"fix",
					"Update the database to include role names for tracked roles without.",
				),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "role",
			Description:    "Interface for vanity roles",
			NoDMPermission: true,
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"add",
					"Add a vanity role to yourself",
					discord.NewRoleOption("role", "The role to add", true),
				),
				discord.NewSubcommandOption(
					"remove",
					"Remove a vanity role from yourself",
					discord.NewRoleOption("role", "The role to remove", true),
				),
				discord.NewSubcommandOption(
					"list",
					"List all the vanity roles available",
				),
			},
		},
	},
	{
		Scoped: scroteGuildID,
		Data: api.CreateCommandData{
			Name:                     "egoraptor",
			Description:              "Moderator commands for configuring the Egoraptor functionality",
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageGuild),
			NoDMPermission:           true,
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"toggle",
					"Toggle timing users out",
					discord.NewBooleanOption("toggled", "Whether to timeout or not", true),
				),
				discord.NewSubcommandOption(
					"settimeout",
					"Set the length of timeouts",
					discord.NewIntegerOption(
						"seconds",
						"Number of seconds to time out for",
						true,
					),
				),
			},
		},
	},
	{
		Scoped: scroteGuildID,
		Data: api.CreateCommandData{
			Name:                     "mettic",
			Description:              "Moderator commands for updating Mettic's role on an interval",
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageGuild),
			NoDMPermission:           true,
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"toggle",
					"Toggle updating Mettic's role",
					discord.NewBooleanOption("toggled", "Whether to update or not", true),
				),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "q",
			Description:    "Query guild quotes",
			NoDMPermission: true,
			Options: discord.CommandOptions{
				discord.NewIntegerOption("index", "Quote index", false),
				discord.NewUserOption("user", "Member to query", false),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "quote",
			Description:    "Interact with guild quotes",
			NoDMPermission: true,
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"delete",
					"Delete a quote",
					discord.NewIntegerOption(
						"index",
						"The index to delete",
						true,
					),
				),
				discord.NewSubcommandOption(
					"list",
					"Get a list of guild quotes",
				),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:                     "settings",
			Description:              "Configure bot settings for the guild",
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionManageGuild),
			NoDMPermission:           true,
			Options: discord.CommandOptions{
				discord.NewSubcommandGroupOption(
					"quotes",
					"Configure quote settings for the guild",
					discord.NewSubcommandOption(
						"toggle",
						"Toggle the use of quotes in the guild",
						discord.NewBooleanOption(
							"toggled",
							"Toggle quotes on or off",
							true,
						),
					),
					discord.NewSubcommandOption(
						"set-min-required",
						"Sets the minimum amount of reactions required to quote",
						discord.NewIntegerOption(
							"minimum",
							"The minimum required to quote a message",
							true,
						),
					),
				),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:        "topster",
			Description: "Generate a 3x3 Topster chart with Last.fm data",
			Options: discord.CommandOptions{
				discord.NewStringOption("lastfm-user", "Your Last.fm username", true),
				discord.NewStringOption("period", "The period of time to use for data. Valid options are overall, 7day, 1month, 3month, 6month, 12month", false),
				discord.NewStringOption("background-color", "A background color as a hex code", false),
				discord.NewStringOption("text-color", "A text color as a hex code", false),
				discord.NewStringOption("title", "Your chart's title", false),
				discord.NewBooleanOption("show-titles", "Show titles on the chart?", false),
				discord.NewBooleanOption("show-numbers", "Show numbers on the chart?", false),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "tags",
			Description:    "Commands for managing user-generated tags; think custom commands!",
			NoDMPermission: true,
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"create",
					"Create a tag",
					discord.NewStringOption("name", "Name of the tag", true),
					discord.NewStringOption("content", "Content of the tag", true),
				),
				discord.NewSubcommandOption(
					"update",
					"Update a tag's content",
					discord.NewStringOption("name", "Name of the tag", true),
					discord.NewStringOption("content", "New content of the tag", true),
				),
				discord.NewSubcommandOption(
					"rename",
					"Rename a tag",
					discord.NewStringOption("name", "Current name of the tag", true),
					discord.NewStringOption("new-name", "New name of the tag", true),
				),
				discord.NewSubcommandOption(
					"delete",
					"Delete a tag",
					discord.NewStringOption("name", "Name of the tag", true),
				),
				discord.NewSubcommandOption(
					"list",
					"Get a link to a list of tags for the guild",
				),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:           "tag",
			Description:    "Query a tag",
			NoDMPermission: true,
			Options: discord.CommandOptions{
				discord.NewStringOption("name", "Name of the tag", true),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:        "profile",
			Description: "Interface for user profiles",
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"init",
					"Create a profile for yourself",
				),
				discord.NewSubcommandOption(
					"delete",
					"Delete your profile. THIS IS IRREVERSIBLE.",
					discord.NewBooleanOption("are-you-sure", "Confirmation you want to delete your profile", true),
				),
				discord.NewSubcommandOption(
					"get",
					"Get a link to your profile on the site",
				),
				discord.NewSubcommandOption(
					"set-timezone",
					"Set your timezone",
					discord.NewStringOption("timezone", "Your timezone", true),
				),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:        "time",
			Description: "Interface for time-related commands",
			Options: discord.CommandOptions{
				discord.NewSubcommandOption(
					"in",
					"Get the time in a specific timezone or your profile's timezone",
					discord.NewStringOption("timezone", "A valid timezone", false),
				),
				discord.NewSubcommandOption(
					"stamp",
					"Get a Discord timestamp for a given date and time",
					discord.NewStringOption("time", "A time in natural language", true),
					discord.NewStringOption("timezone", "A valid timezone", false),
					discord.NewStringOption("format", "The format of the Discord timestamp", false),
				),
			},
		},
	},
	{
		Scoped: 0,
		Data: api.CreateCommandData{
			Name:        "feedback",
			Description: "Submit bot feedback, ideas, or tell the owner they suck",
			Options: discord.CommandOptions{
				discord.NewStringOption("feedback", "Your feedback", true),
			},
		},
	},
}
