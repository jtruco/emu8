// Package config contains the emulator configuration
package config

import "github.com/jtruco/emu8/emulator/machine"

// Default configuration constants
const (
	DefaultAppTitle        = "emu8"
	DefaultAppFile         = ""
	DefaultMachineModel    = "Speccy"
	DefaultMachineOptions  = ""
	DefaultVideoScale      = 2
	DefaultVideoFullScreen = false
	DefaultAudioFrequency  = machine.DefaultAudioFrequency // 48 KHz
	DefaultAudioMute       = false
)

// -----------------------------------------------------------------------------
// Configuration
// -----------------------------------------------------------------------------

// Config is the main configuration
type Config struct {
	App     AppConfig
	Machine MachineConfig
	Video   VideoConfig
	Audio   AudioConfig
}

// AppConfig is the application configuration
type AppConfig struct {
	Title string // Application base title
	File  string // File to load
}

// MachineConfig is the machine configuration
type MachineConfig struct {
	Model   string // Machine model ID
	Options string // Machine options
}

// VideoConfig is the video configuration
type VideoConfig struct {
	Scale      int  // Video scale
	FullScreen bool // Fullscreen mode
}

// AudioConfig is the audio configuration
type AudioConfig struct {
	Frequency int  // Audio frequency
	Mute      bool // Mute audio
}

// New returns a configuration with the default values
func New() *Config {
	config := new(Config)
	config.App.Title = DefaultAppTitle
	config.App.File = DefaultAppFile
	config.Machine.Model = DefaultMachineModel
	config.Machine.Options = DefaultMachineOptions
	config.Video.Scale = DefaultVideoScale
	config.Video.FullScreen = DefaultVideoFullScreen
	config.Audio.Frequency = DefaultAudioFrequency
	config.Audio.Mute = DefaultAudioMute
	return config
}

// MachineSettings returns the machine construction settings
func (config *Config) MachineSettings() machine.Settings {
	return machine.Settings{
		ModelID:        config.Machine.Model,
		AudioFrequency: config.Audio.Frequency,
		Options:        config.Machine.Options}
}
