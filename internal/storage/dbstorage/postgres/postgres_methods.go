package postgres

import (
	"context"

	"passkeeper/internal/entities"
)

func (pg *PostgresStorage) AddCredential(ctx context.Context, userID string, password *entities.CryptoBlob) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresStorage) GetCredential(ctx context.Context, userID, passwordID string) (password *entities.CryptoBlob, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresStorage) GetAllCredentials(ctx context.Context, userID string) (passwords []*entities.CryptoBlob, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresStorage) DeleteCredential(ctx context.Context, userID, passwordID string) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresStorage) UpdateCredential(ctx context.Context, userID string, password *entities.CryptoBlob) error {
	//TODO implement me
	panic("implement me")
}
