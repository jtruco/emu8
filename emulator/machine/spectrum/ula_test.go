package spectrum

import "testing"

func TestUlaContentionTable(t *testing.T) {
	// 48K pattern 6,5,4,3,2,1,0,0 from Tstate 14335, 16 groups per paper
	// line, every 224 Ts, for 192 lines.
	cases := []struct {
		tstate, delay int
	}{
		{0, 0},
		{14334, 0},
		{14335, 6}, {14336, 5}, {14337, 4}, {14338, 3},
		{14339, 2}, {14340, 1}, {14341, 0}, {14342, 0},
		{14343, 6},
		{14335 + 120, 6},     // last group of the line
		{14335 + 127, 0},     // end of the last group
		{14335 + 128, 0},     // right border
		{14335 + 224, 6},     // next line
		{14335 + 191*224, 6}, // last paper line
		{14335 + 192*224, 0}, // bottom border
		{zxTStates - 1, 0},   // end of frame
	}
	for _, c := range cases {
		if got := ulaDelayTable[c.tstate]; got != c.delay {
			t.Errorf("delay at %d = %d, want %d", c.tstate, got, c.delay)
		}
	}
}

func TestUlaMemoryContention(t *testing.T) {
	cases := []struct {
		name    string
		address uint16
		tstate  int
		delay   int
	}{
		{"video bank contended", 0x4000, 14335, 6},
		{"video bank pattern", 0x5800, 14338, 3},
		{"video bank idle slot", 0x4000, 14341, 0},
		{"video bank outside paper", 0x4000, 1000, 0},
		{"upper RAM not contended", 0x8000, 14335, 0},
		{"ROM not contended", 0x0000, 14335, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.at(c.tstate).memory.Read(c.address)
			if got := m.clock.Tstates() - c.tstate; got != c.delay {
				t.Errorf("read delay = %d, want %d", got, c.delay)
			}
			m.at(c.tstate).memory.Write(c.address, 0)
			if got := m.clock.Tstates() - c.tstate; got != c.delay {
				t.Errorf("write delay = %d, want %d", got, c.delay)
			}
		})
	}
}

func TestUlaIOContention(t *testing.T) {
	// 48K I/O patterns (high byte contended / ULA port bit 0):
	// N:1 C:3, N:4, C:1 C:3, C:1 C:1 C:1 C:1
	cases := []struct {
		name   string
		port   uint16
		tstate int
		end    int
	}{
		{"ULA port, idle", 0x00fe, 0, 4},
		{"other port, idle", 0x00ff, 0, 4},
		{"ULA port contended high, idle", 0x40fe, 0, 4},
		{"other port contended high, idle", 0x40ff, 0, 4},
		{"ULA port N:1 C:3", 0x00fe, 14335, 14344},
		{"other port N:4", 0x00ff, 14335, 14339},
		{"ULA port contended high C:1 C:3", 0x40fe, 14335, 14345},
		{"other port contended high C:1 x4", 0x40ff, 14335, 14351},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.at(c.tstate).ula.Read(c.port)
			if got := m.clock.Tstates(); got != c.end {
				t.Errorf("read ends at %d, want %d", got, c.end)
			}
			m.at(c.tstate).ula.Write(c.port, 7)
			if got := m.clock.Tstates(); got != c.end {
				t.Errorf("write ends at %d, want %d", got, c.end)
			}
		})
	}
}

func TestUlaWriteBorder(t *testing.T) {
	m := newTestSpectrum(t)
	m.at(14335).ula.Write(0x00fe, 0x02)
	if m.tv.border != 2 {
		t.Errorf("border = %d, want 2", m.tv.border)
	}
	// the change happens after the 1 Ts of the I/O cycle (14336), so the
	// beam is rendered up to that Tstate before the change
	if m.tv.tstate != 14336 {
		t.Errorf("tv tstate = %d, want 14336", m.tv.tstate)
	}
	m.tv.EndFrame()
	if p := m.pixel(31, 24); p != colWhite {
		t.Errorf("left border line 64 = %#x, want old %#x", p, colWhite)
	}
	if p := m.pixel(288, 24); p != colRed {
		t.Errorf("right border line 64 = %#x, want new %#x", p, colRed)
	}
}

func TestUlaReadPorts(t *testing.T) {
	cases := []struct {
		name  string
		write byte // value written to port 0xfe before reading
		port  uint16
		want  byte
	}{
		{"ULA no keys, EAR/MIC low (issue 2)", 0x00, 0xfefe, 0xbf},
		{"ULA no keys, MIC high", 0x08, 0xfefe, 0xff},
		{"ULA no keys, EAR high", 0x10, 0xfefe, 0xff},
		{"kempston idle", 0x18, 0x001f, 0x00},
		{"unattached port", 0x18, 0x00ff, 0xff},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.ula.Write(0x00fe, c.write)
			if got := m.ula.Read(c.port); got != c.want {
				t.Errorf("Read(%#04x) = %#02x, want %#02x", c.port, got, c.want)
			}
		})
	}
}

// A CPU write lands at the end of its 3 Ts memory cycle (ulaMemWriteTstates):
// paper blocks fetched before that keep the old byte.
func TestUlaVideoWriteTiming(t *testing.T) {
	cases := []struct {
		name   string
		tstate int
		offset int // attribute offset (line 64 cell)
		x      int
		line64 uint32
	}{
		{"completes before the block fetch", 14333, 0x1800, 32, colRed},
		{"completes at the block fetch", 14341, 0x1802, 48, colRed},
		{"completes inside the fetched block", 14341, 0x1801, 40, colBlue},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.videoPoke(c.offset, 1<<3) // blue paper
			m.at(c.tstate).memory.Write(uint16(0x4000+c.offset), 2<<3)
			m.tv.EndFrame()
			if p := m.pixel(c.x, 24); p != c.line64 {
				t.Errorf("line 64 = %#x, want %#x", p, c.line64)
			}
			if p := m.pixel(c.x, 25); p != colRed {
				t.Errorf("line 65 = %#x, want new %#x", p, colRed)
			}
		})
	}
}

func TestUlaVideoAccessRendering(t *testing.T) {
	cases := []struct {
		name    string
		write   bool
		address uint16
		tstate  int
	}{
		{"video write renders", true, 0x4000, 30000},
		{"attribute write renders", true, 0x5aff, 30000},
		{"write past video memory does not render", true, 0x4000 + tvVideoSize, 0},
		{"video read does not render", false, 0x4000, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := newTestSpectrum(t)
			m.at(30000)
			if c.write {
				m.memory.Write(c.address, 0)
			} else {
				m.memory.Read(c.address)
			}
			if m.tv.tstate != c.tstate {
				t.Errorf("tv tstate = %d, want %d", m.tv.tstate, c.tstate)
			}
		})
	}
}
