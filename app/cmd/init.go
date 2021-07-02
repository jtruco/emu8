package main

import (
	"flag"

	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller/vfs"
)

func init() {
	conf := config.Get()

	// parse config parameters
	flag.StringVar(&conf.App.File, "file", "", "Load file")
	flag.BoolVar(&conf.Emulator.Async, "async", config.DefaultEmulatorAsync, "Asynchronous emulation")
	flag.StringVar(&conf.Machine.Model, "model", config.DefaultMachineModel, "Machine model")
	flag.StringVar(&conf.Machine.Options, "options", "", "Machine options")
	flag.IntVar(&conf.Video.Scale, "scale", config.DefaultVideoScale, "Video scale (1..3)")
	flag.BoolVar(&conf.Video.FullScreen, "fullscreen", config.DefaultVideoFullScreen, "Video in full screen mode")
	flag.BoolVar(&conf.Audio.Mute, "mute", config.DefaultAudioMute, "Audio Mute")
	flag.Parse()
	if len(flag.Args()) > 0 {
		conf.App.File = flag.Args()[0]
	}

	// validate parameters
	if conf.Video.Scale < 1 || conf.Video.Scale > 3 {
		conf.Video.Scale = config.DefaultVideoScale
	}

	// init desktop vfs
	vfs.InitDesktop()
}
