package memstorage

import (
	"context"

	"passkeeper/internal/entities"
)

func (m *MemStorage) AddCredential(ctx context.Context, userID string, password *entities.CryptoBlob) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	if m.blobs[userID] == nil {
		m.blobs[userID] = make([]*entities.CryptoBlob, 1)
		m.blobs[userID][0] = password
		return nil
	}

	m.blobs[userID] = append(m.blobs[userID], password)
	return nil
}

func (m *MemStorage) GetCredential(ctx context.Context, userID, passwordID string) (password *entities.CryptoBlob, err error) {
	m.blobMU.RLock()
	defer m.blobMU.RUnlock()

	if m.blobs[userID] == nil {
		return nil, ErrBlobNotFound
	}

	for _, pass := range m.blobs[userID] {
		if pass.ID == passwordID {
			return pass, nil
		}
	}
	return nil, ErrBlobNotFound
}

func (m *MemStorage) GetAllCredentials(ctx context.Context, userID string) (passwords []*entities.CryptoBlob, err error) {
	m.blobMU.RLock()
	defer m.blobMU.RUnlock()

	pass, ok := m.blobs[userID]
	if !ok {
		return nil, ErrBlobNotFound
	}

	return pass, nil
}

func (m *MemStorage) DeleteCredential(ctx context.Context, userID, passwordID string) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	l := len(m.blobs[userID])
	for ind, pass := range m.blobs[userID] {
		if pass.ID == passwordID {
			for j := ind; j < l-1; j++ {
				m.blobs[userID][j] = m.blobs[userID][j+1]
			}
			m.blobs[userID][l-1] = nil
			m.blobs[userID] = m.blobs[userID][:l-1]
			return nil
		}
	}

	return ErrBlobNotFound
}

func (m *MemStorage) UpdateCredential(ctx context.Context, userID string, password *entities.CryptoBlob) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	for ind, pass := range m.blobs[userID] {
		if pass.ID == password.ID {
			m.blobs[userID][ind] = password
			return nil
		}
	}

	return ErrBlobNotFound
}
