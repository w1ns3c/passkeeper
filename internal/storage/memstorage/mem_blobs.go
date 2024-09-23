package memstorage

import (
	"context"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/myerrors"
)

// AddBlob add blob to memory storage
func (m *MemStorage) AddBlob(ctx context.Context, blob *entities.CryptoBlob) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	if m.blobs[blob.UserID] == nil {
		m.blobs[blob.UserID] = make([]*entities.CryptoBlob, 1)
		m.blobs[blob.UserID][0] = blob
		return nil
	}

	m.blobs[blob.UserID] = append(m.blobs[blob.UserID], blob)
	return nil
}

// GetBlob get a blob from memory storage
func (m *MemStorage) GetBlob(ctx context.Context, userID, blobID string) (blob *entities.CryptoBlob, err error) {
	m.blobMU.RLock()
	defer m.blobMU.RUnlock()

	if m.blobs[userID] == nil {
		return nil, myerrors.ErrBlobNotFound
	}

	for _, pass := range m.blobs[userID] {
		if pass.ID == blobID {
			return pass, nil
		}
	}
	return nil, myerrors.ErrBlobNotFound
}

// GetAllBlobs return all blobs from memory storage for specific user
func (m *MemStorage) GetAllBlobs(ctx context.Context, userID string) (blobs []*entities.CryptoBlob, err error) {
	m.blobMU.RLock()
	defer m.blobMU.RUnlock()

	pass, ok := m.blobs[userID]
	if !ok {
		return nil, myerrors.ErrBlobNotFound
	}

	return pass, nil
}

// DeleteBlob delete blob from memory storage
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

	return myerrors.ErrBlobNotFound
}

// UpdateBlob change blob in memory storage
func (m *MemStorage) UpdateBlob(ctx context.Context, blob *entities.CryptoBlob) error {
	m.blobMU.Lock()
	defer m.blobMU.Unlock()

	for ind, pass := range m.blobs[blob.UserID] {
		if pass.ID == blob.ID {
			m.blobs[blob.UserID][ind] = blob
			return nil
		}
	}

	return myerrors.ErrBlobNotFound
}
