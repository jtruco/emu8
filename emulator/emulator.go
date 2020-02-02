// Package emulator implements the core of 8bit machine emulator
package emulator

import (
	"log"
	"sync"
	"time"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/machine"
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
	emulator := &Emulator{}
	emulator.machine = machine
	emulator.controller = controller.New()
	emulator.machine.SetController(emulator.controller)
	return emulator
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
func (emulator *Emulator) Reset() { emulator.machine.Reset() }

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

// Emulation

// runEmulation the emulation loop goroutine
func (emulator *Emulator) runEmulation() {
	// sync
	emulator.wg.Add(1)
	defer emulator.wg.Done()

	// emulation speed
	frame := emulator.machine.Config().FrameTime
	ticker := time.NewTicker(time.Duration(frame))
	defer ticker.Stop()

	// emulation loop
	for emulator.running {
		select {
		case <-ticker.C:
			{
				// do frame
				emulator.flushInput()
				emulator.emulateFrame()
				emulator.refreshUI()
			}
		}
	}
}

// emulateFrame emulates the frame
func (emulator *Emulator) emulateFrame() {
	machine := emulator.machine
	clock := machine.CPU().Clock()
	config := machine.Config()

	// frame emulation loop
	clock.Restart(config.FrameTStates)
	machine.BeginFrame()
	for clock.Tstates() < config.FrameTStates {
		machine.Emulate()
	}
	machine.EndFrame()
}

// flushInput flushes input events
func (emulator *Emulator) flushInput() {
	// Keyboard events
	emulator.controller.Keyboard().Flush()
}

// refreshUI refresh UI asynchronusly
func (emulator *Emulator) refreshUI() {
	// Video & Audio refresh
	emulator.controller.Audio().Flush()
	emulator.controller.Video().Refresh()
}
