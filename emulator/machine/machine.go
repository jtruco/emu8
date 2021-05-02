// Package machine contains 8bit machines and coponents
package machine

import (
	"time"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/device"
)

// -----------------------------------------------------------------------------
// Machine
// -----------------------------------------------------------------------------

// Machine is a 8bit machine
type Machine interface {
	device.Device // Is a device
	// Config gets the machine configuration
	Config() *Config
	// Clock the machine main clock
	Clock() device.Clock
	// Components the machine components
	Components() *device.Components
	// InitControl connects the machine to the emulator controller
	InitControl(*controller.Controller)
	// BeginFrame begin emulation frame tasks
	BeginFrame()
	// Emulate one machine step
	Emulate()
	// EndFrame end emulation frame tasks
	EndFrame()
	// LoadFile loads a file into machine
	LoadFile(name string)
	// TakeSnap takes and saves snapshop of the machine state
	TakeSnapshot()
}

// Config machine configuration
type Config struct {
	Model        int           // Machine model
	Name         string        // Model name
	FPS          float32       // Frames per second
	FrameTime    time.Duration // Duration of a frame
	FrameTStates int           // TStates per frame
}

// SetFPS sets FPS and FrameDuration
func (conf *Config) SetFPS(FPS float32) {
	conf.FPS = FPS
	conf.FrameTime = time.Duration(1e9 / FPS)
}
