// Package audio contains audio devices and components
package audio

import (
	"github.com/jtruco/emu8/emulator/device"
)

// -----------------------------------------------------------------------------
// Audio
// -----------------------------------------------------------------------------

// Audio device
type Audio interface {
	device.Device    // Is a device
	Buffer() *Buffer // Buffer returns the audio buffer
	Config() *Config // Config returns the audio configuration
	EndFrame()       // EndFrame audio frame is finished and buffer is ready
}

// Config is the device audio configuration
type Config struct {
	Frequency int     // Audio frequency
	Fps       int     // Frames per second
	Samples   int     // Samples per frame
	TStates   int     // Tstates per frame ( device frequency = tstates * fps )
	Rate      float32 // Device sampling rate
}

// NewConfig creates audio config
func NewConfig(frecuency, fps, tstates int) *Config {
	c := new(Config)
	c.Frequency = frecuency
	c.Fps = fps
	c.TStates = tstates
	c.Samples = c.Frequency / c.Fps
	c.Rate = float32(c.Samples) / float32(tstates)
	return c
}
