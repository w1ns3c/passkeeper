package memstorage

import (
	"context"
	"fmt"
	"sync"

	"passkeeper/internal/entities"
	"passkeeper/internal/storage"
)

var (
	ErrBlobNotFound = fmt.Errorf("blob not exist")
	ErrUserNotFound = fmt.Errorf("user not exist")
)

type MemStorage struct {
	users   map[string]*entities.User
	usersMU *sync.RWMutex
	blobs   map[string][]*entities.CryptoBlob
	blobMU  *sync.RWMutex
}

func (m *MemStorage) Close() error {
	m.usersMU.Lock()
	m.blobMU.Lock()
	defer m.usersMU.Unlock()
	defer m.blobMU.Unlock()

	m.users = nil
	m.blobs = nil

	return nil
}

func NewMemStorage(ctx context.Context, options ...MemOptions) storage.Storage {
	storage := NewEmptyMemStorage(ctx)
	for _, option := range options {
		option(storage)
	}
	return storage
}

func NewEmptyMemStorage(ctx context.Context) *MemStorage {
	var m = &MemStorage{
		users:   make(map[string]*entities.User),
		usersMU: &sync.RWMutex{},
		blobs:   make(map[string][]*entities.CryptoBlob),
		blobMU:  &sync.RWMutex{},
	}

	err := m.Init(ctx)
	if err != nil {

		return nil
	}

	return m
}

type MemOptions func(*MemStorage)

// WithUsers func allow init storage
// with already registered usersUC
func WithUsers(users map[string]*entities.User) MemOptions {
	return func(s *MemStorage) {
		s.users = users
	}
}

// WithBlobs func allow init storage
// with encrypted blobs
func WithBlobs(passwords map[string][]*entities.CryptoBlob) MemOptions {
	return func(s *MemStorage) {
		s.blobs = passwords
	}
}

func (m *MemStorage) Init(ctx context.Context) error {
	if m.usersMU == nil {
		m.usersMU = &sync.RWMutex{}
	}

	if m.blobMU == nil {
		m.blobMU = &sync.RWMutex{}
	}

	if m.users == nil {
		m.users = make(map[string]*entities.User)
	}

	if m.blobs == nil {
		m.blobs = make(map[string][]*entities.CryptoBlob)
	}

	return nil
}

func (m *MemStorage) CheckConnection() error {
	return nil
}
