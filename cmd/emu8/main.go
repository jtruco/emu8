// package main contains the main console aplicacion
package main

import (
	"log"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/ui/sdl"
)

// main program
func main() {
	// initialize emulator
	emu, err := emulator.GetDefault()
	if err != nil {
		log.Fatal("App : Could not initialize emulator: ", err.Error())
	}
	log.Println("App : Emulator for machine:", emu.Machine().Config().Name)

	app := sdl.NewApp()
	if err := app.Init(emu); err != nil {
		log.Fatal("App : Could not initialize application: ", err.Error())
	}
	app.Run()
	app.End()
}
