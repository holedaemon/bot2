package bot

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type egoraptorData struct {
	Count            int            `json:"count"`
	LastTimestamp    time.Time      `json:"last_timestamp"`
	LastUser         discord.UserID `json:"last_user"`
	TimeoutOnMention bool           `json:"timeout_on_mention"`
	TimeoutLength    int64          `json:"timeout_length"`
	mu               sync.Mutex     `json:"-"`
}

func (d *egoraptorData) Update(sf discord.UserID) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Count++
	d.LastTimestamp = time.Now()
	d.LastUser = sf
}

func (d *egoraptorData) SetTimeout(toggled bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.TimeoutOnMention = toggled
}

func (d *egoraptorData) SetTimeoutLength(secs int64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if secs == 0 {
		return
	}

	d.TimeoutLength = secs
}

func (b *Bot) loadEgoraptorData() (*egoraptorData, error) {
	if b.egoraptorData != nil {
		return b.egoraptorData, nil
	}

	file, err := os.Open("egoraptor.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			b.egoraptorData = &egoraptorData{
				LastTimestamp: time.Now(),
				TimeoutLength: 60,
			}
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
