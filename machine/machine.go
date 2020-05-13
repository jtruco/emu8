// Package machine contains 8bit machine and coponents implementation
package machine

import (
	"time"

	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/cpu"
	"github.com/jtruco/emu8/emulator/controller"
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
	SetController(controller.Controller)
	// BeginFrame begin emulation frame tasks
	BeginFrame()
	// Emulate one machine step
	Emulate()
	// EndFrame end emulation frame tasks
	EndFrame()
	// LoadFile loads a file into machine
	LoadFile(name string)
}

// -----------------------------------------------------------------------------
// Machine configuration
// -----------------------------------------------------------------------------

// Config machine configuration
type Config struct {
	Model        int           // Machine model
	FPS          float32       // Frames per second
	FrameTime    time.Duration // Duration of a frame
	FrameTStates int           // TStates per frame
}

// SetFPS sets FPS and FrameDuration
func (conf *Config) SetFPS(FPS float32) {
	conf.FPS = FPS
	conf.FrameTime = time.Duration(1e9 / FPS)
}
