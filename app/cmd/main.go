package main

import (
	"github.com/jtruco/emu8/app/cmd/ui/sdl"
	"github.com/jtruco/emu8/emulator"
)

// main program
func main() {
	emu := emulator.GetDefault()
	if emu == nil {
		return
	}
	app := sdl.NewApp()
	if !app.Init(emu) {
		return
	}
	app.Run()
	app.End()
}
