package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"

	"passkeeper/internal/storage"
)

type PostgresStorage struct {
	db  *sql.DB
	url string
	log *zerolog.Logger
}

const (
	TableUsers = "users"
	TableBlobs = "blobs"

	fieldUserID = "userID"
	fieldLogin  = "login"
	fieldHash   = "hash"
	fieldSalt   = "salt"
	fieldSecret = "secret"
	fieldEmail  = "email"
	fieldPhone  = "phone"

	fieldBlobID   = "blobID"
	fiedlBlobData = "blobData"
)

// NewStorage is constructor for correct connect to DB and init tables
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
		"%s varchar primary KEY UNIQUE,"+
		"%s varchar NOT NULL UNIQUE,"+
		"%s varchar NOT NULL,"+
		"%s varchar NOT NULL,"+
		"%s varchar NOT NULL,"+
		"%s varchar,"+ // currently not using
		"%s varchar,"+ // currently not using
		"CONSTRAINT users_fk FOREIGN KEY (userID) REFERENCES public.%s(userID));",
		TableUsers,
		fieldUserID, fieldLogin, fieldHash, fieldSalt, fieldSecret, fieldEmail, fieldPhone,
		TableUsers)

	queryTb2 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"%s varchar primary KEY,"+
		"%s varchar NOT NULL UNIQUE,"+
		"%s varchar NOT NULL,"+
		"CONSTRAINT blobs_fk FOREIGN KEY (userID) REFERENCES public.%s(userID));",
		TableBlobs,
		fieldUserID, fieldBlobID, fiedlBlobData,
		TableBlobs)

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
