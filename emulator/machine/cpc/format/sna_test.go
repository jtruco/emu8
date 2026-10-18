package format

import "testing"

func newTestCpcSnapshot() *Snapshot {
	snap := NewSnapshot()
	snap.A, snap.PC, snap.SP, snap.IM = 0x42, 0x4000, 0xbfe0, 1
	snap.GaSelectedPen = 3
	snap.GaPenColours[16] = 0x14
	snap.CrtcRegisters[1] = 40
	snap.PsgRegisters[7] = 0x3f
	snap.Memory[0x0000] = 0x11
	snap.Memory[0xffff] = 0x22
	return snap
}

func TestCpcSNASaveLoadRoundtrip(t *testing.T) {
	snap := LoadSNA(newTestCpcSnapshot().SaveSNA())
	if snap == nil {
		t.Fatal("LoadSNA failed")
	}
	if snap.A != 0x42 || snap.PC != 0x4000 || snap.SP != 0xbfe0 || snap.IM != 1 {
		t.Errorf("Z80 state not restored: A %#02x PC %#04x SP %#04x IM %d", snap.A, snap.PC, snap.SP, snap.IM)
	}
	if snap.GaSelectedPen != 3 || snap.GaPenColours[16] != 0x14 || snap.CrtcRegisters[1] != 40 || snap.PsgRegisters[7] != 0x3f {
		t.Error("devices state not restored")
	}
	if snap.Memory[0x0000] != 0x11 || snap.Memory[0xffff] != 0x22 {
		t.Error("memory not restored")
	}
}

func TestCpcSNALoadValidation(t *testing.T) {
	valid := newTestCpcSnapshot().SaveSNA()
	withSize := func(kb byte) []byte {
		d := append([]byte(nil), valid...)
		d[0x6b] = kb
		return d
	}
	cases := []struct {
		name string
		data []byte
	}{
		{"empty", nil},
		{"short id", []byte("MV - S")},
		{"truncated header", valid[:0x80]},
		{"wrong id", append([]byte("XX - SNA"), valid[8:]...)},
		{"memory size over 64k", withSize(128)},
		{"truncated memory dump", valid[:len(valid)-1]},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if snap := LoadSNA(c.data); snap != nil {
				t.Error("LoadSNA accepted invalid data")
			}
		})
	}
	if snap := LoadSNA(withSize(0)); snap == nil {
		t.Error("LoadSNA rejected a snapshot without memory dump")
	}
}

func FuzzLoadCpcSNA(f *testing.F) {
	f.Add(newTestCpcSnapshot().SaveSNA())
	f.Add(newTestCpcSnapshot().SaveSNA()[:0x100])
	f.Fuzz(func(t *testing.T, data []byte) {
		LoadSNA(data)
	})
}
