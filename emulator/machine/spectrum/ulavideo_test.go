package spectrum

import (
	"testing"

	"github.com/jtruco/emu8/emulator/machine"
	"github.com/jtruco/emu8/emulator/machine/spectrum/format"
)

// Reference positions (48K): the first visible pixel (0,0) is drawn at
// Tstate 8944 (line 40, 16 Ts of left border), the first paper pixel at
// Tstate 14336, shown at screen (32,24).

func TestUlaTvGeometry(t *testing.T) {
	cases := []struct {
		name      string
		got, want int
	}{
		{"frame tstates", tvFrameTstates, zxTStates},
		{"frame tstates value", tvFrameTstates, 69888},
		{"first paper tstate", tvScreenTstate, 14336},
		{"screen width", tvScreenWidth, 320},
		{"screen height", tvScreenHeight, 240},
		{"block tstates", tvBlockTstates, 4},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}
	m := newTestSpectrum(t)
	if w, h := m.tv.Screen().Width(), m.tv.Screen().Height(); w != 320 || h != 240 {
		t.Errorf("screen size = %dx%d, want 320x240", w, h)
	}
}

func TestUlaTvPaperAddressTable(t *testing.T) {
	cases := []struct {
		y          int
		data, attr int
	}{
		{0, 0x0000, 0x1800},
		{1, 0x0100, 0x1800},
		{7, 0x0700, 0x1800},
		{8, 0x0020, 0x1820},
		{63, 0x07e0, 0x18e0},
		{64, 0x0800, 0x1900},
		{191, 0x17e0, 0x1ae0},
	}
	for _, c := range cases {
		if got := tvPaperAddr[c.y]; got[0] != c.data || got[1] != c.attr {
			t.Errorf("tvPaperAddr[%d] = %#04x/%#04x, want %#04x/%#04x",
				c.y, got[0], got[1], c.data, c.attr)
		}
	}
}

func TestUlaTvReset(t *testing.T) {
	m := newTestSpectrum(t)
	m.at(30000).tv.SetBorder(3)
	m.tv.flash = true
	m.tv.Reset()

	if m.tv.tstate != 0 {
		t.Errorf("tstate = %d, want 0", m.tv.tstate)
	}
	if m.tv.border != 7 || m.tv.borderRgb != colWhite {
		t.Errorf("border = %d/%#x, want 7/%#x", m.tv.border, m.tv.borderRgb, colWhite)
	}
	if m.tv.flash {
		t.Error("flash = true, want false")
	}
	if p := m.pixel(0, 0); p != colBlack {
		t.Errorf("pixel(0,0) = %#x, want cleared %#x", p, colBlack)
	}
}

func TestUlaTvFullFrame(t *testing.T) {
	m := newTestSpectrum(t)
	m.tv.SetBorder(2)
	m.tv.EndFrame()

	cases := []struct {
		name string
		x, y int
		want uint32
	}{
		{"top-left border", 0, 0, colRed},
		{"top-right border", 319, 0, colRed},
		{"left border", 31, 24, colRed},
		{"first paper pixel", 32, 24, colBlack},
		{"last paper pixel", 287, 215, colBlack},
		{"right border", 288, 215, colRed},
		{"bottom-right border", 319, 239, colRed},
	}
	for _, c := range cases {
		if p := m.pixel(c.x, c.y); p != c.want {
			t.Errorf("%s pixel(%d,%d) = %#x, want %#x", c.name, c.x, c.y, p, c.want)
		}
	}
}

func TestUlaTvPaperPixels(t *testing.T) {
	cases := []struct {
		name   string
		bitmap byte
		attr   byte
		flash  bool
		want   [8]uint32
	}{
		{"ink and paper", 0xa5, 5<<3 | 2, false,
			[8]uint32{colRed, colCyan, colRed, colCyan, colCyan, colRed, colCyan, colRed}},
		{"bright", 0xf0, 0x40 | 2<<3 | 4, false,
			[8]uint32{colBrGreen, colBrGreen, colBrGreen, colBrGreen, colBrRed, colBrRed, colBrRed, colBrRed}},
		{"flash off", 0xf0, 0x80 | 5<<3 | 2, false,
			[8]uint32{colRed, colRed, colRed, colRed, colCyan, colCyan, colCyan, colCyan}},
		{"flash on swaps ink and paper", 0xf0, 0x80 | 5<<3 | 2, true,
			[8]uint32{colCyan, colCyan, colCyan, colCyan, colRed, colRed, colRed, colRed}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.videoPoke(0x0000, c.bitmap)
			m.videoPoke(0x1800, c.attr)
			m.tv.flash = c.flash
			m.tv.EndFrame()
			for i, want := range c.want {
				if p := m.pixel(32+i, 24); p != want {
					t.Errorf("pixel %d = %#x, want %#x", i, p, want)
				}
			}
		})
	}
}

func TestUlaTvPaperLayout(t *testing.T) {
	// paper line 65, column 3 : bitmap 0x0903, attribute row 8 0x1903
	m := newTestSpectrum(t)
	m.videoPoke(0x0903, 0x80)
	m.videoPoke(0x1903, 2) // ink red, paper black
	m.tv.EndFrame()

	x, y := 32+3*8, 24+65
	if p := m.pixel(x, y); p != colRed {
		t.Errorf("pixel(%d,%d) = %#x, want ink %#x", x, y, p, colRed)
	}
	if p := m.pixel(x+1, y); p != colBlack {
		t.Errorf("pixel(%d,%d) = %#x, want paper %#x", x+1, y, p, colBlack)
	}
	if p := m.pixel(x, y-1); p != colBlack {
		t.Errorf("pixel(%d,%d) = %#x, want paper %#x (line 64)", x, y-1, p, colBlack)
	}
}

