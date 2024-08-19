package memstorage

import (
	"context"
	"fmt"
	"sync"

	"passkeeper/internal/entities"
)

var (
	ErrPassNotFound = fmt.Errorf("password not exist")
	ErrUserNotFound = fmt.Errorf("user not exist")
)

type MemStorage struct {
	users     map[string]*entities.User
	usersMU   *sync.RWMutex
	passwords map[string][]*entities.Credential
	passMU    *sync.RWMutex
}

func NewMemStorage(options ...MemOptions) *MemStorage {
	storage := NewEmptyMemStorage()
	for _, option := range options {
		option(storage)
	}
	return storage
}

func NewEmptyMemStorage() *MemStorage {
	return &MemStorage{
		users:     make(map[string]*entities.User),
		usersMU:   &sync.RWMutex{},
		passwords: make(map[string][]*entities.Credential),
		passMU:    &sync.RWMutex{},
	}
}

type MemOptions func(*MemStorage)

// WithUsers func allow init storage
// with already registered usersUC
func WithUsers(users map[string]*entities.User) MemOptions {
	return func(s *MemStorage) {
		s.users = users
	}
}

// WithPasswords func allow init storage
// with passwords
func WithPasswords(passwords map[string][]*entities.Credential) MemOptions {
	return func(s *MemStorage) {
		s.passwords = passwords
	}
}

func (m *MemStorage) Init(ctx context.Context) error {
	if m.usersMU == nil {
		m.usersMU = &sync.RWMutex{}
	}

	if m.passMU == nil {
		m.passMU = &sync.RWMutex{}
	}

	if m.users == nil {
		m.users = make(map[string]*entities.User)
	}

	if m.passwords == nil {
		m.passwords = make(map[string][]*entities.Credential)
	}

	return nil
}
