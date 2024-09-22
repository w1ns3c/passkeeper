package blobsUC

import (
	"context"

	"passkeeper/internal/entities"
)

func (u *BlobUsecase) GetBlob(ctx context.Context, userToken, blobID string) (blob *entities.CryptoBlob, err error) {
	return u.storage.GetBlob(ctx, userToken, blobID)
}

func (u *BlobUsecase) AddBlob(ctx context.Context,
	userID string, blob *entities.CryptoBlob) error {

	//blob.Password, err = EncryptPass(blob.Password)
	//blob.ID = hashes.GeneratePassID(sec, salt)

	//blob.ID = hashes.GeneratePassID2()
	//VerifyCredDate(blob)

	return u.storage.AddBlob(ctx, userID, blob)
}

func (u *BlobUsecase) UpdBlob(ctx context.Context,
	userID string, blob *entities.CryptoBlob) error {

	return u.storage.UpdateBlob(ctx, userID, blob)
}

func (u *BlobUsecase) DelBlob(ctx context.Context,
	userToken, blobID string) error {

	return u.storage.DeleteBlob(ctx, userToken, blobID)
}

func (u *BlobUsecase) ListBlobs(ctx context.Context,
	userID string) (blobs []*entities.CryptoBlob, err error) {

	return u.storage.GetAllBlobs(ctx, userID)
}

func (u *BlobUsecase) GetUserSalt(ctx context.Context, userID string) string {
	//TODO implement me
	panic("implement me")
}
