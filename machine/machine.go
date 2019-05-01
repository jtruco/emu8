// Package machine contains 8bit machine and coponents implementation
package machine

import (
	"github.com/jtruco/emu8/cpu"
	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/audio"
	"github.com/jtruco/emu8/device/io/keyboard"
	"github.com/jtruco/emu8/device/video"
)

// -----------------------------------------------------------------------------
// Machine
// -----------------------------------------------------------------------------

// Machine is a 8bit machine
type Machine interface {
	device.Device // Is a device
	// Config gets the machine configuration
	Config() *Config
	// CPU the machine main CPU
	CPU() cpu.CPU
	// Components the machine components
	Components() *device.Components
	// SetController connects the machine to the controller
	SetController(Controller)
	// LoadFile
	LoadFile(filename string)
	// BeginFrame begin emulation frame tasks
	BeginFrame()
	// Emulate one machine step
	Emulate()
	// EndFrame end emulation frame tasks
	EndFrame()
}

// Controller is the machine controller
type Controller interface {
	Keyboard() *keyboard.Controller
	Video() *video.Controller
	Audio() *audio.Controller
}

// -----------------------------------------------------------------------------
// Machine configuration
// -----------------------------------------------------------------------------

// Config machine configuration
type Config struct {
	Model         int     // Machine model
	FPS           float32 // Frames per second
	FrameDuration int     // Duration of a frame in Nanos (= 1e9 / FPS)
	FrameTStates  int     // TStates per frame
}

// SetFPS sets FPS and FrameDuration
func (conf *Config) SetFPS(FPS float32) {
	conf.FPS = FPS
	conf.FrameDuration = int(float32(1e9) / FPS)
}
