package topster

import (
	"fmt"
	"net/http"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
)

// Client is an API client that interacts with a companion
// microgopster instance.
type Client struct {
	addr   string
	client *http.Client
}

// New creates a new Topster client.
func New(addr string, opts ...Option) (*Client, error) {
	if addr == "" {
		return nil, fmt.Errorf("topster: address is blank")
	}

	t := &Client{
		addr: addr,
	}

	for _, o := range opts {
		o(t)
	}

	if t.client == nil {
		t.client = httpx.DefaultClient
	}

	return t, nil
}
