package config

// Default configuration constants
const (
	DefaultAppTitle     = "emu8"
	DefaultMachineModel = "Speccy"
)

// Config is the main configuration
type Config struct {
	AppTitle       string
	MachineModel   string
	MachineOptions string
	FileName       string
	Video          VideoConfig
	Audio          AudioConfig
}

// config is the application main configuration
var config = new(Config)

// Get gets the main configuration
func Get() *Config {
	return config
}

// Default configuration constants
const (
	DefaultVideoScale      = 2
	DefaultVideoFullScreen = false
)

// VideoConfig the video configuration
type VideoConfig struct {
	Scale      int
	FullScreen bool
}

// Default configuration constants
const (
	DefaultAudioFrecuency = 48000 // 48 KHz
	DefaultAudioMute      = false
)

// AudioConfig the audio configuration
type AudioConfig struct {
	Frequency int
	Mute      bool
}

// init initializes configuration
func init() {
	config.AppTitle = DefaultAppTitle
	config.MachineModel = DefaultMachineModel
	config.Video.Scale = DefaultVideoScale
	config.Video.FullScreen = DefaultVideoFullScreen
	config.Audio.Frequency = DefaultAudioFrecuency
	config.Audio.Mute = DefaultAudioMute
}
