package format

import (
	"testing"

	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/io/tape"
)

// even number of pulses : the EAR ends low, so end of tape rewinds and pauses
// (which force the EAR low) don't hide the last toggle
const tzxTestTonePulses = 4

// tzxImage builds a TZX v1.20 image from raw blocks
func tzxImage(blocks ...[]byte) []byte {
	data := []byte("ZXTape!\x1a\x01\x14")
	for _, b := range blocks {
		data = append(data, b...)
	}
	return data
}

var (
	tzxStopBlock     = []byte{0x20, 0x00, 0x00}                          // pause 0 : stop the tape
	tzxPauseBlock    = []byte{0x20, 0x01, 0x00}                          // pause 1 ms
	tzxStop48KBlock  = []byte{0x2a, 0x00, 0x00, 0x00, 0x00}              // stop the tape if in 48K mode
	tzxPureToneBlock = []byte{0x12, 0x10, 0x00, tzxTestTonePulses, 0x00} // pulses of 16 Ts
)

func newTestDrive(t *testing.T, data []byte) *tape.Drive {
	t.Helper()
	tzx := NewTzx()
	if !tzx.Load(data) {
		t.Fatal("TZX load failed")
	}
	drive := tape.New(device.NewClock())
	drive.Insert(tzx)
	return drive
}

// playUntilStop plays the tape until it stops and returns the EAR toggles
func playUntilStop(t *testing.T, drive *tape.Drive) int {
	t.Helper()
	drive.Play()
	ear, toggles := drive.Ear(), 0
	for i := 0; i < 10000 && drive.IsPlaying(); i++ {
		drive.Emulate(8)
		if drive.Ear() != ear {
			ear = drive.Ear()
			toggles++
		}
	}
	if drive.IsPlaying() {
		t.Fatal("tape did not stop")
	}
	return toggles
}

func TestTzxResumeAfterStopTape(t *testing.T) {
	cases := []struct {
		name string
		stop []byte
	}{
		{"pause 0 block", tzxStopBlock},
		{"stop if 48K block", tzxStop48KBlock},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			drive := newTestDrive(t, tzxImage(tzxPureToneBlock, c.stop, tzxPureToneBlock))
			if toggles := playUntilStop(t, drive); toggles != tzxTestTonePulses {
				t.Errorf("first part: %d EAR toggles, want %d", toggles, tzxTestTonePulses)
			}
			// user presses play again (multi-load) : the next block plays
			if toggles := playUntilStop(t, drive); toggles != tzxTestTonePulses {
				t.Errorf("after resume: %d EAR toggles, want %d", toggles, tzxTestTonePulses)
			}
		})
	}
}

func TestTzxStopTapeAsLastBlock(t *testing.T) {
	drive := newTestDrive(t, tzxImage(tzxPureToneBlock, tzxStopBlock))
	playUntilStop(t, drive)
	// the stop block was the last one : the tape is rewound and plays again
	if toggles := playUntilStop(t, drive); toggles != tzxTestTonePulses {
		t.Errorf("after the end of the tape: %d EAR toggles, want %d", toggles, tzxTestTonePulses)
	}
}

func TestTzxPauseDoesNotStop(t *testing.T) {
	drive := newTestDrive(t, tzxImage(tzxPureToneBlock, tzxPauseBlock, tzxPureToneBlock))
	if toggles := playUntilStop(t, drive); toggles != 2*tzxTestTonePulses {
		t.Errorf("%d EAR toggles, want both tones %d without stopping", toggles, 2*tzxTestTonePulses)
	}
}

func TestTzxUserStopResumesInBlock(t *testing.T) {
	drive := newTestDrive(t, tzxImage(tzxPureToneBlock))
	drive.Play()
	ear, toggles := drive.Ear(), 0
	for toggles < 2 {
		drive.Emulate(8)
		if drive.Ear() != ear {
			ear = drive.Ear()
			toggles++
		}
	}
	drive.Stop() // user stops (F7) in the middle of the block
	if rest := playUntilStop(t, drive); toggles+rest != tzxTestTonePulses {
		t.Errorf("%d + %d EAR toggles, want %d", toggles, rest, tzxTestTonePulses)
	}
}
