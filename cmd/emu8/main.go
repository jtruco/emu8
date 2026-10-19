// package main contains the main console application
package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/ui"
)

// main program
func main() {
	// configuration from command line
	conf, err := parseFlags(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	} else if err != nil {
		os.Exit(2) // flag package already printed the error and usage
	}

	// desktop file system : current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Main : Could not get working directory: ", err.Error())
	}

	// initialize emulator
	emu, err := emulator.Create(conf.MachineSettings(), vfs.NewDesktopFileSystem(cwd))
	if err != nil {
		log.Fatal("Main : Could not initialize emulator: ", err.Error())
	}
	log.Println("Main : Emulator for machine:", emu.Machine().Config().Name)
	if conf.App.File != "" {
		if err := emu.LoadFile(conf.App.File); err != nil {
			log.Fatal("Main : Could not load file: ", err.Error())
		}
	}

	// initialize UI application
	app := ui.GetApp(conf)
	app.Connect(emu)
	if err := app.Init(); err != nil {
		log.Fatal("Main : Could not initialize UI: ", err.Error())
	}

	// run application
	defer app.Quit()
	app.Run()
}
