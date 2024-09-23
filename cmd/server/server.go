package main

import (
	"context"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/config"
	cnf "passkeeper/internal/entities/config/server"
	"passkeeper/internal/entities/logger"
	"passkeeper/internal/server"
	"passkeeper/internal/storage/dbstorage/postgres"
	"passkeeper/internal/storage/memstorage"
	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"
)

func main() {

	conf := cnf.SrvParseArgs()

	saltLen := config.UserPassSaltLen

	ctx := context.Background()

	lg := logger.Init(conf.LogLevel)
	lg.Info().Msg("[+] Logger init:  done")

	// try to connect to DB
	storage, err := postgres.NewStorage(ctx, conf.DBurl)
	if err != nil {
		lg.Warn().Err(err).
			Msg("fail to init DB storage, use memory storage")

		storage = memstorage.NewMemStorage(ctx)
	} else {
		lg.Info().Str("db", entities.HideDBpass(conf.DBurl)).
			Msg("[+] Storage init: done (successfully connected to DB)")
	}

	// 	 TODO Change this
	users, blobs := InitTestData()

	//storage := memstorage.NewMemStorage(
	//	memstorage.WithUsers(users),
	//	memstorage.WithBlobs(blobs))
	//lg.Info().Msg("[i] Storage init: done")

	errs := InitTestDB(storage, users, blobs)
	for _, err = range errs {
		lg.Error().Err(err).Send()
	}

	err = TestGetUser(storage)
	if err != nil {
		lg.Error().Err(err).Send()
	}

	credsUC := blobsUC.NewBlobUCWithOpts(
		blobsUC.WithContext(ctx),
		blobsUC.WithStorage(storage),
		blobsUC.WithLogger(lg),
	)

	uc := usersUC.NewUserUsecase(ctx).
		SetStorage(storage).
		SetLog(lg).
		SetSaltLen(saltLen)

	lg.Info().Msg("[+] Usecase init: done")

	srv, err := server.NewServer(
		server.WithAddr(conf.Addr),
		server.WithUCblobs(credsUC),
		server.WithUCusers(uc),
		server.WithLogger(lg),
	)
	if err != nil {
		lg.Error().Err(err).
			Msg("Failed to create server")
		return
	}

	lg.Info().Msg("[+] Server init:  done")

	err = srv.Run()
	if err != nil {
		lg.Error().Err(err).
			Msg("[+] Server stopping...")

		defer srv.Stop()
	}

}
