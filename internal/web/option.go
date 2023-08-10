package web

// Option configures a Server.
type Option func(*Server)

// WithAddr sets Server's address.
func WithAddr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}
