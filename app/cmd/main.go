// package main contains the main console aplicacion
package main

import (
	"log"

	"github.com/jtruco/emu8/app/cmd/ui/sdl"
	"github.com/jtruco/emu8/emulator"
)

// main program
func main() {
	// initialize emulator
	emu, err := emulator.GetDefault()
	if err != nil {
		log.Fatal("App : Could not initialize emulator")
	}
	log.Println("App : Emulator for machine:", emu.Machine().Config().Name)

	app := sdl.NewApp()
	if !app.Init(emu) {
		log.Fatal("App : Could not initialize application")
	}
	app.Run()
	app.End()
}
