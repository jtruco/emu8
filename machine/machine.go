// Package machine contains 8bit machine and coponents implementation
package machine

import (
	"errors"
	"time"

	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/cpu"
	"github.com/jtruco/emu8/emulator/controller"
)

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
	// TakeSnap takes and saves snapshop of the machine state
	TakeSnapshot()
}

// -----------------------------------------------------------------------------
// Machine Factory
// -----------------------------------------------------------------------------

// Factory is a machine constructor function
type Factory func(int) Machine

// Registered machine factories
var factories = map[int]Factory{}

// Register registers a machine
func Register(id int, factory Factory) {
	factories[id] = factory
}

// Create returns a machine from a model name
func Create(model string) (Machine, error) {
	return CreateFromModel(GetModel(model))
}

// CreateFromModel returns a machine from a model id
func CreateFromModel(model int) (Machine, error) {
	id := GetMachineFromModel(model)
	if id == UnknownMachine {
		return nil, errors.New("Machine : unknown machine model")
	}
	factory := factories[id]
	if factory == nil {
		return nil, errors.New("Machine : unsupported machine model")
	}
	return factory(model), nil
}
