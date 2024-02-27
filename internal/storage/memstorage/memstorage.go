package memstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/w1nsec/passkeeper/internal/entities"
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
