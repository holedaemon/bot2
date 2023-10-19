package topster

import "net/http"

// Option configures a Topster client.
type Option func(*Client)

// WithHTTPClient sets a Topster client's underlying HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(t *Client) {
		t.client = client
	}
}
