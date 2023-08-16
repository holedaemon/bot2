package web

import (
	"fmt"
	"net/http"

	"github.com/holedaemon/bot2/internal/web/templates"
)

func (s *Server) basePage(r *http.Request) templates.BasePage {
	return templates.BasePage{
		Username: s.sessionManager.GetString(r.Context(), sessionUsername),
	}
}

func (s *Server) errorPage(w http.ResponseWriter, r *http.Request, code int, msg string, args ...interface{}) {
	if msg != "" {
		msg = fmt.Sprintf(msg, args...)
	}

	w.WriteHeader(code)
	templates.WritePageTemplate(w, &templates.ErrorPage{
		BasePage:  s.basePage(r),
		ErrorHead: http.StatusText(code),
		ErrorText: msg,
	})
}
