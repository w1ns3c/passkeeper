package server

import (
	"context"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	mygrpc "passkeeper/internal/transport/grpc"
	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

func TestNewServer(t *testing.T) {
	type args struct {
		opts []SrvOption
	}

	var (
		ctx  = context.Background()
		uc1  = usersUC.NewUserUsecase(ctx)
		uc2  = blobsUC.NewBlobUCWithOpts()
		addr = "127.0.0.1:123"

		logger = zerolog.New(io.Discard)
		srv    = &Server{
			addr:  addr,
			users: uc1,
			blobs: uc2,
			log:   &logger,
		}
	)

	tests := []struct {
		name    string
		args    args
		wantSrv *Server
		wantErr bool
	}{
		{
			name: "TestNewServer 1",
			args: args{
				opts: []SrvOption{
					WithLogger(&logger),
					WithAddr(addr),
					WithUCusers(uc1),
					WithUCblobs(uc2),
				},
			},
			wantSrv: srv,
			wantErr: false,
		},
		{
			name: "TestNewServer 2: Blobs nil",
			args: args{
				opts: []SrvOption{
					WithUCusers(uc1),
				},
			},
			wantSrv: &Server{users: uc1},
			wantErr: true,
		},
		{
			name:    "TestNewServer 3: Users nil",
			args:    args{},
			wantSrv: &Server{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSrv, err := NewServer(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			require.Equal(t, tt.wantSrv.log, gotSrv.log)
			require.Equal(t, tt.wantSrv.addr, gotSrv.addr)
			require.Equal(t, tt.wantSrv.users, gotSrv.users)
			require.Equal(t, tt.wantSrv.blobs, gotSrv.blobs)
		})
	}
}

func TestServer_Stop(t *testing.T) {
	type fields struct {
		addr      string
		users     usersUC.UserUsecaseInf
		blobs     blobsUC.BlobUsecaseInf
		transport *mygrpc.TransportGRPC
		log       *zerolog.Logger
	}

	var (
		ctx    = context.Background()
		uc1    = usersUC.NewUserUsecase(ctx)
		uc2    = blobsUC.NewBlobUCWithOpts()
		logger = zerolog.New(io.Discard)
	)

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestServer_Stop 1",
			fields: fields{
				transport: &mygrpc.TransportGRPC{},
				log:       &logger,
				users:     uc1,
				blobs:     uc2,
				addr:      "127.0.0.1:123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{
				addr:      tt.fields.addr,
				users:     tt.fields.users,
				blobs:     tt.fields.blobs,
				transport: tt.fields.transport,
				log:       tt.fields.log,
			}
			if err := s.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Nil(t, s.users)
			require.Nil(t, s.blobs)
			require.Nil(t, s.transport)
		})
	}
}
