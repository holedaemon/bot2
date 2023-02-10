package jerkcity

import (
	"net/http"
)

const root = "https://bonequest.com/api/v2"

// Client is an HTTP Client responsible for pulling data from Jerkcity's API.
type Client struct {
	cli *http.Client
}

// New creates a new Client.
func New(opts ...Option) *Client {
	c := &Client{}

	for _, o := range opts {
		o(c)
	}

	if c.cli == nil {
		c.cli = http.DefaultClient
	}

	return c
}
