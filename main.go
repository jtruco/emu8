package main

import (
	"github.com/jtruco/emu8/app"
	"github.com/jtruco/emu8/config"
)

// main program start
func main() {
	if !config.Init() {
		return
	}
	app := app.GetDefaultApp()
	if !app.Init() {
		return
	}
	app.Run()
	app.End()
}
