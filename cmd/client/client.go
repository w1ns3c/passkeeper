package main

import (
	"passkeeper/internal/client/tui"
)

func main() {
	addr := "localhost:8001"
	debug := "debug"
	tuiApp, _ := tui.NewTUI(addr, debug)
	if err := tuiApp.App.Run(); err != nil {
		panic(err)
	}
}
