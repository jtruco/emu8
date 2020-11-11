package main

import (
	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
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
	app := NewApp()
	if !app.Init(emu) {
		return
	}
	app.Run()
	app.End()
}
