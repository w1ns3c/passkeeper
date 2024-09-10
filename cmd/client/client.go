package main

import (
	"fmt"
	"os"
	"passkeeper/internal/client/tui"
	"passkeeper/internal/config/client"
	"path/filepath"
)

func main() {
	args := client.CliParseArgs()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	fmt.Println(ex)
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	tuiApp, _ := tui.NewTUIconf(args)
	if err := tuiApp.Run(); err != nil {
		fmt.Println(err)
		tuiApp.Stop()
	}
}
