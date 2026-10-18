package format

import "testing"

func newTestSnapshot() *Snapshot {
	snap := NewSnapshot()
	snap.A, snap.F, snap.B, snap.C = 0x12, 0x34, 0x56, 0x78
	snap.I, snap.R, snap.IM = 0x3f, 0x5a, 1
	snap.IXh, snap.IXl, snap.IYh, snap.IYl = 0x9a, 0xbc, 0xde, 0xf0
	snap.IFF1, snap.IFF2 = true, true
	snap.Border = 5
	snap.PC = 0x8123
	snap.SP = 0xff00
	return snap
}

// stackWord reads the 16 bit word at a machine address of the snapshot RAM
func stackWord(snap *Snapshot, address uint16) uint16 {
	i := int(address) - _SNARamAddress
	return uint16(snap.Memory[i]) | uint16(snap.Memory[i+1])<<8
}

func TestSNASaveLoadRoundtrip(t *testing.T) {
	snap := newTestSnapshot()
	data := snap.SaveSNA()
	if len(data) != _SNAFileLength {
		t.Fatalf("SNA length = %d, want %d", len(data), _SNAFileLength)
	}
	loaded := LoadSNA(data)
	if loaded == nil {
		t.Fatal("LoadSNA failed")
	}

	// PC is pushed on the stack and restored by RETN at 0x72
	if loaded.PC != 0x72 {
		t.Errorf("loaded PC = %#04x, want RETN address 0x0072", loaded.PC)
	}
	if loaded.SP != 0xfefe {
		t.Errorf("loaded SP = %#04x, want 0xfefe", loaded.SP)
	}
	if pc := stackWord(loaded, loaded.SP); pc != 0x8123 {
		t.Errorf("PC on stack = %#04x, want 0x8123", pc)
	}

	// registers
	if loaded.A != 0x12 || loaded.F != 0x34 || loaded.B != 0x56 || loaded.C != 0x78 ||
		loaded.I != 0x3f || loaded.R != 0x5a || loaded.IM != 1 ||
		loaded.IXh != 0x9a || loaded.IXl != 0xbc || loaded.IYh != 0xde || loaded.IYl != 0xf0 ||
		!loaded.IFF1 || loaded.Border != 5 {
		t.Errorf("registers not restored: %+v", loaded.State)
	}
}

func TestSNASaveKeepsSnapshot(t *testing.T) {
	snap := newTestSnapshot()
	snap.SaveSNA()
	if snap.SP != 0xff00 {
		t.Errorf("snapshot SP modified to %#04x", snap.SP)
	}
	if w := stackWord(snap, 0xfefe); w != 0 {
		t.Errorf("snapshot memory modified at stack: %#04x", w)
	}
}

func TestSNASaveStackEdges(t *testing.T) {
	cases := []struct {
		name   string
		sp     uint16
		wantSP uint16
		low    int // expected PC low byte memory index (-1 : not written)
		high   int // expected PC high byte memory index (-1 : not written)
	}{
		{"stack at top of RAM wraps", 0x0000, 0xfffe, 0xbffe, 0xbfff},
		{"first RAM bytes", 0x4002, 0x4000, 0x0000, 0x0001},
		{"high byte in RAM, low byte in ROM", 0x4001, 0x3fff, -1, 0x0000},
		{"stack in ROM", 0x2000, 0x1ffe, -1, -1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			snap := newTestSnapshot()
			snap.SP = c.sp
			loaded := LoadSNA(snap.SaveSNA())
			if loaded.SP != c.wantSP {
				t.Errorf("SP = %#04x, want %#04x", loaded.SP, c.wantSP)
			}
			if c.low >= 0 && loaded.Memory[c.low] != 0x23 {
				t.Errorf("PC low byte = %#02x, want 0x23", loaded.Memory[c.low])
			}
			if c.high >= 0 && loaded.Memory[c.high] != 0x81 {
				t.Errorf("PC high byte = %#02x, want 0x81", loaded.Memory[c.high])
			}
		})
	}
}

func TestLoadSNAInvalidLength(t *testing.T) {
	for _, n := range []int{0, 27, _SNAFileLength - 1, _SNAFileLength + 1} {
		if LoadSNA(make([]byte, n)) != nil {
			t.Errorf("LoadSNA accepted %d bytes", n)
		}
	}
}

// SNA byte 19 bit 2 stores IFF2 (the loader sets IFF1 = IFF2)
func TestSNASaveIFF2(t *testing.T) {
	cases := []struct{ iff1, iff2, want bool }{
		{false, true, true}, {true, false, false}, {true, true, true}, {false, false, false},
	}
	for _, c := range cases {
		snap := newTestSnapshot()
		snap.IFF1, snap.IFF2 = c.iff1, c.iff2
		if got := snap.SaveSNA()[19]&0x04 != 0; got != c.want {
			t.Errorf("IFF1 %v IFF2 %v: interrupt bit = %v, want %v", c.iff1, c.iff2, got, c.want)
		}
	}
}
