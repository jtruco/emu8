package main

import (
	"github.com/jtruco/emu8/app"
	"github.com/jtruco/emu8/config"
	"github.com/jtruco/emu8/emulator"
)

// main program
func main() {
	if !config.Init() {
		return
	}
	emu := emulator.GetDefault()
	if emu == nil {
		return
	}
	app := app.GetDefault()
	if !app.Init(emu) {
		return
	}
	app.Run()
	app.End()
}
