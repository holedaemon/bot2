package web

import (
	"database/sql"

	"golang.org/x/oauth2"
)

// Option configures a Server.
type Option func(*Server)

// WithDebug toggles debug mode on a Server.
func WithDebug(debug bool) Option {
	return func(s *Server) {
		s.Debug = debug
	}
}

// WithAddr sets Server's address.
func WithAddr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

// WithDB sets a Server's DB connection.
func WithDB(db *sql.DB) Option {
	return func(s *Server) {
		s.DB = db
	}
}

// WithOAuth2 sets a Server's OAuth2 config.
func WithOAuth2(oa *oauth2.Config) Option {
	return func(s *Server) {
		s.OAuth2 = oa
	}
}

// WithAdmins sets a Server's admin accounts.
func WithAdmins(admins map[string]string) Option {
	return func(s *Server) {
		s.Admins = admins
	}
}
