package topster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Error is an error sent by a microgopster instance.
type Error struct {
	HTTPStatus int    `json:"-"`
	Message    string `json:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	var sb strings.Builder

	if e.HTTPStatus > 0 {
		sb.WriteString(
			fmt.Sprintf("HTTP %d", e.HTTPStatus),
		)
	}

	if e.Message != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}

		sb.WriteString(e.Message)
	}

	return sb.String()
}

func makeError(res *http.Response) error {
	var te *Error
	if err := json.NewDecoder(res.Body).Decode(&te); err != nil {
		return err
	}

	te.HTTPStatus = res.StatusCode
	return te
}
