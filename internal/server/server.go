package server

import (
	"passkeeper/internal/usecase/srv/credentialsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

type Server struct {
	addr    string // ex: localhost:8000
	usersUC usersUC.UserUsecaseInf
	credsUC credentialsUC.CredUsecaseInf
}

func NewServer() *Server {
	return &Server{}
}

func (s Server) Run() error {

	return nil
}
