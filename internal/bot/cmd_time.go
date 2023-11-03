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
	"go.uber.org/zap"
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

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return respondError("Unable to parse given timezone. Are you sure it's real?")
	}

	now := time.Now().In(loc)
	return respondf(
		"It is %s in `%s`",
		now.Format("January 2, 2006 03:04 PM"),
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
	t := data.Options.Find("time").String()
	tz := data.Options.Find("timezone").String()
	format := data.Options.Find("format").String()

	switch format {
	case "t", "T", "d", "D", "f", "F", "R":
	case "":
		format = "t"
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

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return respondError("Unable to parse given timezone. Are you sure it's real?")
	}

	now := time.Now().In(loc)
	res, err := b.when.Parse(t, now)
	if err != nil {
		ctxlog.Error(ctx, "error parsing time", zap.Error(err))
		return respondError("Unable to parse given time...")
	}

	if res == nil {
		return respondError("Unable to match given time...")
	}

	formatStr := formatMap[format]
	formatStr = fmt.Sprintf(formatStr, res.Time.Unix())
	return respondSilentf("`%s`", formatStr)
}
