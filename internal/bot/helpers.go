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

// Code adapted from https://github.com/hako/durafmt/blob/master/durafmt.go#L87
func fmtDur(d time.Duration) string {
	remaining := int64(d / time.Microsecond)

	var (
		sb strings.Builder

		weeks   int64
		days    int64
		hours   int64
		minutes int64
		seconds int64
	)

	weeks = remaining / (7 * 24 * 3600 * 1000000)
	if weeks > 0 {
		remaining -= weeks * 7 * 24 * 3600 * 1000000
	}

	days = remaining / (24 * 3600 * 1000000)
	if days > 0 {
		remaining -= days * 24 * 3600 * 1000000
	}

	hours = remaining / (3600 * 1000000)
	if hours > 0 {
		remaining -= hours * 3600 * 1000000
	}

	minutes = remaining / (60 * 1000000)
	if minutes > 0 {
		remaining -= minutes * 60 * 1000000
	}

	seconds = remaining / 1000000
	if seconds > 0 {
		remaining -= seconds * 1000000
	}

	type durMap struct {
		name string
		dur  int64
	}

	durs := []durMap{
		{name: "week", dur: weeks},
		{name: "day", dur: days},
		{name: "hour", dur: hours},
		{name: "minute", dur: minutes},
		{name: "second", dur: seconds},
	}

	for i, dm := range durs {
		dur := dm.dur
		if dur <= 0 {
			continue
		}

		if dur == 1 {
			sb.WriteString("1 " + dm.name)
		} else {
			sb.WriteString(
				fmt.Sprintf("%d %ss", dur, dm.name),
			)
		}

		if i != len(durs)-1 {
			sb.WriteString(", ")
		}
	}

	if sb.Len() == 0 {
		return "0 seconds"
	}

	return strings.TrimSuffix(sb.String(), ", ")
}
