package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
)

type ownedGameResponse struct {
	Response *OwnedGames `json:"response"`
}

// OwnedGame represents a game a Steam user owns.
type OwnedGame struct {
	AppID           int `json:"appid"`
	Playtime2Weeks  int `json:"playtime_2weeks"`
	PlaytimeForever int `json:"playtime_forever"`
}

// OwnedGames represents the list of games a Steam user owns.
type OwnedGames struct {
	GameCount int          `json:"game_count"`
	Games     []*OwnedGame `json:"games"`
}

// GetOwnedGames fetches a Steam user's owned games.
func (c *Client) GetOwnedGames(ctx context.Context, id string) (*OwnedGames, error) {
	u := fmt.Sprintf("%s/IPlayerService/GetOwnedGames/v0001/", root)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("key", c.key)
	q.Set("steamid", id)
	q.Set("format", "json")

	req.URL.RawQuery = q.Encode()
	req.Header.Set("User-Agent", httpx.UserAgent)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if !httpx.IsOK(res.StatusCode) {
		return nil, fmt.Errorf("%w: HTTP %d (%s)", httpx.ErrStatus, res.StatusCode, http.StatusText(res.StatusCode))
	}

	var og *ownedGameResponse
	if err := json.NewDecoder(res.Body).Decode(&og); err != nil {
		return nil, err
	}

	return og.Response, nil
}
