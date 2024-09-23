package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"passkeeper/internal/client/tui"
	"passkeeper/internal/entities/config/client"
)

var (
	BuildVersion = "N/A"
	BuildCommit  = "N/A"
	BuildDate    = "N/A"
)

func main() {
	args := client.CliParseArgs()

	if args.ShowVersion {
		fmt.Printf("Passkeeper Client\n"+
			" - version:    %s     %s\n"+
			" - build date: %s\n",
			BuildVersion, BuildCommit, BuildDate)

		return
	}

	tuiApp, _ := tui.NewTUIconf(args)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()

	if err := tuiApp.Run(ctx); err != nil {
		fmt.Println(err)
		tuiApp.Stop()
	}

}