// A border change at T is visible from the 4 Ts block containing T (FUSE).
func TestUlaTvBorderChangePosition(t *testing.T) {
	cases := []struct {
		name    string
		tstate  int
		x, y    int // first pixel with the new colour
		oldX    int // last border pixel with the old colour (-1 : none)
		oldY    int
		visible bool
	}{
		{"before the visible area", 100, 0, 0, -1, 0, true},
		{"first visible block", 8944, 0, 0, -1, 0, true},
		{"inside first block rounds down", 8945, 0, 0, -1, 0, true},
		{"end of first block rounds down", 8947, 0, 0, -1, 0, true},
		{"second block", 8948, 8, 0, 7, 0, true},
		{"inside second block", 8951, 8, 0, 7, 0, true},
		{"left border of first paper line", 14320, 0, 24, 319, 23, true},
		{"during paper shows at right border", 14336, 288, 24, 31, 24, true},
		{"right border of first paper line", 14464, 288, 24, 31, 24, true},
		{"during retrace shows on next line", 14492, 0, 25, 319, 24, true},
		{"after the visible area", 69000, 0, 0, -1, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.at(c.tstate).tv.SetBorder(2)
			m.tv.EndFrame()
			x, y, found := m.firstPixel(colRed)
			if found != c.visible {
				t.Fatalf("change visible = %v, want %v", found, c.visible)
			}
			if !found {
				return
			}
			if x != c.x || y != c.y {
				t.Errorf("first changed pixel = (%d,%d), want (%d,%d)", x, y, c.x, c.y)
			}
			if c.oldX >= 0 && m.pixel(c.oldX, c.oldY) != colWhite {
				t.Errorf("pixel(%d,%d) = %#x, want old border %#x",
					c.oldX, c.oldY, m.pixel(c.oldX, c.oldY), colWhite)
			}
		})
	}
}

func TestUlaTvUpdateRendersCompletedBlocks(t *testing.T) {
	cases := []struct {
		name   string
		update int
		cell0  uint32
	}{
		{"block not completed keeps it pending", 14339, colRed},
		{"block completed is rendered", 14340, colBlue},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.videoPoke(0x1800, 1<<3) // blue paper
			m.videoPoke(0x1801, 1<<3)
			m.tv.Update(c.update)
			m.videoPoke(0x1800, 2<<3) // red paper
			m.videoPoke(0x1801, 2<<3)
			m.tv.EndFrame()
			if p := m.pixel(32, 24); p != c.cell0 {
				t.Errorf("cell 0 line 64 = %#x, want %#x", p, c.cell0)
			}
			if p := m.pixel(40, 24); p != colRed {
				t.Errorf("cell 1 line 64 = %#x, want new %#x", p, colRed)
			}
			if p := m.pixel(32, 25); p != colRed {
				t.Errorf("cell 0 line 65 = %#x, want new %#x", p, colRed)
			}
		})
	}
}

func TestUlaTvUpdateBeamPosition(t *testing.T) {
	cases := []struct {
		name    string
		updates []int
		want    int
	}{
		{"rounds down to block", []int{14343}, 14340},
		{"does not rewind", []int{20000, 10000}, 20000},
		{"clamps to frame end", []int{tvFrameTstates + 100}, tvFrameTstates},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			for _, u := range c.updates {
				m.tv.Update(u)
			}
			if m.tv.tstate != c.want {
				t.Errorf("tstate = %d, want %d", m.tv.tstate, c.want)
			}
		})
	}
}

func TestUlaTvEndFrame(t *testing.T) {
	m := newTestSpectrum(t)
	m.tv.Update(30000)
	m.tv.EndFrame()
	if m.tv.tstate != 0 {
		t.Errorf("tstate after EndFrame = %d, want 0", m.tv.tstate)
	}
	// flash toggles every 16 frames
	cases := map[int]bool{1: true, 15: true, 16: false, 31: false, 32: true}
	for n := 2; n <= 32; n++ {
		m.tv.EndFrame()
		if want, ok := cases[n]; ok && m.tv.flash != want {
			t.Errorf("frame %d flash = %v, want %v", n, m.tv.flash, want)
		}
	}
}

func TestUlaTvResync(t *testing.T) {
	cases := []struct {
		name   string
		clock  int
		tstate int
	}{
		{"aligns down to block", 14343, 14340},
		{"clamps to frame end", tvFrameTstates + 20, tvFrameTstates},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.at(c.clock).tv.Resync()
			if m.tv.tstate != c.tstate {
				t.Errorf("tstate = %d, want %d", m.tv.tstate, c.tstate)
			}
			if p := m.pixel(0, 0); p != colBlack {
				t.Errorf("Resync rendered pixel(0,0) = %#x, want untouched %#x", p, colBlack)
			}
		})
	}
}

func TestSpectrumResetRestartsTv(t *testing.T) {
	m := newTestSpectrum(t)
	m.at(30000).tv.Update(30000)
	m.Reset()
	if m.clock.Tstates() != 0 || m.tv.tstate != 0 {
		t.Errorf("after Reset clock/tv = %d/%d, want 0/0", m.clock.Tstates(), m.tv.tstate)
	}
}

func TestSpectrumSnapshotResyncsTv(t *testing.T) {
	m := newTestSpectrum(t)
	m.at(30000).tv.Update(30000)
	sna := make([]byte, 49179)
	sna[26] = 2 // border
	m.LoadState(machine.State{Format: format.SNA, Data: sna})
	if m.tv.tstate != m.clock.Tstates() {
		t.Errorf("tv tstate = %d, want clock %d", m.tv.tstate, m.clock.Tstates())
	}
	if m.tv.border != 2 {
		t.Errorf("border = %d, want 2", m.tv.border)
	}
}
