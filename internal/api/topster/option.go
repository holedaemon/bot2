package topster

import "net/http"

// Option configures a Topster client.
type Option func(*Topster)

// WithHTTPClient sets a Topster client's underlying HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(t *Topster) {
		t.client = client
	}
}
