package httpx

import (
	"errors"
	"net/http"
	"time"
)

var ErrStatus = errors.New("non-OK status code")

// UserAgent is Bot2's default user-agent.
const UserAgent = "Bot2"

// DefaultClient is an HTTP Client with sane defaults.
var DefaultClient = &http.Client{
	Timeout: time.Second * 10,
}

// IsOK tests if a status code is within the 2xx range.
func IsOK(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}
