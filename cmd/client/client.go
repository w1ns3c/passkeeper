package main

import (
	"fmt"
	"passkeeper/internal/client/tui"
	"passkeeper/internal/config/client"
)

func main() {
	args := client.CliParseArgs()

	tuiApp, _ := tui.NewTUIconf(args)
	if err := tuiApp.Run(); err != nil {
		fmt.Println(err)
	}
}
