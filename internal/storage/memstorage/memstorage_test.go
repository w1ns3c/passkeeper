package memstorage

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/structs"
)

func TestMemStorage_Init(t *testing.T) {
	type fields struct {
		users   map[string]*structs.User
		usersMU *sync.RWMutex
		blobs   map[string][]*structs.CryptoBlob
		blobMU  *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Test 1: init nil usersMU",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				users:   tt.fields.users,
				usersMU: tt.fields.usersMU,
				blobs:   tt.fields.blobs,
				blobMU:  tt.fields.blobMU,
			}
			if err := m.Init(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_CheckConnection(t *testing.T) {
	var m = &MemStorage{}
	if err := m.CheckConnection(); err != nil {
		t.Errorf("CheckConnection() error = %v", err)
	}
}

func TestMemStorage_Close(t *testing.T) {
	m := &MemStorage{}
	m.Init(context.Background())

	if err := m.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	require.Nil(t, m.users)
	require.Nil(t, m.blobs)
}

func TestNewMemStorage(t *testing.T) {
	var (
		users   = map[string]*structs.User{}
		userID1 = "11111"
		blobs   = map[string][]*structs.CryptoBlob{
			userID1: {
				{
					ID:     "123213",
					UserID: userID1,
					Blob:   "ffffffffff",
				},
			},
		}
		options = []MemOptions{
			WithBlobs(blobs),
			WithUsers(users),
		}

		u = &MemStorage{}
	)

	for _, opt := range options {
		opt(u)
	}

	got := NewMemStorage(context.Background(), options...)
	gotBlobs, _ := got.GetAllBlobs(context.Background(), userID1)

	require.Equal(t, blobs[userID1], gotBlobs)
}
