package web

import (
	"context"
	"embed"
	"io/fs"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/holedaemon/bot2/internal/web/templates"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

//go:embed static
var static embed.FS

var assetsDir fs.FS

func init() {
	var err error
	assetsDir, err = fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
}

// Server is the HTTP server responsible for serving BOT/2's website.
type Server struct {
	Addr string
}

// New creates a new Server.
func New(opts ...Option) (*Server, error) {
	srv := &Server{}

	for _, o := range opts {
		o(srv)
	}

	return srv, nil
}

// Run starts a Server.
func (s *Server) Run(ctx context.Context) error {
	r := chi.NewMux()

	r.Use(s.recoverer)

	logger := ctxlog.FromContext(ctx)
	r.Use(requestLogger(logger))

	r.Get("/", s.index)
	r.Get("/about", s.about)

	r.NotFound(s.notFound)

	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(assetsDir))))
	r.Handle("/favicon.ico", http.RedirectHandler("/static/favicon.ico", http.StatusFound))

	srv := &http.Server{
		Addr:        s.Addr,
		Handler:     r,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			ctxlog.Error(ctx, "error shutting down server", zap.Error(err))
			return
		}
	}()

	return srv.ListenAndServe()
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.IndexPage{})
}

func (s *Server) about(w http.ResponseWriter, r *http.Request) {
	templates.WritePageTemplate(w, &templates.AboutPage{})
}
