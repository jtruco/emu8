// Package emulator implements the core of 8bit machine emulator
package emulator

import (
	"log"
	"time"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// -----------------------------------------------------------------------------
// Emulator
// -----------------------------------------------------------------------------

// Emulator is the emulator main controller. The UI drives the emulation
// calling Emulate and Sync once per frame.
type Emulator struct {
	machine  machine.Machine        // The hosted machine
	control  *controller.Controller // The emulator controller
	running  bool                   // Indicates emulation is running
	tstates  int                    // Tstates per frame
	duration time.Duration          // Frame duration
	sleep    time.Duration          // Sleep duration
	current  time.Time              // Current time
	lost     bool                   // Lost frame
}

// New creates a ready to use emulator of a machine, using a file system for
// its files. It connects the machine to the emulator controller, then
// initializes it (ROMs are loaded).
func New(machine machine.Machine, fs vfs.FileSystem) (*Emulator, error) {
	emulator := new(Emulator)
	emulator.machine = machine
	emulator.control = controller.New(machine, fs)
	if err := emulator.initialize(); err != nil {
		return nil, err
	}
	return emulator, nil
}

// Emulator factory

// Create returns an emulator of a machine built from its settings, using a
// file system for its files
func Create(settings machine.Settings, fs vfs.FileSystem) (*Emulator, error) {
	machine, err := machine.Create(settings)
	if err != nil {
		return nil, err
	}
	return New(machine, fs)
}

// Machine controller

// Control gets the emulator controller
func (emulator *Emulator) Control() *controller.Controller {
	return emulator.control
}

// Machine gets the hosted machine
func (emulator *Emulator) Machine() machine.Machine {
	return emulator.machine
}

// Machine emulation control

// IsRunning the emulation
func (emulator *Emulator) IsRunning() bool { return emulator.running }

// initialize initializes the machine and the emulation timings
func (emulator *Emulator) initialize() error {
	if err := emulator.machine.Init(); err != nil {
		return err
	}
	emulator.tstates = emulator.machine.Config().TStates
	emulator.duration = emulator.machine.Config().Duration
	return nil
}

// Reset the emulation
func (emulator *Emulator) Reset() {
	if emulator.running {
		emulator.Stop()
		defer emulator.Start()
	}
	emulator.machine.Reset()
}

// Start the emulation
func (emulator *Emulator) Start() {
	if !emulator.running {
		emulator.running = true
		emulator.sleep = emulator.duration
		emulator.current = time.Now()
		log.Println("Emulator : Started")
	}
}

// Stop the emulation
func (emulator *Emulator) Stop() {
	if emulator.running {
		emulator.running = false
		log.Println("Emulator : Stopped")
	}
}

// LoadFile loads file into the emulator
func (emulator *Emulator) LoadFile(name string) {
	if emulator.running {
		emulator.Stop()
		defer emulator.Start()
	}
	emulator.control.LoadFile(name)
}

// TakeSnapshot takes and saves snapshot of the machine state
func (emulator *Emulator) TakeSnapshot() {
	if emulator.running {
		emulator.Stop()
		defer emulator.Start()
	}
	emulator.control.TakeSnapshot()
}

// Emulation

// Emulate one frame loop
func (emulator *Emulator) Emulate() {
	if emulator.running {
		emulator.emulateFrame()
	}
}

// Sync synchronizes next frame loop
func (emulator *Emulator) Sync() {
	if emulator.running {
		emulator.syncFrame()
	}
}

// emulateFrame emulates the frame
func (emulator *Emulator) emulateFrame() {
	// first we process input events and begin the frame
	emulator.control.Scan()
	emulator.machine.BeginFrame()

	// the clock is restarted every frame to adjust the emulation rate
	clock := emulator.machine.Clock()
	clock.Restart(emulator.tstates)

	// machine frame-cycle emulation
	for clock.Tstates() < emulator.tstates {
		emulator.machine.Emulate()
	}

	// last we finish the frame and update the output
	emulator.machine.EndFrame()
	emulator.control.Refresh()
}

// syncFrame adjusts the frame timing
func (emulator *Emulator) syncFrame() {
	emulator.sleep += emulator.duration - time.Since(emulator.current)
	emulator.current = time.Now()
	if emulator.sleep > 0 {
		emulator.lost = false
		time.Sleep(emulator.sleep) // sleep until next frame
	} else {
		emulator.lost = true // lost frame
		emulator.sleep = 0   // reset sleep control
	}
}
