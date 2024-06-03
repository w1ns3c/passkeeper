package main

import (
	"github.com/w1ns3c/passkeeper/internal/client/tui"
)

func main() {
	tuiApp, _ := tui.NewTUI()
	if err := tuiApp.App.Run(); err != nil {
		panic(err)
	}
}
