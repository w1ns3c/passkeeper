package grpc

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"passkeeper/internal/transport/grpc/handlers"
	"passkeeper/internal/transport/grpc/interceptors"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

var (
	errHndNotRegistered = fmt.Errorf("can't register some handlers")
	errNotEnoughOptions = fmt.Errorf("not enough options for grpc constructor")
)

type TransportOption func(*TransportGRPC)

type TransportGRPC struct {
	//ctx    context.Context
	addr   *net.TCPAddr
	srvRPC *grpc.Server

	hndCreds      *handlers.BlobsHandler
	hndUsers      *handlers.UsersHandler
	hndChangePass *handlers.UserChangePassHandler

	authInterceptor *interceptors.AuthInterceptor

	users usersUC.UserUsecaseInf
	creds blobsUC.BlobUsecaseInf

	log *zerolog.Logger
}

func NewTransportGRPC(opts ...TransportOption) (srv *TransportGRPC, err error) {
	srv = new(TransportGRPC)

	for _, opt := range opts {
		opt(srv)
	}

	// check users usecase
	if srv.users == nil {
		if srv.log != nil {
			srv.log.Error().
				Err(err).Msg("no users usecase")
		}

		return nil, errNotEnoughOptions
	}

	// check creds usecase
	if srv.creds == nil {
		if srv.log != nil {
			srv.log.Error().
				Err(err).Msg("no creds usecase")
		}

		return nil, errNotEnoughOptions
	}

	// init intercepters
	srv.authInterceptor = interceptors.NewAuthInterceptor(srv.users)

	// init handlers
	srv.hndCreds = handlers.NewBlobsHandler(srv.log, srv.creds)
	srv.hndUsers = handlers.NewUsersHandler(srv.log, srv.users)

	srv.hndChangePass = handlers.NewUserChangePassHandler(srv.log, srv.users)

	srv.srvRPC = grpc.NewServer(
		grpc.ChainUnaryInterceptor(srv.authInterceptor.AuthFunc()),
	)

	return
}

func WithLogger(logger *zerolog.Logger) TransportOption {
	return func(srv *TransportGRPC) {
		srv.log = logger
	}
}

func WithUCusers(users usersUC.UserUsecaseInf) TransportOption {
	return func(srv *TransportGRPC) {
		srv.users = users
	}
}

func WithUCcreds(creds blobsUC.BlobUsecaseInf) TransportOption {
	return func(srv *TransportGRPC) {
		srv.creds = creds
	}
}

func WithAddr(addr string) TransportOption {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Printf("Can't resolv tcp addr")
		return nil
	}

	return func(srv *TransportGRPC) {
		srv.addr = netAddr
	}
}

// Run will start GRPC server
func (srv *TransportGRPC) Run() error {
	if srv.hndCreds == nil || srv.hndUsers == nil || srv.hndChangePass == nil {
		if srv.log != nil {
			srv.log.Error().
				Msgf("grpc handlers is <nil>")
		}
		return errHndNotRegistered
	}

	pb.RegisterUserSvcServer(srv.srvRPC, srv.hndUsers)
	pb.RegisterBlobSvcServer(srv.srvRPC, srv.hndCreds)
	pb.RegisterUserChangePassSvcServer(srv.srvRPC, srv.hndChangePass)

	l, err := net.Listen("tcp4", srv.addr.String())
	if err != nil {
		if srv.log != nil {
			srv.log.Error().
				Err(err).Msgf("can't start grpc server")
		}
		return err
	}

	if srv.log != nil {
		srv.log.Info().Msgf("[+] Server started on: %s", srv.addr.String())
	}

	return srv.srvRPC.Serve(l)
}

func (srv *TransportGRPC) Stop() error {
	srv.srvRPC.GracefulStop()
	srv.log.Info().Msgf("[+] Server stopped")

	return nil
}
