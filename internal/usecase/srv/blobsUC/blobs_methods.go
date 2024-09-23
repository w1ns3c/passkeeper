package blobsUC

import (
	"context"

	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/entities/structs"
)

// GetBlob return blob from storage by userID and blobID
func (u *BlobUsecase) GetBlob(ctx context.Context, userID, blobID string) (blob *structs.CryptoBlob, err error) {
	return u.storage.GetBlob(ctx, userID, blobID)
}

// AddBlob add blob to specific user
func (u *BlobUsecase) AddBlob(ctx context.Context,
	userID string, blob *structs.CryptoBlob) error {

	// user try to add blob to someone else's
	if userID != blob.UserID {
		return myerrors.ErrBlobsUserIDdifferent
	}

	return u.storage.AddBlob(ctx, blob)
}

// UpdBlob change specific blob
func (u *BlobUsecase) UpdBlob(ctx context.Context,
	userID string, blob *structs.CryptoBlob) error {

	// user try to update someone else's blob
	if userID != blob.UserID {
		return myerrors.ErrBlobsUserIDdifferent
	}

	return u.storage.UpdateBlob(ctx, blob)
}

// DelBlob delete specific blob
func (u *BlobUsecase) DelBlob(ctx context.Context,
	userID, blobID string) error {

	blob, err := u.storage.GetBlob(ctx, userID, blobID)
	if err != nil {
		return err
	}

	// user try to delete someone else's blob
	if blob.UserID != userID {
		return myerrors.ErrBlobsUserIDdifferent
	}

	return u.storage.DeleteBlob(ctx, userID, blobID)
}

// ListBlobs return all blobs for specific user
func (u *BlobUsecase) ListBlobs(ctx context.Context,
	userID string) (blobs []*structs.CryptoBlob, err error) {

	return u.storage.GetAllBlobs(ctx, userID)
}
