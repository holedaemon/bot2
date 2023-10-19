package topster

import (
	"fmt"
	"net/http"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
)

// Topster is an API client that interacts with a companion
// microgopster instance.
type Topster struct {
	addr   string
	client *http.Client
}

// New creates a new Topster client.
func New(addr string, opts ...Option) (*Topster, error) {
	if addr == "" {
		return nil, fmt.Errorf("topster: address is blank")
	}

	t := &Topster{
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
