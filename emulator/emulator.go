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
	wg         sync.WaitGroup        // Sync control
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

// Init the emulation
func (emulator *Emulator) Init() { emulator.machine.Init() }

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
	emulator.emulateFrame()
}

// Start the emulation
func (emulator *Emulator) Start() {
	if !emulator.running {
		emulator.running = true
		go emulator.runEmulation()
		log.Println("Emulator : Started")
	}
}

// Stop the emulation
func (emulator *Emulator) Stop() {
	if emulator.running {
		emulator.running = false
		emulator.wg.Wait()
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

// runEmulation the emulation loop goroutine
func (emulator *Emulator) runEmulation() {
	emulator.wg.Add(1)
	defer emulator.wg.Done()

	// emulation timmings
	frameTime := emulator.machine.Config().FrameTime
	sleep := frameTime

	// emulation loop
	for emulator.running {
		start := time.Now()
		emulator.emulateFrame()
		time.Sleep(sleep) // sleep until next frame
		sleep += frameTime - time.Since(start)
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
