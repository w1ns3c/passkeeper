package main

import (
	"context"
	"passkeeper/internal/config"
	"passkeeper/internal/logger"
	"passkeeper/internal/server"
	"passkeeper/internal/usecase/srv/credentialsUC"
	"passkeeper/internal/usecase/srv/usersUC"
	"time"

	"passkeeper/internal/entities"
	"passkeeper/internal/storage/memstorage"
	"passkeeper/internal/utils/hashes"
)

func main() {

	listenAddr := "localhost:8001"
	logLevel := "DEBUG"
	saltLen := config.UserPassSaltLen

	var (
		login1 = "user1"
		login2 = "user2"
		hash1  = hashes.Hash("123")
		hash2  = hashes.Hash("password")

		user1 = &entities.User{
			ID:    login1,
			Login: login1,
			Hash:  hash1,
		}

		user2 = &entities.User{
			ID:    login2,
			Login: login2,
			Hash:  hash2,
		}

		// Passwords
		password1 = &entities.Credential{
			ID:          "ID1111",
			Date:        time.Now(),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &entities.Credential{
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * 500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &entities.Credential{
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}
	)

	users := map[string]*entities.User{
		login1: user1,
		login2: user2,
	}

	passwords := map[string][]*entities.Credential{
		login1: {
			password1, password2,
		},
		login2: {
			password3,
		},
	}

	ctx := context.Background()

	lg := logger.Init(logLevel)

	lg.Info().Msg("Logger init: done")

	storage := memstorage.NewMemStorage(memstorage.WithUsers(users),
		memstorage.WithPasswords(passwords))
	lg.Info().Msg("Storage init: done")

	credsUC := credentialsUC.NewCredUCWithOpts(
		credentialsUC.WithContext(ctx),
		credentialsUC.WithStorage(storage),
		credentialsUC.WithLogger(lg),
	)

	usersUC := usersUC.NewUserUsecase().
		SetContext(ctx).
		SetStorage(storage).
		SetLog(lg).
		SetSaltLen(saltLen)

	lg.Info().Msg("Usecase init: done")

	srv, err := server.NewServer(
		server.WithAddr(listenAddr),
		server.WithUCcreds(credsUC),
		server.WithUCusers(usersUC),
		server.WithLogger(lg),
	)
	if err != nil {
		lg.Error().Err(err).
			Msg("Failed to create server")
		return
	}

	lg.Info().Msg("Server init: done")

	err = srv.Run()
	if err != nil {
		lg.Error().Err(err).Msg("Server stopped...")
		defer srv.Stop()
	}

}
