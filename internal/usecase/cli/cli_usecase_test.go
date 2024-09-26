package cli

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/structs"
)

func TestClientUC_CredsLen(t *testing.T) {
	type args struct {
		usecase *ClientUC
	}

	tests := []struct {
		name  string
		args  args
		wantL int
	}{
		{
			name: "Test 1: 3 len",
			args: args{usecase: &ClientUC{
				Creds: []*structs.Credential{
					&structs.Credential{}, &structs.Credential{}, &structs.Credential{},
				},
				m: &sync.RWMutex{},
			}},
			wantL: 3,
		},
		{
			name: "Test len 2: 0 len",
			args: args{usecase: &ClientUC{
				Creds: []*structs.Credential{},
				m:     &sync.RWMutex{},
			}},
			wantL: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.usecase.CredsLen(); got != tt.wantL {
				t.Errorf("CredsLen() = %v, want %v", got, tt.wantL)
			}
		})
	}
}

func TestClientUC_CredsNotNil(t *testing.T) {
	type args struct {
		usecase *ClientUC
	}

	tests := []struct {
		name      string
		args      args
		gotNotNil bool
	}{
		{
			name: "Test CredsNil 1: creds exist",
			args: args{usecase: &ClientUC{
				Creds: []*structs.Credential{
					&structs.Credential{}, &structs.Credential{}, &structs.Credential{},
				},
				m: &sync.RWMutex{},
			}},
			gotNotNil: true,
		},
		{
			name: "Test CredsNil 2: creds are nil",
			args: args{usecase: &ClientUC{
				m: &sync.RWMutex{},
			}},
			gotNotNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.usecase.CredsNotNil(); got != tt.gotNotNil {
				t.Errorf("CredsNotNil() = %v, want %v", got, tt.gotNotNil)
			}
		})
	}
}

func TestClientUC_GetSyncTime(t *testing.T) {
	type args struct {
		usecase *ClientUC
	}

	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "Test time 1: one hour",
			args: args{usecase: &ClientUC{
				SyncTime: time.Hour,
			}},
			want: time.Hour,
		},
		{
			name: "Test time 2: time is empty",
			args: args{usecase: &ClientUC{}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.usecase.GetSyncTime(); got != tt.want {
				t.Errorf("GetSyncTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientUC_GetToken(t *testing.T) {
	type args struct {
		usecase *ClientUC
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test GetToken 1: exist",
			args: args{usecase: &ClientUC{
				Token: "some token",
			}},
			want: "some token",
		},
		{
			name: "Test GetToken 2: token is empty",
			args: args{usecase: &ClientUC{}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.usecase.GetToken(); got != tt.want {
				t.Errorf("GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientUC(t *testing.T) {
	var (
		lg = zerolog.New(io.Discard)
	)

	tests := []struct {
		name    string
		opts    []ClientUCoptions
		wantCli *ClientUC
		wantErr bool
	}{
		{
			name: "Test 1: with all options",
			opts: []ClientUCoptions{
				WithAddr("127.0.0.1:9999"),
				WithLogger(&lg),
				WithSyncTime(time.Second * 10),
			},
			wantCli: &ClientUC{
				SyncTime: time.Second * 10,
				log:      &lg,
				Addr:     "127.0.0.1:9999",
			},
		},
		{
			name: "Test 2: with all options (incorrect)",
			opts: []ClientUCoptions{
				WithAddr(""),
				WithLogger(&lg),
				WithSyncTime(time.Hour * 10),
			},
			wantCli: &ClientUC{
				log:  &lg,
				Addr: "localhost:8000",
			},
		},
		{
			name: "Test 3: without any options",
			opts: []ClientUCoptions{},
			wantCli: &ClientUC{
				log:  &lg,
				Addr: "localhost:8000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientUC(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientUC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Less(t, got.SyncTime, config.SyncMax)
			require.Greater(t, got.SyncTime, config.SyncMin)

			require.Equal(t, tt.wantCli.log, got.log)
			require.Equal(t, tt.wantCli.Addr, got.Addr)
			require.NotNil(t, got.FilesUC)
			require.NotNil(t, got.userSvc)
			require.NotNil(t, got.blobsSvc)
			require.NotNil(t, got.passSvc)
			require.NotNil(t, got.m)

		})
	}
}
