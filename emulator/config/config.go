// Package config contains the emulator configuration
package config

// -----------------------------------------------------------------------------
// Configuration
// -----------------------------------------------------------------------------

// Default configuration constants
const (
	// General
	DefaultAppTitle      = "emu8"
	DefaultMachineModel  = "Speccy"
	DefaultEmulatorAsync = false
	// Video
	DefaultVideoScale      = 2
	DefaultVideoFullScreen = false
	// Audio
	DefaultAudioFrecuency = 48000 // 48 KHz
	DefaultAudioMute      = false
)

// Config is the main configuration
type Config struct {
	AppTitle       string
	FileName       string
	MachineModel   string
	MachineOptions string
	EmulatorAsync  bool
	Video          VideoConfig
	Audio          AudioConfig
}

// VideoConfig the video configuration
type VideoConfig struct {
	Scale      int
	FullScreen bool
}

// AudioConfig the audio configuration
type AudioConfig struct {
	Frequency int
	Mute      bool
}

// config is the application main configuration
var config = new(Config)

// Get gets the main configuration
func Get() *Config {
	return config
}

// init initializes configuration
func init() {
	// General
	config.AppTitle = DefaultAppTitle
	config.MachineModel = DefaultMachineModel
	config.EmulatorAsync = DefaultEmulatorAsync
	// Video
	config.Video.Scale = DefaultVideoScale
	config.Video.FullScreen = DefaultVideoFullScreen
	// Audio
	config.Audio.Frequency = DefaultAudioFrecuency
	config.Audio.Mute = DefaultAudioMute
}
