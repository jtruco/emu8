package format

import (
	"testing"

	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/io/tape"
)

// tapImage builds a TAP image from raw blocks (flag + data + checksum)
func tapImage(blocks ...[]byte) []byte {
	var data []byte
	for _, b := range blocks {
		data = append(data, byte(len(b)), byte(len(b)>>8))
		data = append(data, b...)
	}
	return data
}

// tapHeaderBlock is a standard 19 bytes program header block
var tapHeaderBlock = append([]byte{0x00, 0x00}, append([]byte("PROGRAM   "), 0x10, 0x00, 0x0a, 0x00, 0x10, 0x00, 0x00)...)

// playTape plays a tape until it stops for good (the end of the tape rewinds
// and stops it) or a step limit, resuming it a few times when a block stops
// the tape. A single step that never returns hangs the test.
func playTape(t *testing.T, tp tape.Tape) {
	t.Helper()
	const maxPlays = 4
	drive := tape.New(device.NewClock())
	drive.Insert(tp)
	for step, plays := 0, 0; step < 100000; step++ {
		if !drive.IsPlaying() {
			if plays == maxPlays {
				return
			}
			plays++
			drive.Play()
		}
		drive.Emulate(4096)
	}
}

func TestReadIntN(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	cases := []struct{ n, want int }{
		{1, 0x01}, {2, 0x0201}, {3, 0x030201}, {4, 0x04030201},
	}
	for _, c := range cases {
		if got := readIntN(data, 0, c.n); got != c.want {
			t.Errorf("readIntN(%d) = %#x, want %#x", c.n, got, c.want)
		}
	}
}

func TestTapLoadValidation(t *testing.T) {
	cases := []struct {
		name   string
		data   []byte
		valid  bool
		blocks int
	}{
		{"empty", nil, false, 0},
		{"one length byte", []byte{0x13}, false, 0},
		{"truncated block data", []byte{0x13, 0x00, 0x00, 0x00}, false, 0},
		{"zero length block", []byte{0x00, 0x00}, false, 0},
		{"truncated second block", append(tapImage(tapHeaderBlock), 0x05, 0x00, 0xff), false, 0},
		{"header and data", tapImage(tapHeaderBlock, []byte{0xff, 0x01, 0x02, 0x03}), true, 2},
		{"short header block", tapImage([]byte{0x00, 0x00, 0x41}), true, 1},
		{"one byte block", tapImage([]byte{0xff}), true, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tap := NewTap()
			if valid := tap.Load(c.data); valid != c.valid {
				t.Fatalf("Load = %v, want %v", valid, c.valid)
			}
			if c.valid {
				if n := len(tap.Blocks()); n != c.blocks {
					t.Errorf("%d blocks, want %d", n, c.blocks)
				}
				playTape(t, tap)
			}
		})
	}
}

func TestTzxLoadValidation(t *testing.T) {
	customInfo := append([]byte{0x35}, []byte("CUSTOM INFO     ")...)
	customInfo = append(customInfo, 0x2c, 0x01, 0x00, 0x00) // 300 bytes
	customInfo = append(customInfo, make([]byte, 300)...)
	cases := []struct {
		name   string
		data   []byte
		valid  bool
		blocks int
	}{
		{"empty", nil, false, 0},
		{"short signature", []byte("ZX"), false, 0},
		{"wrong signature", []byte("ZXTapo!\x1a\x01\x14"), false, 0},
		{"truncated header", []byte("ZXTape!\x1a\x01"), false, 0},
		{"header only", tzxImage(), true, 1},
		{"unknown block", tzxImage([]byte{0x99}), false, 0},
		{"truncated fixed part", tzxImage([]byte{0x11, 0x00, 0x00}), false, 0},
		{"truncated data", tzxImage([]byte{0x10, 0xe8, 0x03, 0x10, 0x00, 0xff}), false, 0},
		{"standard block", tzxImage([]byte{0x10, 0xe8, 0x03, 0x02, 0x00, 0xff, 0x00}), true, 2},
		{"4 bytes length (custom info)", tzxImage(customInfo, tzxPureToneBlock), true, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tzx := NewTzx()
			if valid := tzx.Load(c.data); valid != c.valid {
				t.Fatalf("Load = %v, want %v", valid, c.valid)
			}
			if c.valid {
				if n := len(tzx.Blocks()); n != c.blocks {
					t.Errorf("%d blocks, want %d", n, c.blocks)
				}
				playTape(t, tzx)
			}
		})
	}
}

