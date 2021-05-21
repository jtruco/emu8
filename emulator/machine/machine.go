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
	device.Device                       // Is a device
	Config() *Config                    // Config gets the machine configuration
	Clock() device.Clock                // Clock the machine main clock
	Components() *device.Components     // Components the machine components
	InitControl(*controller.Controller) // InitControl connects the machine to the emulator controller
	Emulate()                           // Emulate one machine step
	BeginFrame()                        // BeginFrame begin emulation frame tasks
	EndFrame()                          // EndFrame end emulation frame tasks
	LoadFile(name string)               // LoadFile loads a file into machine
	TakeSnapshot()                      // TakeSnap takes and saves snapshop of the machine state
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
