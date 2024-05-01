package memstorage

import (
	"context"
	"fmt"

	"github.com/w1nsec/passkeeper/internal/entities"
)

var (
	ErrUserNotExist = fmt.Errorf("user not exist")
)

func (m *MemStorage) CheckUserExist(ctx context.Context, login string) (exist bool, err error) {
	m.usersMU.RLock()
	defer m.usersMU.RUnlock()

	_, ok := m.users[login]
	return ok, nil
}

func (m *MemStorage) LoginUser(ctx context.Context, login, hash string) (user *entities.User, err error) {
	m.usersMU.RLock()
	defer m.usersMU.RUnlock()

	user, ok := m.users[login]
	if !ok {
		return nil, ErrUserNotExist
	}

	return user, err
}

func (m *MemStorage) SaveUser(ctx context.Context, u *entities.User) error {
	m.usersMU.Lock()
	defer m.usersMU.Unlock()

	m.users[u.Login] = u
	//m.passwords[u.ID] = make([]entities.Credential, 0)
	return nil
}

func (m *MemStorage) GetUserByID(cxt context.Context, userID string) (user *entities.User, err error) {
	for _, u := range m.users {
		if u.ID == userID {
			return nil, err
		}
	}

	return nil, ErrUserNotExist
}
