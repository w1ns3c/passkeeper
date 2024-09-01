package memstorage

import (
	"context"

	"passkeeper/internal/entities"
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
		return nil, ErrUserNotFound
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
	m.usersMU.Lock()
	defer m.usersMU.Unlock()
	for _, u := range m.users {
		if u.ID == userID {
			return u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (m *MemStorage) GetUserByLogin(cxt context.Context, login string) (user *entities.User, err error) {
	m.usersMU.Lock()
	defer m.usersMU.Unlock()

	u, ok := m.users[login]
	if !ok {
		return nil, ErrUserNotFound
	}

	return u, nil
}

func (m *MemStorage) SaveChallenge(ctx context.Context, challenge string) error {
	//TODO implement me
	panic("implement me")
}
