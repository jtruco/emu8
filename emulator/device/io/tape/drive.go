package tape

import (
	"log"

	"github.com/jtruco/emu8/emulator/device"
)

// -----------------------------------------------------------------------------
// Tape Drive
// -----------------------------------------------------------------------------

// Drive tape device
type Drive struct {
	control Control      // Tape control data
	clock   device.Clock // Clock
	tape    Tape         // Loaded tape
}

// New creates a new Tape Drive
func New(clock device.Clock) *Drive {
	drive := new(Drive)
	drive.clock = clock
	drive.Reset()
	return drive
}

// Init intializes tape drive
func (drive *Drive) Init() { drive.Reset() }

// Reset resets tape drive
func (drive *Drive) Reset() {
	drive.Stop()
	drive.control.reset()
}

// HasTape if there is a tape
func (drive *Drive) HasTape() bool { return drive.tape != nil }

// IsPlaying if tape drive is playing
func (drive *Drive) IsPlaying() bool { return drive.control.Playing }

// Ear tape value
func (drive *Drive) Ear() byte { return drive.control.Ear }

// EarHigh tape state is high
func (drive *Drive) EarHigh() bool { return (drive.control.Ear & LevelMask) != 0 }

// EarLow tape state is high
func (drive *Drive) EarLow() bool { return (drive.control.Ear & LevelMask) == 0 }

// Insert loads the tape into the drive
func (drive *Drive) Insert(tape Tape) {
	drive.tape = tape
	drive.control.NumBlocks = len(tape.Blocks())
	drive.Reset()
	log.Println("Tape : Tape inserted:", tape.Info().Name)
}

// Eject ejects the tape from drive
func (drive *Drive) Eject() {
	drive.tape = nil
	drive.control.NumBlocks = 0
	drive.Reset()
	log.Println("Tape : Tape ejected")
}

// Play starts tape playback
func (drive *Drive) Play() {
	if drive.IsPlaying() || !drive.HasTape() {
		return
	}
	drive.control.Playing = true
	log.Println("Tape : Playback started")
}

// Stop stops tape playback
func (drive *Drive) Stop() {
	if !drive.IsPlaying() {
		return
	}
	drive.control.Playing = false
	log.Println("Tape : Playback stopped")
}

// Rewind rewinds the tape to start
func (drive *Drive) Rewind() {
	drive.Reset()
	log.Println("Tape : Tape rewinded")
}

// Emulate emulates the tape drive
func (drive *Drive) Emulate(tstates int) {
	if !drive.IsPlaying() {
		return
	}
	// control tstates timeout
	drive.control.Timeout -= tstates
	if drive.control.Timeout > 0 {
		return
	}
	padding := drive.control.Timeout
	// play until next state
	drive.control.Timeout = 0
	for drive.IsPlaying() && drive.control.Timeout == 0 {
		drive.tape.Play(&drive.control)
	}
	drive.control.Timeout += padding
	// control end of tape playback
	if !drive.IsPlaying() {
		if drive.control.EndOfTape() {
			log.Println("Tape : End of tape")
			drive.Rewind()
		} else {
			log.Println("Tape : End of playback")
		}
	}
}
