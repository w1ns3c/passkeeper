package grpc

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

func TestNewTransportGRPC(t *testing.T) {
	type args struct {
		opts []TransportOption
	}
	var (
		ctx        = context.Background()
		uc1        = usersUC.NewUserUsecase(ctx)
		uc2        = blobsUC.NewBlobUCWithOpts()
		addr       = "127.0.0.1:123"
		netAddr, _ = net.ResolveTCPAddr("tcp", addr)
		logger     = zerolog.New(io.Discard)
		trans      = &TransportGRPC{
			log:   &logger,
			addr:  netAddr,
			users: uc1,
			blobs: uc2,
		}
	)

	tests := []struct {
		name    string
		args    args
		wantSrv *TransportGRPC
		wantErr bool
	}{
		{
			name: "Test 1: success",
			args: args{
				opts: []TransportOption{
					WithLogger(&logger),
					WithAddr(addr),
					WithUCusers(uc1),
					WithUCcreds(uc2),
				},
			},
			wantSrv: trans,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSrv, err := NewTransportGRPC(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransportGRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.wantSrv.log, gotSrv.log)
			require.Equal(t, tt.wantSrv.addr, gotSrv.addr)
			require.Equal(t, tt.wantSrv.users, gotSrv.users)
			require.Equal(t, tt.wantSrv.blobs, gotSrv.blobs)

		})
	}
}
