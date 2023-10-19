package bot

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/holedaemon/bot2/internal/api/topster"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

func (b *Bot) cmdTopster(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.topster == nil {
		return respondError("The Topster command has been disabled! Contact your administrator if you feel this is a mistake.")
	}

	lastUser := data.Options.Find("lastfm-user").String()
	backgroundColor := data.Options.Find("background-color").String()
	textColor := data.Options.Find("text-color").String()
	title := data.Options.Find("title").String()
	period := data.Options.Find("period").String()

	showTitles, err := data.Options.Find("show-titles").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing option as bool", zap.Error(err))
		return respondError("Error parsing option, oops!!")
	}

	showNumbers, err := data.Options.Find("show-numbers").BoolValue()
	if err != nil {
		ctxlog.Error(ctx, "error parsing option as bool", zap.Error(err))
		return respondError("Error parsing option, oops!!")
	}

	if lastUser == "" {
		return respondError("You must provide a last.fm user!!")
	}

	opts := []topster.ChartOption{
		topster.BackgroundColor(backgroundColor),
		topster.TextColor(textColor),
		topster.Title(title),
		topster.Period(period),
		topster.Gap(20),
	}

	if showTitles {
		opts = append(opts, topster.ShowTitles())
	}

	if showNumbers {
		opts = append(opts, topster.ShowNumbers())
	}

	chart, err := b.topster.CreateChart(ctx, lastUser, opts...)
	if err != nil {
		if errors.Is(err, topster.ErrChartOption) {
			return respondErrorf("Your chart options are wrong: %s", err.Error())
		}

		ctxlog.Error(ctx, "error creating topster chart", zap.Error(err))

		switch err.(type) {
		case *topster.Error:
			return respondErrorf("Rut roh, something went wrong. Topster says: %s", err.Error())
		default:
			return respondError("Oops, an unknown error has occurred. Try again later.")
		}
	}

	raw, err := base64.StdEncoding.DecodeString(chart)
	if err != nil {
		ctxlog.Error(ctx, "error decoding base64 into image", zap.Error(err))
		return respondError("Error decoding Topster image from base64...")
	}

	output := bytes.NewBuffer(raw)
	return respondImage("topster.jpg", output)
}
