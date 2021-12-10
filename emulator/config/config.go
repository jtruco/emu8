// Package config contains the emulator configuration
package config

// Default configuration constants
const (
	DefaultAppTitle        = "emu8"
	DefaultAppFile         = ""
	DefaultEmulatorAsync   = false
	DefaultMachineModel    = "Speccy"
	DefaultMachineOptions  = ""
	DefaultVideoScale      = 2
	DefaultVideoFullScreen = false
	DefaultAudioFrecuency  = 48000 // 48 KHz
	DefaultAudioMute       = false
)

// -----------------------------------------------------------------------------
// Configuration
// -----------------------------------------------------------------------------

// Config is the main configuration
type Config struct {
	App      AppConfig
	Emulator EmulatorConfig
	Machine  MachineConfig
	Video    VideoConfig
	Audio    AudioConfig
}

// AppConfig is the application configuration
type AppConfig struct {
	Title string // Application base title
	File  string // File to load
}

// EmulatorConfig is the emulation configuration
type EmulatorConfig struct {
	Async bool // Async emulation
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
	Frequency int  // Audio frecuency
	Mute      bool // Mute autio
}

// -----------------------------------------------------------------------------
// Configuration Singleton
// -----------------------------------------------------------------------------

// config is the application main configuration
var config = new(Config)

// Get gets the main configuration
func Get() *Config { return config }

// init initializes configuration
func init() {
	config.App.Title = DefaultAppTitle
	config.App.File = DefaultAppFile
	config.Emulator.Async = DefaultEmulatorAsync
	config.Machine.Model = DefaultMachineModel
	config.Machine.Options = DefaultMachineOptions
	config.Video.Scale = DefaultVideoScale
	config.Video.FullScreen = DefaultVideoFullScreen
	config.Audio.Frequency = DefaultAudioFrecuency
	config.Audio.Mute = DefaultAudioMute
}
