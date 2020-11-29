package main

import (
	"flag"

	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller/vfs"
)

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
	// desktop file system
	vfs.InitDesktop()
}
