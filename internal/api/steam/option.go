package steam

import "net/http"

// Option conigures a Steam client.
type Option func(*Client)

// WithHTTPClient sets a Steam client's underlying HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}
