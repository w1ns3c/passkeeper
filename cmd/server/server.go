package main

import (
	"context"

	"passkeeper/internal/config"
	cnf "passkeeper/internal/config/server"
	"passkeeper/internal/logger"
	"passkeeper/internal/server"
	"passkeeper/internal/usecase/srv/blobsUC"
	"passkeeper/internal/usecase/srv/usersUC"

	"passkeeper/internal/storage/memstorage"
)

func main() {

	conf := cnf.SrvParseArgs()

	saltLen := config.UserPassSaltLen

	ctx := context.Background()

	lg := logger.Init(conf.LogLevel)
	lg.Info().Msg("[i] Logger init:  done")

	users, blobs := InitTestData()

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
