package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"passkeeper/internal/client/tui"
	"passkeeper/internal/entities/config/client"
)

func main() {
	args := client.CliParseArgs()

	//ex, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(ex)
	//exPath := filepath.Dir(ex)
	//fmt.Println(exPath)

	tuiApp, _ := tui.NewTUIconf(args)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()

	if err := tuiApp.Run(ctx); err != nil {
		fmt.Println(err)
		tuiApp.Stop()
	}

}
