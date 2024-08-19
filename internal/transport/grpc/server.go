package grpc

import (
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"passkeeper/internal/transport/grpc/handlers"
	pb "passkeeper/internal/transport/grpc/protofiles/proto"
	"passkeeper/internal/usecase/srv/credentialsUC"
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

	hndCreds      *handlers.CredsHandler
	hndUsers      *handlers.UsersHandler
	hndChangePass *handlers.UserChangePassHandler

	users usersUC.UserUsecaseInf
	creds credentialsUC.CredUsecaseInf

	log *zerolog.Logger
}

func NewTransportGRPC(opts ...TransportOption) (srv *TransportGRPC, err error) {
	srv = new(TransportGRPC)

	for _, opt := range opts {
		opt(srv)
	}

	// check users usecase
	if srv.users != nil {
		if srv.log != nil {
			srv.log.Error().
				Err(err).Msg("no users usecase")
		}

		return nil, errNotEnoughOptions
	}

	// check creds usecase
	if srv.creds != nil {
		if srv.log != nil {
			srv.log.Error().
				Err(err).Msg("no creds usecase")
		}

		return nil, errNotEnoughOptions
	}

	// register handlers
	srv.hndCreds = handlers.NewCredsHandler(srv.log, srv.creds)
	srv.hndUsers = handlers.NewUsersHandler(srv.log, srv.users)
	srv.hndChangePass = handlers.NewUserChangePassHandler(srv.log, srv.users)

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

func WithUCcreds(creds credentialsUC.CredUsecaseInf) TransportOption {
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
	pb.RegisterCredSvcServer(srv.srvRPC, srv.hndCreds)
	pb.RegisterUserSvcServer(srv.srvRPC, srv.hndUsers)
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
