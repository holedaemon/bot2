package bot

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

var egoraptorRegexp = regexp.MustCompile(`.*(egoraptor|arin\shanson|arin).*(cunnilingus|pussy|cunt|vagina).*`)

var azLoc *time.Location

func init() {
	loc, err := time.LoadLocation("America/Phoenix")
	if err != nil {
		panic(err)
	}

	azLoc = loc
}

func (b *Bot) onScroteMessage(ctx context.Context, m *gateway.MessageCreateEvent) {
	if egoraptorRegexp.MatchString(m.Content) {
		data, err := modelsx.FetchMention(ctx, b.DB, m.GuildID)
		if err != nil {
			ctxlog.Error(ctx, "error querying egoraptor mention", zap.Error(err))
			return
		}

		timestamp := data.LastTimestamp
		since := time.Since(timestamp)
		dur := fmtDur(since)

		content := fmt.Sprintf("It has been %s since the last mention of egoraptor eating pussy", dur)

		image := fakePNG("egopussy")
		err = b.SendImage(m.ChannelID, content, image)
		if err != nil {
			ctxlog.Error(ctx, "error sending image", zap.Error(err))
			return
		}

		if data.TimeoutOnMention {
			duration := time.Duration(data.TimeoutLength) * time.Second
			t := time.Now().Add(duration)
			ts := discord.NewTimestamp(t)

			err = b.State.ModifyMember(m.GuildID, m.Author.ID, api.ModifyMemberData{
				CommunicationDisabledUntil: &ts,
			})
			if err != nil {
				ctxlog.Error(ctx, "error timing user out", zap.Error(err))
			}
		}

		data.Count++
		data.LastTimestamp = time.Now().In(azLoc)
		data.LastUser = m.Author.ID.String()

		err = modelsx.UpsertMention(ctx, b.DB, data)
		if err != nil {
			ctxlog.Error(ctx, "error upserting egoraptor mention", zap.Error(err))
		}
	}
}
