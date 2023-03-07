package bot

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

var egoraptorRegexp = regexp.MustCompile(`.*(egoraptor|arin\shanson|arin).*(cunnilingus|pussy|cunt|vagina).*`)

func (b *Bot) onScroteMessage(ctx context.Context, m *gateway.MessageCreateEvent) {
	if egoraptorRegexp.MatchString(m.Content) {
		data, err := b.loadEgoraptorData()
		if err != nil {
			ctxlog.Error(ctx, "error loading egoraptor data", zap.Error(err))
			return
		}

		since := time.Since(data.LastTimestamp)
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

		data.Update(m.Author.ID)
		if err := b.writeEgoraptorData(); err != nil {
			ctxlog.Error(ctx, "error writing egoraptor data", zap.Error(err))
			return
		}
	}
}
