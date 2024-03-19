package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/holedaemon/bot2/internal/db/modelsx"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

var egoraptorNames = []string{"egoraptor", "arin hanson", "arin"}
var egoraptorThings = []string{"cunnilingus", "pussy", "cunt", "vagina"}

var azLoc *time.Location

func init() {
	loc, err := time.LoadLocation("America/Phoenix")
	if err != nil {
		panic(err)
	}

	azLoc = loc
}

func (b *Bot) onScroteMessage(ctx context.Context, m *gateway.MessageCreateEvent) {
	namesCheck := wordInContent(m.Content, egoraptorNames)
	thingsCheck := wordInContent(m.Content, egoraptorThings)

	if namesCheck && thingsCheck {
		data, err := modelsx.FetchSetting(ctx, b.db, m.GuildID)
		if err != nil {
			ctxlog.Error(ctx, "error querying egoraptor mention", zap.Error(err))
			return
		}

		timestamp := data.LastTimestamp
		since := time.Since(timestamp)
		dur := fmtDur(since)
		var content string

		if data.TimeoutOnMention {
			content = fmt.Sprintf("It has been %s since the last mention of egoraptor eating pussy. You have 10 seconds to speak your piece. Please turn and face the wall.", dur)
		} else {
			content = fmt.Sprintf("It has been %s since the last mention of egoraptor eating pussy", dur)
		}

		image := fakePNG("egopussy")
		err = b.sendImage(ctx, m.ChannelID, content, image)
		if err != nil {
			ctxlog.Error(ctx, "error sending image", zap.Error(err))
			return
		}

		data.Count++
		data.LastTimestamp = time.Now().In(azLoc)
		data.LastUser = m.Author.ID.String()

		err = modelsx.UpsertSetting(ctx, b.db, data)
		if err != nil {
			ctxlog.Error(ctx, "error upserting egoraptor mention", zap.Error(err))
		}

		if data.TimeoutOnMention {
			time.Sleep(time.Second * 10)
			duration := time.Duration(data.TimeoutLength) * time.Second
			t := time.Now().Add(duration)
			ts := discord.NewTimestamp(t)

			err = b.state.ModifyMember(m.GuildID, m.Author.ID, api.ModifyMemberData{
				CommunicationDisabledUntil: &ts,
			})
			if err != nil {
				ctxlog.Error(ctx, "error timing user out", zap.Error(err))
			}
		}
	}
}

func (b *Bot) onScroteDeparture(ctx context.Context, e *gateway.GuildMemberRemoveEvent) {
	guildName := "SCROTEGANG"

	guild, err := b.state.Guild(e.GuildID)
	if err != nil {
		ctxlog.Error(ctx, "error retrieving guild from cache", zap.Error(err))
	} else {
		guildName = guild.Name
	}

	msg := fmt.Sprintf("%s (%s) has left %s please advise", e.User.Username, e.User.ID.String(), guildName)

	ch, err := b.state.CreatePrivateChannel(dylanDMID)
	if err != nil {
		ctxlog.Error(ctx, "error creating private channel", zap.Error(err))
		return
	}

	_, err = b.state.SendMessage(ch.ID, msg)
	if err != nil {
		ctxlog.Error(ctx, "error sending message leave notification", zap.Error(err))
	}
}
