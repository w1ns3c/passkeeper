package server

import (
	"fmt"

	"github.com/rs/zerolog"

	"passkeeper/internal/config"
	"passkeeper/internal/logger"
	mygrpc "passkeeper/internal/transport/grpc"
	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

type Server struct {
	addr  string // ex: localhost:8000
	users usersUC.UserUsecaseInf
	creds blobsUC.BlobUsecaseInf

	transport *mygrpc.TransportGRPC

	log *zerolog.Logger
}

var (
	errNoUCusers = fmt.Errorf("no users usecase")
	errNoUCcreds = fmt.Errorf("no creds usecase")
)

func newEmptyServer() *Server {
	return &Server{
		addr: config.DefaultAddr,
		log:  logger.Init(config.Level),
	}
}

type SrvOption func(*Server)

func NewServer(opts ...SrvOption) (srv *Server, err error) {
	srv = newEmptyServer()

	for _, opt := range opts {
		opt(srv)
	}

	if srv.users == nil {
		srv.log.Error().
			Err(errNoUCusers).Send()

		return nil, errNoUCusers
	}
	if srv.creds == nil {
		srv.log.Error().
			Err(errNoUCcreds).Send()

		return nil, errNoUCcreds
	}

	transport, err := mygrpc.NewTransportGRPC(
		mygrpc.WithUCcreds(srv.creds),
		mygrpc.WithUCusers(srv.users),
		mygrpc.WithAddr(srv.addr),
		mygrpc.WithLogger(srv.log),
	)
	if err != nil {
		srv.log.Error().
			Err(err).Send()

		return nil, err
	}

	srv.transport = transport

	return srv, nil
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

func WithUCcreds(uc blobsUC.BlobUsecaseInf) SrvOption {
	return func(srv *Server) {
		srv.creds = uc
	}
}

func WithLogger(lg *zerolog.Logger) SrvOption {
	return func(srv *Server) {
		srv.log = lg
	}
}

func (s Server) Run() error {
	return s.transport.Run()
}

func (s Server) Stop() error {
	s.users = nil
	s.creds = nil
	s.transport = nil
	return nil
}
