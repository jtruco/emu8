package config

import "github.com/jtruco/emu8/machine"

// Default configuration constants
const (
	defaultAppTitle     = "emu8"
	defaultAppUI        = "sdl"
	defaultMachineModel = machine.ZXSpectrum48k
	defaultVideoScale   = 2
	defaultFullScreen   = false
)

// Config is the main configuration
type Config struct {
	AppTitle     string
	AppUI        string
	MachineModel int
	Snapfile     string
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
func Init() bool { return true }

// init initializes configuration
func init() {
	config = &Config{}
	config.AppTitle = defaultAppTitle
	config.AppUI = defaultAppUI
	config.MachineModel = defaultMachineModel
	config.VideoScale = defaultVideoScale
	config.FullScreen = defaultFullScreen
}
