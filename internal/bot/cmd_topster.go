package bot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

type topsterBody struct {
	User            string  `json:"user"`
	Period          string  `json:"period"`
	Title           string  `json:"title"`
	BackgroundColor string  `json:"background_color"`
	TextColor       string  `json:"text_color"`
	Gap             float64 `json:"gap"`
	ShowNumbers     bool    `json:"show_numbers"`
	ShowTitles      bool    `json:"show_titles"`
}

type topsterResponse struct {
	Image string `json:"image"`
}

type topsterError struct {
	Message string `json:"message"`
}

func (b *Bot) cmdTopster(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	if b.topsterAddr == "" {
		return respondError("The Topster command isn't configured correctly! Contact your administrator or something")
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

	switch period {
	case "overall", "7day", "1month", "3month", "6month", "12month":
	case "":
		period = "overall"
	default:
		return respondError("Period should be one of the following: overall, 7day, 1month, 3month, 6month, 12month")
	}

	opts := &topsterBody{
		User:            lastUser,
		Period:          period,
		Title:           title,
		BackgroundColor: backgroundColor,
		TextColor:       textColor,
		Gap:             20,
		ShowNumbers:     showNumbers,
		ShowTitles:      showTitles,
	}

	var input bytes.Buffer
	if err := json.NewEncoder(&input).Encode(&opts); err != nil {
		return respondError("Error encoding Topster options as JSON. How embarrassing...")
	}

	res, err := http.Post(b.topsterAddr, "application/json", &input)
	if err != nil {
		ctxlog.Error(ctx, "error POSTing to topster addr", zap.Error(err))
		return respondError("Error sending request to Topster!!")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var te *topsterError
		if err := json.NewDecoder(res.Body).Decode(&te); err != nil {
			ctxlog.Error(ctx, "error decoding json body", zap.Error(err))
			return respondError("Oops! Something went wrong. Unable to decode Topster error into something readable...")
		}

		return respondErrorf("Something went wrong, Topster says: %s", te.Message)
	}

	var tr *topsterResponse
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		ctxlog.Error(ctx, "error decoding topster response", zap.Error(err))
		return respondError("Error decoding the image sent by Topster...")
	}

	raw, err := base64.StdEncoding.DecodeString(tr.Image)
	if err != nil {
		ctxlog.Error(ctx, "error decoding base64 into image", zap.Error(err))
		return respondError("Error decoding Topster image from base64...")
	}

	output := bytes.NewBuffer(raw)
	return respondImage("topster.jpg", output)
}
