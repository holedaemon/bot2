package jerkcity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
)

var (
	// ErrEmptyQuery is returned when a search query is empty.
	ErrEmptyQuery = errors.New("jerkcity: empty search query")

	loc *time.Location
)

func init() {
	var err error
	loc, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic("jerkcity: error loading location " + err.Error())
	}
}

// Search represents a request to the /search endpoint.
type Search struct {
	Episodes []*Episode     `json:"episodes"`
	Search   *EpisodeSearch `json:"search"`
}

type EpisodeSearchSums struct {
	Dates    int `json:"dates"`
	Episodes int `json:"episodes"`
	Tags     int `json:"tags"`
	Titles   int `json:"titles"`
	Words    int `json:"words"`
}

type EpisodeSearch struct {
	Query   string             `json:"query"`
	Runtime float64            `json:"runtime"`
	Sums    *EpisodeSearchSums `json:"sums"`
	Version int                `json:"version"`
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

	// Only present when calling /search
	Search *EpisodeSearch `json:"search"`
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

type metaResponse struct {
	Meta *Meta `json:"meta"`
}

type Meta struct {
	High int `json:"high"`
}

func (c *Client) FetchMeta(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		root,
		nil,
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", httpx.UserAgent)

	res, err := c.cli.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if !httpx.IsOK(res.StatusCode) {
		return 0, fmt.Errorf("%w: %d", httpx.ErrStatus, res.StatusCode)
	}

	var m *metaResponse
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return 0, err
	}

	return m.Meta.High, nil
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

func (c *Client) FetchSearch(ctx context.Context, query string) (*Search, error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}

	u := url.Values{}
	u.Add("q", query)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		root+"/search",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = u.Encode()

	res, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if !httpx.IsOK(res.StatusCode) {
		return nil, fmt.Errorf("%w: %d", httpx.ErrStatus, res.StatusCode)
	}

	var s *Search
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, err
	}

	return s, nil
}
