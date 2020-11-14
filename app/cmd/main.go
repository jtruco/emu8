package main

import (
	"flag"

	"github.com/jtruco/emu8/app/cmd/io"
	"github.com/jtruco/emu8/app/cmd/ui/sdl"
	"github.com/jtruco/emu8/emulator"
	"github.com/jtruco/emu8/emulator/config"
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

// init
func init() {
	// parse config parameters
	flag.StringVar(&config.Get().MachineModel, "m", config.DefaultMachineModel, "Machine model")
	flag.StringVar(&config.Get().MachineOptions, "o", "", "Machine options")
	flag.StringVar(&config.Get().FileName, "f", "", "Load file")
	flag.IntVar(&config.Get().Video.Scale, "vs", config.DefaultVideoScale, "Video scale (1..4)")
	flag.BoolVar(&config.Get().Video.FullScreen, "vf", config.DefaultVideoFullScreen, "Video in full screen mode")
	flag.BoolVar(&config.Get().Audio.Mute, "am", config.DefaultAudioMute, "Audio Mute")
	flag.Parse()
	if len(flag.Args()) > 0 {
		config.Get().FileName = flag.Args()[0]
	}
	// default file system
	io.DefaultFileSystem()
}
