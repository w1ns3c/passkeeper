package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/storage"
)

type PostgresStorage struct {
	db  *sql.DB
	url string
	log *zerolog.Logger
}

func NewStorage(ctx context.Context, dbURL string) (repo storage.Storage, err error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	repo = &PostgresStorage{
		db:  db,
		url: dbURL,
	}

	err = repo.CheckConnection()
	if err != nil {

		return nil, err
	}

	err = repo.Init(ctx)
	if err != nil {

		return nil, err
	}

	return repo, nil
}

// Init creat necessary tables in DB, if they not exist
func (pg *PostgresStorage) Init(ctx context.Context) error {
	if pg.db == nil {
		return fmt.Errorf("db not created")
	}

	queryTb1 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"userID varchar primary KEY UNIQUE,"+
		"login varchar NOT NULL UNIQUE,"+
		"hash varchar NOT NULL,"+
		"salt varchar NOT NULL,"+
		"secret varchar NOT NULL,"+
		"email varchar,"+ // currently not using
		"phone varchar,"+ // currently not using
		"CONSTRAINT users_fk FOREIGN KEY (userID) REFERENCES public.%s(userID));",
		config.TableUsers, config.TableUsers)

	queryTb2 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"userID varchar primary KEY,"+
		"blobID varchar NOT NULL UNIQUE,"+
		"blobData varchar NOT NULL,"+
		"CONSTRAINT blobs_fk FOREIGN KEY (userID) REFERENCES public.%s(userID));",
		config.TableBlobs, config.TableBlobs)

	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, queryTb1)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, queryTb2)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Close is wrapper for postgres Close
func (pg *PostgresStorage) Close() error {
	return pg.db.Close()
}

// CheckConnection check connection to Postgresql DB
// return error if it is not active
func (pg *PostgresStorage) CheckConnection() error {
	var err error

	pg.db, err = sql.Open("pgx", pg.url)
	if err != nil {
		pg.log.Error().Err(err).Send()

		return err
	}

	err = pg.db.Ping()
	if err != nil {

		return err
	}

	return nil
}
