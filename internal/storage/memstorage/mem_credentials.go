package memstorage

import (
	"context"

	"passkeeper/internal/entities"
)

func (m *MemStorage) AddCredential(ctx context.Context, userID string, password *entities.Credential) error {
	m.passMU.Lock()
	defer m.passMU.Unlock()

	if m.passwords[userID] == nil {
		m.passwords[userID] = make([]*entities.Credential, 1)
		m.passwords[userID][0] = password
		return nil
	}

	m.passwords[userID] = append(m.passwords[userID], password)
	return nil
}

func (m *MemStorage) GetCredential(ctx context.Context, userID, passwordID string) (password *entities.Credential, err error) {
	m.passMU.RLock()
	defer m.passMU.RUnlock()

	if m.passwords[userID] == nil {
		return nil, ErrPassNotFound
	}

	for _, pass := range m.passwords[userID] {
		if pass.ID == passwordID {
			return pass, nil
		}
	}
	return nil, ErrPassNotFound
}

func (m *MemStorage) GetAllCredentials(ctx context.Context, userID string) (passwords []*entities.Credential, err error) {
	m.passMU.RLock()
	defer m.passMU.RUnlock()

	pass, ok := m.passwords[userID]
	if !ok {
		return nil, ErrPassNotFound
	}

	return pass, nil
}

func (m *MemStorage) DeleteCredential(ctx context.Context, userID, passwordID string) error {
	m.passMU.Lock()
	defer m.passMU.Unlock()

	l := len(m.passwords[userID])
	for ind, pass := range m.passwords[userID] {
		if pass.ID == passwordID {
			for j := ind; j < l-1; j++ {
				m.passwords[userID][j] = m.passwords[userID][j+1]
			}
			m.passwords[userID][l-1] = nil
			m.passwords[userID] = m.passwords[userID][:l-1]
			return nil
		}
	}

	return ErrPassNotFound
}

func (m *MemStorage) UpdateCredential(ctx context.Context, userID string, password *entities.Credential) error {
	m.passMU.Lock()
	defer m.passMU.Unlock()

	for ind, pass := range m.passwords[userID] {
		if pass.ID == password.ID {
			m.passwords[userID][ind] = password
			return nil
		}
	}

	return ErrPassNotFound
}
