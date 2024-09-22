package blobsUC

import (
	"context"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/myerrors"
)

func (u *BlobUsecase) GetBlob(ctx context.Context, userID, blobID string) (blob *entities.CryptoBlob, err error) {
	return u.storage.GetBlob(ctx, userID, blobID)
}

func (u *BlobUsecase) AddBlob(ctx context.Context,
	userID string, blob *entities.CryptoBlob) error {

	// user try to add blob to someone else's
	if userID != blob.UserID {
		return myerrors.ErrBlobsUserIDdifferent
	}

	return u.storage.AddBlob(ctx, blob)
}

func (u *BlobUsecase) UpdBlob(ctx context.Context,
	userID string, blob *entities.CryptoBlob) error {

	// user try to update someone else's blob
	if userID != blob.UserID {
		return myerrors.ErrBlobsUserIDdifferent
	}

	return u.storage.UpdateBlob(ctx, blob)
}

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

func (u *BlobUsecase) ListBlobs(ctx context.Context,
	userID string) (blobs []*entities.CryptoBlob, err error) {

	return u.storage.GetAllBlobs(ctx, userID)
}
