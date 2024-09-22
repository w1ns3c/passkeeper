package memstorage

import (
	"context"

	"passkeeper/internal/entities"
)

func (m *MemStorage) AddBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	if m.blobs[userID] == nil {
		m.blobs[userID] = make([]*entities.CryptoBlob, 1)
		m.blobs[userID][0] = blob
		return nil
	}

	m.blobs[userID] = append(m.blobs[userID], blob)
	return nil
}

func (m *MemStorage) GetBlob(ctx context.Context, userID, blobID string) (blob *entities.CryptoBlob, err error) {
	m.blobMU.RLock()
	defer m.blobMU.RUnlock()

	if m.blobs[userID] == nil {
		return nil, ErrBlobNotFound
	}

	for _, pass := range m.blobs[userID] {
		if pass.ID == blobID {
			return pass, nil
		}
	}
	return nil, ErrBlobNotFound
}

func (m *MemStorage) GetAllBlobs(ctx context.Context, userID string) (blobs []*entities.CryptoBlob, err error) {
	m.blobMU.RLock()
	defer m.blobMU.RUnlock()

	pass, ok := m.blobs[userID]
	if !ok {
		return nil, ErrBlobNotFound
	}

	return pass, nil
}

func (m *MemStorage) DeleteBlob(ctx context.Context, userID, blobID string) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	l := len(m.blobs[userID])
	for ind, pass := range m.blobs[userID] {
		if pass.ID == blobID {
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

func (m *MemStorage) UpdateBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	for ind, pass := range m.blobs[userID] {
		if pass.ID == blob.ID {
			m.blobs[userID][ind] = blob
			return nil
		}
	}

	return ErrBlobNotFound
}
