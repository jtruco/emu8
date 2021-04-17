package main

import (
	"flag"

	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller/vfs"
)

func init() {
	conf := config.Get()

	// parse config parameters
	flag.StringVar(&conf.FileName, "file", "", "Load file")
	flag.StringVar(&conf.MachineModel, "model", config.DefaultMachineModel, "Machine model")
	flag.StringVar(&conf.MachineOptions, "mopts", "", "Machine options")
	flag.IntVar(&conf.Video.Scale, "scale", config.DefaultVideoScale, "Video scale (1..3)")
	flag.BoolVar(&conf.Video.FullScreen, "fullscreen", config.DefaultVideoFullScreen, "Video in full screen mode")
	flag.BoolVar(&conf.Audio.Mute, "mute", config.DefaultAudioMute, "Audio Mute")
	flag.Parse()
	if len(flag.Args()) > 0 {
		conf.FileName = flag.Args()[0]
	}

	// validate parameters
	if conf.Video.Scale < 1 || conf.Video.Scale > 3 {
		conf.Video.Scale = config.DefaultVideoScale
	}

	// init desktop vfs
	vfs.InitDesktop()
}
