// package main contains the main console application
package main

import (
	"log"
	"os"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/ui"
)

// main program
func main() {
	// desktop file system : current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Main : Could not get working directory: ", err.Error())
	}
	// initialize emulator
	emu, err := emulator.GetEmulator(vfs.NewDesktopFileSystem(cwd))
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
