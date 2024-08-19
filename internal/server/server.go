package server

import (
	"passkeeper/internal/config"
	"passkeeper/internal/usecase/srv/credentialsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

type Server struct {
	addr  string // ex: localhost:8000
	users usersUC.UserUsecaseInf
	creds credentialsUC.CredUsecaseInf
}

func newEmptyServer() *Server {
	return &Server{
		addr: config.DefaultAddr,
	}
}

type SrvOption func(*Server)

func NewServer(opts ...SrvOption) *Server {
	srv := newEmptyServer()

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func WithAddr(addr string) SrvOption {
	return func(srv *Server) {
		srv.addr = addr
	}
}

func WithUCusers(uc usersUC.UserUsecaseInf) SrvOption {
	return func(srv *Server) {
		srv.users = uc
	}
}

func WithUCcreds(uc credentialsUC.CredUsecaseInf) SrvOption {
	return func(srv *Server) {
		srv.creds = uc
	}
}

func (s Server) Run() error {

	return nil
}
