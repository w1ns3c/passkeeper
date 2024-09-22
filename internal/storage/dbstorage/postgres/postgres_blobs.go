package postgres

import (
	"context"
	"fmt"
	"time"

	"passkeeper/internal/entities"
)

func (pg *PostgresStorage) AddBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error {
	var (
		query = fmt.Sprintf("insert into %s (%s, %s, %s) "+
			"values ($1, $2, $3);",
			TableBlobs,
			fieldUserID, fieldBlobID, fieldBlobData)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Second*8)
	defer cancel()

	tx, err := pg.db.Begin()
	if err != nil {

		return err
	}
	_, err = tx.ExecContext(ctx, query, userID, blob.ID, blob.Blob)
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (pg *PostgresStorage) GetBlob(ctx context.Context, userID, blobID string) (blob *entities.CryptoBlob, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresStorage) GetAllBlobs(ctx context.Context, userID string) (blobs []*entities.CryptoBlob, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresStorage) DeleteBlob(ctx context.Context, userID, blobID string) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresStorage) UpdateBlob(ctx context.Context, userID string, blob *entities.CryptoBlob) error {
	//TODO implement me
	panic("implement me")
}
