package main

import (
	"context"
	"github.com/w1ns3c/go-examples/crypto"
	"passkeeper/internal/config"
	cnf "passkeeper/internal/config/server"
	"passkeeper/internal/logger"
	"passkeeper/internal/server"
	"passkeeper/internal/usecase/srv/credentialsUC"
	"passkeeper/internal/usecase/srv/usersUC"
	"time"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/storage/memstorage"
)

func main() {

	conf := cnf.SrvParseArgs()

	saltLen := config.UserPassSaltLen

	var (
		login1  = "user1"
		pass1   = "123"
		userid1 = login1 + "_ID"

		login2  = "user2"
		pass2   = "password"
		userid2 = login2 + "_ID"

		hash1 = hashes.Hash(pass1)
		hash2 = hashes.Hash(pass2)

		userSalt1, _ = crypto.GenRandStr(config.UserPassSaltLen)
		userSalt2, _ = crypto.GenRandStr(config.UserPassSaltLen)

		cryptHash1, _ = hashes.GenerateCryptoHash(hash1, userSalt1)
		cryptHash2, _ = hashes.GenerateCryptoHash(hash2, userSalt2)

		secret1, _      = hashes.GenerateSecret(config.UserPassSaltLen)
		secret2, _      = hashes.GenerateSecret(config.UserPassSaltLen)
		cryptSecret1, _ = hashes.EncryptSecret(secret1, hash1)
		cryptSecret2, _ = hashes.EncryptSecret(secret2, hash2)

		user1 = &entities.User{
			ID:     userid1,
			Login:  login1,
			Hash:   cryptHash1,
			Salt:   userSalt1,
			Secret: cryptSecret1,
		}

		user2 = &entities.User{
			ID:     userid2,
			Login:  login2,
			Hash:   cryptHash2,
			Salt:   userSalt2,
			Secret: cryptSecret2,
		}

		// Passwords
		password1 = &entities.Credential{
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &entities.Credential{
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
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

	key1, _ := hashes.GenerateCredsSecret(pass1, user1.ID, cryptSecret1)
	key2, _ := hashes.GenerateCredsSecret(pass2, user2.ID, cryptSecret2)

	users := map[string]*entities.User{
		login1: user1,
		login2: user2,
	}
	//passwords := map[string][]*entities.Credential{
	//	login1: {
	//		password1, password2,
	//	},
	//	login2: {
	//		password3,
	//	},
	//}

	blob1, _ := hashes.EncryptBlob(password1, key1)
	blob2, _ := hashes.EncryptBlob(password2, key1)
	blob3, _ := hashes.EncryptBlob(password3, key2)

	blobs := map[string][]*entities.CredBlob{
		userid1: {
			blob1, blob2,
		},
		userid2: {
			blob3,
		},
	}

	ctx := context.Background()

	lg := logger.Init(conf.LogLevel)
	lg.Info().Msg("[i] Logger init:  done")

	storage := memstorage.NewMemStorage(
		memstorage.WithUsers(users),
		memstorage.WithBlobs(blobs))
	lg.Info().Msg("[i] Storage init: done")

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

	lg.Info().Msg("[i] Usecase init: done")

	srv, err := server.NewServer(
		server.WithAddr(conf.Addr),
		server.WithUCcreds(credsUC),
		server.WithUCusers(usersUC),
		server.WithLogger(lg),
	)
	if err != nil {
		lg.Error().Err(err).
			Msg("Failed to create server")
		return
	}

	lg.Info().Msg("[i] Server init:  done")

	err = srv.Run()
	if err != nil {
		lg.Error().Err(err).Msg("[+] Server stopped...")
		defer srv.Stop()
	}

}
