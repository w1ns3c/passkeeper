package memstorage

import (
	"context"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/myerrors"
)

// CheckUserExist check user's login in memory storage
func (m *MemStorage) CheckUserExist(ctx context.Context, login string) (exist bool, err error) {
	m.usersMU.RLock()
	defer m.usersMU.RUnlock()

	_, ok := m.users[login]
	return ok, nil
}

// SaveUser save entities.User in memory storage
func (m *MemStorage) SaveUser(ctx context.Context, u *entities.User) error {
	m.usersMU.Lock()
	defer m.usersMU.Unlock()

	m.users[u.Login] = u
	m.blobs[u.ID] = make([]*entities.CryptoBlob, 0)

	return nil
}

// GetUserByID return entities.User by userID from memory storage
func (m *MemStorage) GetUserByID(cxt context.Context, userID string) (user *entities.User, err error) {
	m.usersMU.Lock()
	defer m.usersMU.Unlock()
	for _, u := range m.users {
		if u.ID == userID {
			return u, nil
		}
	}

	return nil, myerrors.ErrUserNotFound
}

// GetUserByLogin return entities.User by login from memory storage
func (m *MemStorage) GetUserByLogin(cxt context.Context, login string) (user *entities.User, err error) {
	m.usersMU.Lock()
	defer m.usersMU.Unlock()

	u, ok := m.users[login]
	if !ok {
		return nil, myerrors.ErrUserNotFound
	}

	return u, nil
}
