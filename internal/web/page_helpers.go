package web

import (
	"fmt"
	"net/http"

	"github.com/holedaemon/bot2/internal/web/templates"
)

func (s *Server) errorPage(w http.ResponseWriter, code int, head, text string) {
	if head == "" {
		head = http.StatusText(code)
	}

	w.WriteHeader(code)
	templates.WritePageTemplate(w, &templates.ErrorPage{
		ErrorHead: head,
		ErrorText: text,
	})
}

func (s *Server) notFound(w http.ResponseWriter, _ *http.Request) {
	s.errorPage(w, http.StatusNotFound, "", "Whatever you're looking for ain't here")
}

func (s *Server) internalError(w http.ResponseWriter, msg string, args ...interface{}) {
	if msg != "" {
		msg = fmt.Sprintf(msg, args...)
	}

	s.errorPage(w, http.StatusInternalServerError, "", msg)
}
