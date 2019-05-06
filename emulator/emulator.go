// Package emulator implements the core of 8bit machine emulator
package emulator

import (
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
	machine  machine.Machine                // The emulated machine
	file     *controller.FileManager        // The files manager
	video    *controller.VideoController    // The video controller
	audio    *controller.AudioController    // The audio controller
	keyboard *controller.KeyboardController // The keyboard controlller
	running  bool
	wg       sync.WaitGroup
}

// New creates a machine emulator
func New(machine machine.Machine) *Emulator {
	emulator := &Emulator{}
	emulator.machine = machine
	emulator.file = controller.DefaultFileManager()
	emulator.video = controller.NewVideoController()
	emulator.audio = controller.NewAudioController()
	emulator.keyboard = controller.NewKeyboardController()
	emulator.machine.SetController(emulator)
	return emulator
}

// Machine controller

// Machine gets the emulated machine
func (emulator *Emulator) Machine() machine.Machine {
	return emulator.machine
}

// FileManager the filesystem manager
func (emulator *Emulator) File() *controller.FileManager {
	return emulator.file
}

// Video the video controller
func (emulator *Emulator) Video() *controller.VideoController {
	return emulator.video
}

// Audio the audio controller
func (emulator *Emulator) Audio() *controller.AudioController {
	return emulator.audio
}

// Keyboard the keyboard controller
func (emulator *Emulator) Keyboard() *controller.KeyboardController {
	return emulator.keyboard
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
	}
}

// Stop the emulation
func (emulator *Emulator) Stop() {
	if emulator.running {
		emulator.running = false
		emulator.wg.Wait()
	}
}

// Emulation

// runEmulation the emulation loop goroutine
func (emulator *Emulator) runEmulation() {
	// sync
	emulator.wg.Add(1)
	defer emulator.wg.Done()

	// emulation speed
	config := emulator.machine.Config()
	ticker := time.NewTicker(time.Duration(config.FrameDuration))
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
	// FIXME : emulator.keyboard.Flush()
}

// refreshUI refresh UI asynchronusly
func (emulator *Emulator) refreshUI() {
	// Video & Audio refresh
	emulator.video.Refresh()
	emulator.audio.Flush()
}
