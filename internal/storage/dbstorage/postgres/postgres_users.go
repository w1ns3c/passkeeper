package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/myerrors"
)

func (pg *PostgresStorage) GetUserByID(cxt context.Context, userID string) (user *entities.User, err error) {
	return pg.getUserByAny(cxt, fieldLogin, userID)
}

func (pg *PostgresStorage) GetUserByLogin(ctx context.Context, login string) (user *entities.User, err error) {
	return pg.getUserByAny(ctx, fieldLogin, login)
}

func (pg *PostgresStorage) getUserByAny(ctx context.Context, field string, value string) (user *entities.User, err error) {
	var (
		query = fmt.Sprintf("SELECT *"+
			"FROM %s WHERE %s=$1;", TableUsers, field)
	)

	rows, err := pg.db.QueryContext(ctx, query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*entities.User, 0)
	for rows.Next() {
		var (
			userID string
			login  string
			hash   string
			secret string
			salt   string

			email sql.NullString
			phone sql.NullString
		)

		err = rows.Scan(&userID, &login, &hash, &secret, &salt, &email, &phone)
		if err != nil {
			return nil, err
		}

		result = append(result, &entities.User{
			ID:     userID,
			Login:  login,
			Hash:   hash,
			Salt:   secret,
			Secret: salt,
			Phone:  email.String,
			Email:  phone.String,
		})
	}

	rerr := rows.Close()
	if rerr != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// User not exist
	if len(result) == 0 {
		return nil, myerrors.ErrUsersNotExist
	}

	// More than one user trying to return
	if len(result) != 1 {
		return nil, myerrors.ErrUsersWrongResult
	}

	return result[0], nil
}

func (pg *PostgresStorage) CheckUserExist(ctx context.Context, login string) (exist bool, err error) {
	_, err = pg.GetUserByLogin(ctx, login)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (pg *PostgresStorage) SaveUser(ctx context.Context, user *entities.User) error {
	var (
		query = fmt.Sprintf("insert into %s (%s, %s, %s, %s, %s,%s, %s) "+
			"values ($1, $2, $3, $4, $5, $6, $7);",
			TableUsers,
			fieldUserID, fieldLogin, fieldHash,
			fieldSalt, fieldSecret, fieldEmail, fieldPhone)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Second*8)
	defer cancel()

	tx, err := pg.db.Begin()
	if err != nil {

		return err
	}
	_, err = tx.ExecContext(ctx, query, user.ID, user.Login, user.Hash,
		user.Salt, user.Secret, user.Email, user.Phone)
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}
