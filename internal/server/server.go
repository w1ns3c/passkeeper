package server

import (
	"github.com/rs/zerolog"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/logger"
	"passkeeper/internal/entities/myerrors"
	mygrpc "passkeeper/internal/transport/grpc"
	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

// Server consist of all server' parts: transport, storage, usecase
type Server struct {
	addr  string // ex: localhost:8000
	users usersUC.UserUsecaseInf
	blobs blobsUC.BlobUsecaseInf

	transport *mygrpc.TransportGRPC

	log *zerolog.Logger
}

// newEmptyServer is an empty constructor for Server
func newEmptyServer() *Server {
	return &Server{
		addr: config.DefaultAddr,
		log:  logger.Init(config.Level),
	}
}

// SrvOption new type for Server constructor
type SrvOption func(*Server)

// NewServer is a constructor for Server with SrvOption
func NewServer(opts ...SrvOption) (srv *Server, err error) {
	srv = newEmptyServer()

	for _, opt := range opts {
		opt(srv)
	}

	if srv.users == nil {
		srv.log.Error().
			Err(myerrors.ErrNoUCusers).Send()

		return nil, myerrors.ErrNoUCusers
	}
	if srv.blobs == nil {
		srv.log.Error().
			Err(myerrors.ErrNoUCcreds).Send()

		return nil, myerrors.ErrNoUCcreds
	}

	transport, err := mygrpc.NewTransportGRPC(
		mygrpc.WithUCcreds(srv.blobs),
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

// WithAddr setup add for Server
func WithAddr(addr string) SrvOption {
	return func(srv *Server) {
		srv.addr = addr
	}
}

// WithUCusers setup user usecase for Server
func WithUCusers(uc usersUC.UserUsecaseInf) SrvOption {
	return func(srv *Server) {
		srv.users = uc
	}
}

// WithUCblobs setup blob usecase for Server
func WithUCblobs(uc blobsUC.BlobUsecaseInf) SrvOption {
	return func(srv *Server) {
		srv.blobs = uc
	}
}

// WithLogger setup add for Server
func WithLogger(lg *zerolog.Logger) SrvOption {
	return func(srv *Server) {
		srv.log = lg
	}
}

// Run will start Server
func (s *Server) Run() error {
	return s.transport.Run()
}

// Stop will stop Server
func (s *Server) Stop() error {
	s.users = nil
	s.blobs = nil
	s.transport = nil
	return nil
}
