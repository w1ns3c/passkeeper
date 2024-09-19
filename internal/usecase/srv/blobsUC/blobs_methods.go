package blobsUC

import (
	"context"

	"passkeeper/internal/entities"
)

func (u *BlobUsecase) GetBlob(ctx context.Context, userToken, blobID string) (blob *entities.CryptoBlob, err error) {
	return u.storage.GetCredential(ctx, userToken, blobID)
}

func (u *BlobUsecase) AddBlob(ctx context.Context,
	userID string, blob *entities.CryptoBlob) error {

	//blob.Password, err = EncryptPass(blob.Password)
	//blob.ID = hashes.GeneratePassID(sec, salt)

	//blob.ID = hashes.GeneratePassID2()
	//VerifyCredDate(blob)

	return u.storage.AddCredential(ctx, userID, blob)
}

func (u *BlobUsecase) UpdBlob(ctx context.Context,
	userID string, blob *entities.CryptoBlob) error {

	return u.storage.UpdateCredential(ctx, userID, blob)
}

func (u *BlobUsecase) DelBlob(ctx context.Context,
	userToken, blobID string) error {

	return u.storage.DeleteCredential(ctx, userToken, blobID)
}

func (u *BlobUsecase) ListBlobs(ctx context.Context,
	userID string) (blobs []*entities.CryptoBlob, err error) {

	return u.storage.GetAllCredentials(ctx, userID)
}

func (u *BlobUsecase) GetUserSalt(ctx context.Context, userID string) string {
	//TODO implement me
	panic("implement me")
}
