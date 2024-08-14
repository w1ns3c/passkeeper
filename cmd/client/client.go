package main

import (
	"github.com/w1ns3c/passkeeper/internal/client/tui"
)

func main() {
	addr := "localhost:8001"
	tuiApp, _ := tui.NewTUI(addr)
	if err := tuiApp.App.Run(); err != nil {
		panic(err)
	}
}
