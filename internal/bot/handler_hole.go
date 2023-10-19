package bot

import (
	"context"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

var regexTimeout = regexp.MustCompile(`(computer|boys), \w+ this \w+`)

func (b *Bot) onHoleMessage(ctx context.Context, m *gateway.MessageCreateEvent) {
	if roleInSlice(holeFortniteRoleID, m.MentionRoleIDs) {
		idx := rand.Intn(4) + 1
		str := strconv.FormatInt(int64(idx), 10)
		image := fakeGif("fortnite-" + str)

		if err := b.SendImage(m.ChannelID, "", image); err != nil {
			ctxlog.Error(ctx, "error sending image to channel", zap.Error(err))
		}

		return
	}

	if regexTimeout.MatchString(m.Content) {
		cache := m.ReferencedMessage
		if cache == nil {
			if err := b.Reply(m.Message, "Who's am I's whackin' 'ere?"); err != nil {
				ctxlog.Error(ctx, "error sending reply", zap.Error(err))
			}
			return
		}

		if cache.Author.ID == m.Author.ID {
			if err := b.SendImage(m.ChannelID, "", fakeJPG("snipes")); err != nil {
				ctxlog.Error(ctx, "error sending image", zap.Error(err))
			}
			return
		}

		til := rand.Intn(60)
		tilDur := time.Duration(til) * time.Second
		dur := time.Now().Add(tilDur)
		ts := discord.NewTimestamp(dur)

		if err := b.state.ModifyMember(m.GuildID, cache.Author.ID, api.ModifyMemberData{
			CommunicationDisabledUntil: &ts,
		}); err != nil {
			er := err.(*httputil.HTTPError)
			if er.Code == 50013 {
				if err := b.Reply(m.Message, "Sorry boss, I's forgots my gun"); err != nil {
					ctxlog.Error(ctx, "error sending message", zap.Error(err))
				}
				return
			}

			ctxlog.Error(ctx, "error timing out user", zap.Error(err))

			if err := b.Reply(m.Message, "Sorry boss, da feds got in da way"); err != nil {
				ctxlog.Error(ctx, "error sending message", zap.Error(err))
			}
			return
		}

		if err := b.Reply(m.Message, "Da jobs done, boss"); err != nil {
			ctxlog.Error(ctx, "error sending message", zap.Error(err))
		}
	}
}
