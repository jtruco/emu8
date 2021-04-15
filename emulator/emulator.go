// Package emulator implements the core of 8bit machine emulator
package emulator

import (
	"log"
	"sync"
	"time"

	"github.com/jtruco/emu8/emulator/config"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/machine"
)

// -----------------------------------------------------------------------------
// Emulator
// -----------------------------------------------------------------------------

// Emulator is the emulator main controller
type Emulator struct {
	machine    machine.Machine       // The hosted machine
	controller controller.Controller // The emulator controller
	running    bool                  // Indicates emulation is running
	async      bool                  // Async emulation goroutine
	wg         sync.WaitGroup        // Sync control
	frame      time.Duration         // Frame duration
	sleep      time.Duration         // Sleep duration
	current    time.Time             // Current time
}

// New creates a machine emulator
func New(machine machine.Machine) *Emulator {
	emulator := new(Emulator)
	emulator.machine = machine
	emulator.controller = controller.New()
	emulator.machine.SetController(emulator.controller)
	return emulator
}

// Emulator factory

// GetDefault returns the configured emulator
func GetDefault() *Emulator {
	return FromModel(config.Get().MachineModel)
}

// FromModel returns an emulator for a machine model name
func FromModel(model string) *Emulator {
	machine, err := machine.Create(model)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return New(machine)
}

// Machine controller

// Controller gets the emulator controller
func (emulator *Emulator) Controller() controller.Controller {
	return emulator.controller
}

// Machine gets the hosted machine
func (emulator *Emulator) Machine() machine.Machine {
	return emulator.machine
}

// Machine emulation control

// IsRunning the emulation
func (emulator *Emulator) IsRunning() bool { return emulator.running }

// IsAsync async emulation is active
func (emulator *Emulator) IsAsync() bool { return emulator.async }

// IsAsync emulation loop is active
func (emulator *Emulator) SetAsync(async bool) {
	if emulator.running {
		emulator.Stop()
	}
	emulator.async = async
}

// Init the emulation
func (emulator *Emulator) Init() {
	emulator.machine.Init()
	emulator.frame = emulator.machine.Config().FrameTime
}

// Reset the emulation
func (emulator *Emulator) Reset() {
	if emulator.running {
		emulator.Stop()
		defer emulator.Start()
	}
	emulator.machine.Reset()
}

// Emulate one frame loop
func (emulator *Emulator) Emulate() {
	if !emulator.running {
		return
	}
	emulator.emulateFrame()
}

// Sync synchronizes next frame loop
func (emulator *Emulator) Sync() {
	time.Sleep(emulator.sleep) // sleep until next frame
	emulator.sleep += emulator.frame - time.Since(emulator.current)
	emulator.current = time.Now()
}

// Start the emulation
func (emulator *Emulator) Start() {
	if !emulator.running {
		emulator.running = true
		emulator.sleep = emulator.frame
		emulator.current = time.Now()
		if emulator.async {
			go emulator.emulationLoop()
		}
		log.Println("Emulator : Started")
	}
}

// Stop the emulation
func (emulator *Emulator) Stop() {
	if emulator.running {
		emulator.running = false
		if emulator.async {
			emulator.wg.Wait()
		}
		log.Println("Emulator : Stopped")
	}
}

// LoadFile loads file into the emulator
func (emulator *Emulator) LoadFile(name string) {
	if emulator.running {
		emulator.Stop()
		defer emulator.Start()
	}
	emulator.machine.LoadFile(name)
}

// TakeSnapshot takes and saves snapshop of the machine state
func (emulator *Emulator) TakeSnapshot() {
	if emulator.running {
		emulator.Stop()
		defer emulator.Start()
	}
	emulator.machine.TakeSnapshot()
}

// Emulation

// emulationLoop the emulation loop goroutine
func (emulator *Emulator) emulationLoop() {
	emulator.wg.Add(1)
	defer emulator.wg.Done()

	// emulation loop
	for emulator.running {
		emulator.emulateFrame()
		emulator.Sync()
	}
}

// emulateFrame emulates the frame
func (emulator *Emulator) emulateFrame() {
	machine := emulator.machine
	clock := machine.CPU().Clock()
	config := machine.Config()

	// pre-frame actions
	emulator.controller.Scan()

	// frame emulation loop
	machine.BeginFrame()
	for clock.Tstates() < config.FrameTStates {
		machine.Emulate()
	}
	machine.EndFrame()

	// post-frame actions
	emulator.controller.Refresh()
	clock.Restart(config.FrameTStates)
}