// Valid TZX blocks with degenerate values must play without panics or hangs
func TestTzxPlayDegenerateBlocks(t *testing.T) {
	cases := []struct {
		name  string
		block []byte
	}{
		{"standard block without data", []byte{0x10, 0x01, 0x00, 0x00, 0x00}},
		{"turbo block without data", append([]byte{0x11}, make([]byte, 18)...)},
		{"pure data block without data", []byte{0x14, 0x10, 0x00, 0x20, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"last byte with 0 bits", []byte{0x14, 0x10, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0xaa}},
		{"last byte with 9 bits", []byte{0x14, 0x10, 0x00, 0x20, 0x00, 0x09, 0x00, 0x00, 0x01, 0x00, 0x00, 0xaa}},
		{"jump to itself", []byte{0x23, 0x00, 0x00}},
		{"loop end without start", []byte{0x25}},
		{"loop of 0", []byte{0x24, 0x00, 0x00, 0x25}},
		{"pulse sequence of 0", []byte{0x13, 0x00}},
		{"pure tone of 0 Ts", []byte{0x12, 0x00, 0x00, 0xff, 0xff}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tzx := NewTzx()
			if !tzx.Load(tzxImage(c.block, tzxPureToneBlock)) {
				t.Fatal("Load failed")
			}
			playTape(t, tzx)
		})
	}
}

func FuzzTapLoadPlay(f *testing.F) {
	f.Add(tapImage(tapHeaderBlock, []byte{0xff, 0x01, 0x02, 0x03}))
	f.Add(tapImage([]byte{0x00, 0x00, 0x41}))
	f.Add([]byte{0x13, 0x00, 0x00, 0x00})
	f.Fuzz(func(t *testing.T, data []byte) {
		tap := NewTap()
		if tap.Load(data) {
			playTape(t, tap)
		}
	})
}

func FuzzTzxLoadPlay(f *testing.F) {
	f.Add(tzxImage(tzxPureToneBlock, tzxStopBlock, tzxPureToneBlock))
	f.Add(tzxImage([]byte{0x10, 0xe8, 0x03, 0x02, 0x00, 0xff, 0x00}))
	f.Add(tzxImage([]byte{0x14, 0x10, 0x00, 0x20, 0x00, 0x03, 0x00, 0x00, 0x01, 0x00, 0x00, 0xaa}))
	f.Add(tzxImage([]byte{0x13, 0x02, 0x10, 0x00, 0x20, 0x00}, tzxStop48KBlock))
	f.Add(tzxImage([]byte{0x24, 0x02, 0x00}, tzxPureToneBlock, []byte{0x25}))
	f.Add(tzxImage([]byte{0x23, 0x02, 0x00}, tzxPauseBlock, tzxPureToneBlock))
	f.Fuzz(func(t *testing.T, data []byte) {
		tzx := NewTzx()
		if tzx.Load(data) {
			playTape(t, tzx)
		}
	})
}

// Jump to block uses a signed relative offset
func TestTzxJumpBackward(t *testing.T) {
	data := tzxImage(
		[]byte{0x23, 0x03, 0x00}, // #1 jump +3 -> #4
		tzxPureToneBlock,         // #2 played after the backward jump
		[]byte{0x23, 0x02, 0x00}, // #3 jump +2 -> end of tape
		[]byte{0x23, 0xfe, 0xff}, // #4 jump -2 -> #2
	)
	drive := newTestDrive(t, data)
	if toggles := playUntilStop(t, drive); toggles != tzxTestTonePulses {
		t.Errorf("%d EAR toggles, want the tone %d", toggles, tzxTestTonePulses)
	}
}

func TestTzxJumpBeforeFirstBlock(t *testing.T) {
	tzx := NewTzx().(*Tzx)
	if !tzx.Load(tzxImage([]byte{0x23, 0xf6, 0xff})) { // #1 jump -10
		t.Fatal("Load failed")
	}
	control := tape.Control{State: tapeStateTzxHeader, BlockIndex: 1, NumBlocks: len(tzx.Blocks())}
	tzx.Play(&control)
	if control.BlockIndex != 0 {
		t.Errorf("block index = %d, want 0", control.BlockIndex)
	}
}

func TestReadString(t *testing.T) {
	data := []byte("A\x01Z\x7f\x80~ ")
	if got := readString(data, 0, len(data)); got != "A?Z\x7f?~ " {
		t.Errorf("readString = %q", got)
	}
}
