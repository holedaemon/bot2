package steam

import (
	"fmt"
	"net/http"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
)

const (
	root = "https://api.steampowered.com"
)

// Client is an API client that interacts with the Steam API.
type Client struct {
	key    string
	client *http.Client
}

// New creates a new Steam client.
func New(key string, opts ...Option) (*Client, error) {
	if key == "" {
		return nil, fmt.Errorf("steam: key is blank")
	}

	c := &Client{
		key: key,
	}
	for _, o := range opts {
		o(c)
	}

	if c.client == nil {
		c.client = httpx.DefaultClient
	}

	return c, nil
}
