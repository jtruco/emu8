package config

import (
	"flag"

	"github.com/jtruco/emu8/machine"
)

// Default configuration constants
const (
	DefaultAppTitle     = "emu8"
	DefaultAppUI        = "sdl"
	DefaultMachineModel = machine.ZXSpectrum48k
	DefaultVideoScale   = 2
	DefaultFullScreen   = false
)

// Config is the main configuration
type Config struct {
	AppTitle     string
	AppUI        string
	MachineModel int
	RomFile      string
	FileName     string
	VideoScale   int
	FullScreen   bool
}

// config is the application main configuration
var config *Config

// Get gets the main configuration
func Get() *Config {
	return config
}

// Init parses configuration and arguments
func Init() bool {
	parseArgs()
	return true
}

// parse command line arguments
func parseArgs() {
	model := flag.String("m", "ZXSpectrum48k", "Machine model")
	flag.StringVar(&config.RomFile, "r", "", "Load rom file")
	flag.StringVar(&config.FileName, "f", "", "Load file")
	flag.IntVar(&config.VideoScale, "vs", DefaultVideoScale, "Video scale (1..4)")
	flag.BoolVar(&config.FullScreen, "vf", DefaultFullScreen, "Video in full screen mode")
	flag.Parse()
	config.MachineModel = machine.GetModel(*model, DefaultMachineModel)
	if len(flag.Args()) > 0 {
		config.FileName = flag.Args()[0]
	}
}

// init initializes configuration
func init() {
	config = &Config{}
	config.AppTitle = DefaultAppTitle
	config.AppUI = DefaultAppUI
	config.MachineModel = DefaultMachineModel
	config.VideoScale = DefaultVideoScale
	config.FullScreen = DefaultFullScreen
}