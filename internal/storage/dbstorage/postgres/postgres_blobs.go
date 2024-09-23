package postgres

import (
	"context"
	"fmt"
	"time"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/myerrors"
)

// AddBlob add new crypto blob
func (pg *PostgresStorage) AddBlob(ctx context.Context, blob *entities.CryptoBlob) error {
	var (
		query = fmt.Sprintf("INSERT INTO %s (%s, %s, %s) "+
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
	_, err = tx.ExecContext(ctx, query, blob.UserID, blob.ID, blob.Blob)
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}

// GetBlob return crypto blob by blobID and userID
func (pg *PostgresStorage) GetBlob(ctx context.Context, userID, blobID string) (blob *entities.CryptoBlob, err error) {
	var (
		query = fmt.Sprintf("SELECT %s,%s,%s "+
			"FROM %s WHERE %s=$1 and %s=$2 ;",
			fieldUserID, fieldBlobID, fieldBlobData,
			TableBlobs, fieldUserID, fieldBlobID)
	)

	rows, err := pg.db.QueryContext(ctx, query, userID, blobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blobs := make([]*entities.CryptoBlob, 0)
	for rows.Next() {
		var tmpBlob = &entities.CryptoBlob{}

		err = rows.Scan(&tmpBlob.UserID, &tmpBlob.ID, &tmpBlob.Blob)
		if err != nil {
			return nil, err
		}

		blobs = append(blobs, tmpBlob)
	}

	rerr := rows.Close()
	if rerr != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// More than one blob trying to return
	if len(blobs) != 1 {
		return nil, myerrors.ErrBlobWrongResult
	}

	return blobs[0], nil
}

// GetAllBlobs return all crypto blobs for userID
func (pg *PostgresStorage) GetAllBlobs(ctx context.Context, userID string) (blobs []*entities.CryptoBlob, err error) {
	var (
		query = fmt.Sprintf("SELECT %s, %s, %s "+
			"FROM %s WHERE %s=$1 ;",
			fieldUserID, fieldBlobID, fieldBlobData,
			TableBlobs, fieldUserID)
	)

	rows, err := pg.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blobs = make([]*entities.CryptoBlob, 0)
	for rows.Next() {
		var blob = &entities.CryptoBlob{}

		err = rows.Scan(&blob.UserID, &blob.ID, &blob.Blob)
		if err != nil {
			return nil, err
		}

		blobs = append(blobs, blob)
	}

	rerr := rows.Close()
	if rerr != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blobs, nil
}

// DeleteBlob delete blob by blobID and userID
func (pg *PostgresStorage) DeleteBlob(ctx context.Context, userID, blobID string) error {
	var (
		query = fmt.Sprintf("DELETE FROM %s WHERE %s=$1 and %s=$2;",
			TableBlobs, fieldUserID, fieldBlobID)
	)

	result, err := pg.db.ExecContext(ctx, query, userID, blobID)
	if err != nil {
		return err
	}
	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if num != 1 {
			return fmt.Errorf("expected to affect 1 row, affected %d", num)
		}
	}

	return nil

}

// UpdateBlob update crypto blob data by blobID and userID
func (pg *PostgresStorage) UpdateBlob(ctx context.Context, blob *entities.CryptoBlob) error {
	var (
		query = fmt.Sprintf("UPDATE %s SET %s=$1 WHERE %s=$2",
			TableBlobs, fieldBlobData, fieldBlobID,
		)
	)

	rows, err := pg.db.ExecContext(ctx, query, blob.Blob, blob.ID)
	if err != nil {
		return err
	}

	count, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", count)
	}
	return nil
}
