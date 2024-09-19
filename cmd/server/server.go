package main

import (
	"context"
	"time"

	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/config"
	cnf "passkeeper/internal/config/server"
	"passkeeper/internal/logger"
	"passkeeper/internal/server"
	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"

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
			Type:        entities.UserCred,
			ID:          "ID1111",
			Date:        time.Now().Add(time.Second * -200),
			Resource:    "localhost1111",
			Login:       "my_favorite_username1111",
			Password:    "my_favorite_password1111",
			Description: "some description1111",
		}
		password2 = &entities.Credential{
			Type:        entities.UserCred,
			ID:          "ID2222",
			Date:        time.Now().Add(time.Second * -500),
			Resource:    "localhost2222",
			Login:       "admin2222",
			Password:    "secret password2222",
			Description: "some new description2222",
		}
		password3 = &entities.Credential{
			Type:        entities.UserCred,
			ID:          "superID3333",
			Date:        time.Now(),
			Resource:    "localhost3333",
			Login:       "my_favorite_username333",
			Password:    "my_favorite_password3333",
			Description: "some description3333",
		}

		testCards = []*entities.Card{
			{
				Type:        entities.UserCard,
				Name:        "test1",
				Bank:        entities.Banks[0],
				Person:      "string",
				Number:      122222222222,
				CVC:         232,
				Expiration:  "33/44",
				PIN:         3333,
				Description: "test description only",
			},
			{
				Type:        entities.UserCard,
				Name:        "test333331",
				Bank:        entities.Banks[2],
				Person:      "Major Tom",
				Number:      234872398472,
				CVC:         23244444,
				Expiration:  "11/11",
				PIN:         11111,
				Description: "test description2",
			},
		}

		testNotes = []*entities.Note{
			{
				Type: entities.UserNote,
				Name: "test1",
				Body: "Hello\nWorld!",
			},
			{
				Type: entities.UserNote,
				Name: "HELLO 222222",
				Body: "Hello\nWorld! 9234928309482390480298340923809840",
			},
			{
				Type: entities.UserNote,
				Name: "New Test Blob",
				Body: "Hello\nWorld! Amigo",
			},
		}
	)

	key1, _ := hashes.GenerateCredsSecret(pass1, user1.ID, cryptSecret1)
	key2, _ := hashes.GenerateCredsSecret(pass2, user2.ID, cryptSecret2)

	users := map[string]*entities.User{
		login1: user1,
		login2: user2,
	}

	blob1, _ := hashes.EncryptBlob(password1, key1)
	blob2, _ := hashes.EncryptBlob(password2, key1)
	blob3, _ := hashes.EncryptBlob(password3, key2)
	blob4, _ := hashes.EncryptBlob(testCards[0], key1)
	blob5, _ := hashes.EncryptBlob(testCards[1], key1)
	blob6, _ := hashes.EncryptBlob(testNotes[0], key1)
	blob7, _ := hashes.EncryptBlob(testNotes[1], key1)
	blob8, _ := hashes.EncryptBlob(testNotes[2], key1)

	blobs := map[string][]*entities.CryptoBlob{
		userid1: {
			blob1, blob2, blob4, blob5, blob6, blob7, blob8,
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

	credsUC := blobsUC.NewBlobUCWithOpts(
		blobsUC.WithContext(ctx),
		blobsUC.WithStorage(storage),
		blobsUC.WithLogger(lg),
	)

	uc := usersUC.NewUserUsecase().
		SetContext(ctx).
		SetStorage(storage).
		SetLog(lg).
		SetSaltLen(saltLen)

	lg.Info().Msg("[i] Usecase init: done")

	srv, err := server.NewServer(
		server.WithAddr(conf.Addr),
		server.WithUCcreds(credsUC),
		server.WithUCusers(uc),
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
