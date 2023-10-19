package topster

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
)

type response struct {
	Image string `json:"image"`
}

// ErrChartOption is returned when an invalid chart option is set.
var ErrChartOption = errors.New("topster: invalid chart option")

// CreateChart creates a chart using the configured microgopster instance
// and returns it as a Base64 encoded string.
func (t *Topster) CreateChart(ctx context.Context, user string, opts ...ChartOption) (string, error) {
	if user == "" {
		return "", fmt.Errorf("%w: user is blank", ErrChartOption)
	}

	co := &chartOptions{
		User: user,
	}

	for _, o := range opts {
		o(co)
	}

	switch co.Period {
	case "overall", "7day", "1month", "3month", "6month", "12month":
	case "":
		co.Period = "7day"
	default:
		return "", fmt.Errorf("%w: period must be one of the following: overall, 7day, 1month, 3month, 6month, 12month", ErrChartOption)
	}

	var in bytes.Buffer
	if err := json.NewEncoder(&in).Encode(&co); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.addr, &in)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", httpx.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	res, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if !httpx.IsOK(res.StatusCode) {
		return "", makeError(res)
	}

	var tr *response
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		return "", err
	}

	return tr.Image, nil
}
