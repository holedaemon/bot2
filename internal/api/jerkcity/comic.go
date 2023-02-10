package jerkcity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
)

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic("jerkcity: error loading location " + err.Error())
	}
}

// Episode is an episode of Jerkcity.
type Episode struct {
	Day     int      `json:"day"`
	Month   int      `json:"month"`
	Year    int      `json:"year"`
	Title   string   `json:"title"`
	Episode int      `json:"episode"`
	Image   string   `json:"image"`
	Players []string `json:"players"`

	// Only present when calling /quote/random
	Quote string `json:"quote,omitempty"`
}

// Time returns an Episode's release, localized in Pacific Time.
func (e *Episode) Time() time.Time {
	return time.Date(
		e.Year,
		time.Month(e.Month),
		e.Day,
		0, 0, 0, 0,
		loc,
	)
}

type episodeResponse struct {
	Episodes []*Episode `json:"episodes"`
}

// FetchEpisode fetches an episode from Jerkcity's API.
func (c *Client) FetchEpisode(ctx context.Context, number int) (*Episode, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"%s/episodes/%d",
			root,
			number,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", httpx.UserAgent)

	res, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if !httpx.IsOK(res.StatusCode) {
		return nil, fmt.Errorf("%w: %d", httpx.ErrStatus, res.StatusCode)
	}

	var eb *episodeResponse
	if err := json.NewDecoder(res.Body).Decode(&eb); err != nil {
		return nil, err
	}

	if len(eb.Episodes) == 0 {
		return nil, nil
	}

	return eb.Episodes[0], nil
}

func (c *Client) FetchQuote(ctx context.Context) (*Episode, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/quote/random", root),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", httpx.UserAgent)

	res, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if !httpx.IsOK(res.StatusCode) {
		return nil, fmt.Errorf("%w: %d", httpx.ErrStatus, res.StatusCode)
	}

	var eb *episodeResponse
	if err := json.NewDecoder(res.Body).Decode(&eb); err != nil {
		return nil, err
	}

	if len(eb.Episodes) == 0 {
		return nil, nil
	}

	return eb.Episodes[0], nil
}
