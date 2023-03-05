package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

type egoraptorData struct {
	Count         int            `json:"count"`
	LastTimestamp time.Time      `json:"last_timestamp"`
	LastUser      discord.UserID `json:"last_user"`
	mu            sync.Mutex     `json:"-"`
}

func (d *egoraptorData) Update(sf discord.UserID) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Count++
	d.LastTimestamp = time.Now()
	d.LastUser = sf
}

func (b *Bot) loadEgoraptorData() (*egoraptorData, error) {
	if b.egoraptorData != nil {
		return b.egoraptorData, nil
	}

	file, err := os.Open("egoraptor.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			b.egoraptorData = &egoraptorData{}
			return b.egoraptorData, nil
		}

		return nil, err
	}

	defer file.Close()

	var d *egoraptorData
	if err := json.NewDecoder(file).Decode(&d); err != nil {
		return nil, err
	}

	b.egoraptorData = d
	return d, nil
}

func (b *Bot) writeEgoraptorData() error {
	data, err := b.loadEgoraptorData()
	if err != nil {
		return err
	}

	file, err := os.Create("egoraptor.json")
	if err != nil {
		return err
	}

	defer file.Close()

	if err := json.NewEncoder(file).Encode(data); err != nil {
		return err
	}

	return nil
}

var egoraptorRegexp = regexp.MustCompile(`.*(egoraptor|arin\shanson|arin).*(cunnilingus|pussy|cunt|vagina).*`)

func (b *Bot) onScroteMessage(ctx context.Context, m *gateway.MessageCreateEvent) {
	if egoraptorRegexp.MatchString(m.Content) {
		data, err := b.loadEgoraptorData()
		if err != nil {
			ctxlog.Error(ctx, "error loading egoraptor data", zap.Error(err))
			return
		}

		image := b.imageCache.Get("egopussy")
		if image == nil {
			url := fakePNG("egopussy")
			err := b.imageCache.Download(url)
			if err != nil {
				ctxlog.Error(ctx, "error downloading image", zap.Error(err), zap.String("url", url))
				return
			}

			image = b.imageCache.Get("egopussy")
		}

		since := time.Since(data.LastTimestamp)
		dur := fmtDur(since)

		content := fmt.Sprintf("It has been %s since the last mention of egoraptor eating pussy", dur)

		files := make([]sendpart.File, 0)
		file := sendpart.File{
			Name:   "egopussy.png",
			Reader: image,
		}
		files = append(files, file)

		_, err = b.State.SendMessageComplex(m.ChannelID, api.SendMessageData{
			Content: content,
			Files:   files,
		})
		if err != nil {
			ctxlog.Error(ctx, "error sending message", zap.Error(err))
			return
		}

		data.Update(m.Author.ID)
		if err := b.writeEgoraptorData(); err != nil {
			ctxlog.Error(ctx, "error writing egoraptor data", zap.Error(err))
			return
		}
	}
}
