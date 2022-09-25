// package main contains the main console aplicacion
package main

import (
	"log"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/ui"
)

// main program
func main() {
	// initialize emulator
	emu, err := emulator.GetEmulator()
	if err != nil {
		log.Fatal("Main : Could not initialize emulator: ", err.Error())
	}
	log.Println("Main : Emulator for machine:", emu.Machine().Config().Name)

	// initialize UI application
	app := ui.GetApp()
	app.Connect(emu)
	if err := app.Init(); err != nil {
		log.Fatal("Main : Could not initialize UI: ", err.Error())
	}

	// run application
	defer app.Quit()
	app.Run()
}
