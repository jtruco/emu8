package main

import (
	"log"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/ui"
)

// wasm main entry point
func main() {

	// TESTING : foce load game sample
	config.Get().App.File = "manicminer.sna"

	// initialize emulator
	emu, err := emulator.GetEmulator()
	if err != nil {
		log.Fatal("App : Could not initialize emulator: ", err.Error())
	}
	log.Println("App : Emulator for machine:", emu.Machine().Config().Name)

	// initialize user interface
	app := ui.GetApp()
	app.Connect(emu)
	if err := app.Init(); err != nil {
		log.Fatal("App : Could not initialize UI: ", err.Error())
	}

	// run emulator ...
	defer app.Quit()
	app.Run()
	<-make(chan bool) // wait for ever
}
