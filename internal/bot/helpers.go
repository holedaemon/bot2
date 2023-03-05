package bot

import (
	"fmt"
	"strings"
	"time"

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

func fakePNG(path string) string {
	return fakeCDN + "/" + path + ".png"
}

func fakeGif(path string) string {
	return fakeCDN + "/" + path + ".gif"
}

func fakeJPG(path string) string {
	return fakeCDN + "/" + path + ".jpg"
}

func fmtDur(d time.Duration) string {
	if d.Seconds() == 0 {
		return "0 seconds"
	}

	d = d.Round(time.Second)

	var sb strings.Builder

	days := int(d.Hours() / 24)
	hours := int(d.Hours())
	minutes := int(d.Minutes())
	seconds := int(d.Seconds())

	if days > 0 {
		if days == 1 {
			sb.WriteString("1 day")
		} else {
			sb.WriteString(fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if days > 0 {
			sb.WriteString(", ")
		}

		if hours == 1 {
			sb.WriteString("1 hour")
		} else {
			sb.WriteString(
				fmt.Sprintf(
					"%d hours",
					hours,
				),
			)
		}
	}

	if minutes > 0 {
		if days > 0 || hours > 0 {
			sb.WriteString(", ")
		}

		if minutes == 1 {
			sb.WriteString("1 minute")
		} else {
			sb.WriteString(fmt.Sprintf("%d minutes", minutes))
		}
	}

	if sb.Len() == 0 {
		if seconds > 0 {
			if seconds == 1 {
				sb.WriteString("1 second")
			} else {
				sb.WriteString(fmt.Sprintf("%d seconds", seconds))
			}
		}
	}

	return sb.String()
}
