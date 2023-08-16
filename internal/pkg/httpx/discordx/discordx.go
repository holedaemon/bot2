// Package discordx implements useful extensions to the arikawa/discord package.
package discordx

import "golang.org/x/oauth2"

// Endpoint is an OAuth2 endpoint for Discord.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}
