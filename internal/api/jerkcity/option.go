package jerkcity

import (
	"net/http"
)

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a Client's underlying HTTP client.
func WithHTTPClient(cli *http.Client) Option {
	return func(c *Client) {
		c.cli = cli
	}
}
