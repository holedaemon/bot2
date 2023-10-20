package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/zikaeroh/ctxlog"
)

func (b *Bot) cmdTimeIn(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	tz := data.Options.Find("timezone").String()

	if tz == "" {
		id := data.Event.SenderID()
		if id == 0 {
			ctxlog.Error(ctx, "sender id is 0")
			return respondError("An unexpected error has occurred, oops!")
		}

		p, err := modelsx.FetchUserProfile(ctx, b.db, id.String())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return respondError("You gotta give me a timezone to work with")
			}

			ctxlog.Error(ctx, "error fetching user profile")
			return dbError
		}

		if !p.Timezone.Valid {
			return respondError("You gotta give me a timezone to work with")
		}

		tz = p.Timezone.String
	}

	if !validTimezone(tz) {
		return respondError("Timezone must be in **Area/Location** format e.g. **America/Phoenix**")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return respondError("Unable to parse given timezone. Are you sure it's real?")
	}

	now := time.Now().In(loc)
	return respondf(
		"It is <t:%d> in `%s`",
		now.Unix(),
		tz,
	)
}

var formatMap = map[string]string{
	"t": "<t:%d:t>",
	"T": "<t:%d:T>",
	"d": "<t:%d:d>",
	"D": "<t:%d:D>",
	"f": "<t:%d:f>",
	"F": "<t:%d:F>",
	"R": "<t:%d:R>",
}

func (b *Bot) cmdTimeStamp(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	dateTime := data.Options.Find("date-time").String()
	tz := data.Options.Find("timezone").String()
	format := data.Options.Find("format").String()

	switch format {
	case "t", "T", "d", "D", "f", "F", "R":
	case "":
		format = "f"
	default:
		return respondError("Invalid format, valid options are: t, T, d, D, f, F, R. See <https://discord.com/developers/docs/reference#message-formatting-timestamp-styles> for more info")
	}

	if tz == "" {
		id := data.Event.SenderID()
		if id == 0 {
			ctxlog.Error(ctx, "sender id is 0")
			return respondError("An unexpected error has occurred, oops!")
		}

		p, err := modelsx.FetchUserProfile(ctx, b.db, id.String())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return respondError("You gotta give me a timezone to work with")
			}

			ctxlog.Error(ctx, "error fetching user profile")
			return dbError
		}

		if !p.Timezone.Valid {
			return respondError("You gotta give me a timezone to work with")
		}

		tz = p.Timezone.String
	}

	if !validTimezone(tz) {
		return respondError("Timezone must be in **Area/Location** format e.g. **America/Phoenix**")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return respondError("Unable to parse given timezone. Are you sure it's real?")
	}

	t, err := time.ParseInLocation("01/02/2006 15:04", dateTime, loc)
	if err != nil {
		return respondError("Unable to parse date and time. Make sure you're providing it in `MM/DD/YYYY HH:MM` format. e.g. 01/02/2006 15:04")
	}

	formatStr := formatMap[format]
	formatStr = fmt.Sprintf(formatStr, t.Unix())
	return respondf("`%s`", formatStr)
}
