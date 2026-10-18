package main

import (
	"flag"

	"github.com/jtruco/emu8/emulator/config"
)

// parseFlags returns the configuration from the command line arguments
func parseFlags(args []string) (*config.Config, error) {
	conf := config.New()

	// parse config parameters
	flags := flag.NewFlagSet("emu8", flag.ContinueOnError)
	flags.StringVar(&conf.App.File, "file", config.DefaultAppFile, "Load file")
	flags.StringVar(&conf.Machine.Model, "model", config.DefaultMachineModel, "Machine model")
	flags.StringVar(&conf.Machine.Options, "options", config.DefaultMachineOptions, "Machine options")
	flags.IntVar(&conf.Video.Scale, "scale", config.DefaultVideoScale, "Video scale (1..3)")
	flags.BoolVar(&conf.Video.FullScreen, "fullscreen", config.DefaultVideoFullScreen, "Video in full screen mode")
	flags.BoolVar(&conf.Audio.Mute, "mute", config.DefaultAudioMute, "Audio Mute")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if flags.NArg() > 0 {
		conf.App.File = flags.Arg(0)
	}

	// validate parameters
	if conf.Video.Scale < 1 || conf.Video.Scale > 3 {
		conf.Video.Scale = config.DefaultVideoScale
	}
	return conf, nil
}
