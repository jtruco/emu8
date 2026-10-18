package format

import (
	"bytes"
	"testing"
)

// z80Image builds a Z80 v2 snapshot (48k) with PC and raw banks
func z80Image(pc uint16, banks ...[]byte) []byte {
	data := make([]byte, _Z80HeaderLength+2+_Z80ExtraMinLength)
	data[_Z80HeaderLength] = _Z80ExtraMinLength
	data[32], data[33] = byte(pc), byte(pc>>8)
	for _, b := range banks {
		data = append(data, b...)
	}
	return data
}

// z80Bank builds a bank block : length, page and data
func z80Bank(page byte, data []byte) []byte {
	return append([]byte{byte(len(data)), byte(len(data) >> 8), page}, data...)
}

func TestZ80LoadBanks(t *testing.T) {
	raw := bytes.Repeat([]byte{0x5a}, _Z80BankSize)
	compressed := []byte{0x01, 0xed, 0xed, 0x03, 0xaa, 0x02} // 1, 3 x 0xaa, 2
	snap := LoadZ80(z80Image(0x8000, z80Bank(4, raw), z80Bank(8, compressed)))
	if snap == nil {
		t.Fatal("LoadZ80 failed")
	}
	if snap.PC != 0x8000 {
		t.Errorf("PC = %#04x, want 0x8000", snap.PC)
	}
	if snap.Memory[0x4000] != 0x5a || snap.Memory[0x7fff] != 0x5a {
		t.Error("uncompressed bank 4 (0x8000) not loaded")
	}
	if got := snap.Memory[0:5]; !bytes.Equal(got, []byte{0x01, 0xaa, 0xaa, 0xaa, 0x02}) {
		t.Errorf("compressed bank 8 (0x4000) = % x", got)
	}
}

func TestZ80LoadV3UncompressedBank(t *testing.T) {
	raw := bytes.Repeat([]byte{0xed, 0xed, 0x03, 0xaa}, _Z80BankSize/4) // not a run
	bank := append([]byte{0xff, 0xff, 5}, raw...)
	snap := LoadZ80(z80Image(0x8000, bank))
	if snap == nil {
		t.Fatal("LoadZ80 rejected a v3 uncompressed bank (length 0xffff)")
	}
	if got := snap.Memory[0x8000:0xc000]; !bytes.Equal(got, raw) {
		t.Error("uncompressed bank 5 (0xc000) was not loaded as is")
	}
}

func TestZ80LoadValidation(t *testing.T) {
	overflow := []byte{0xed, 0xed, 0xff, 0x00}
	overflow = bytes.Repeat(overflow, _Z80BankSize/0xff+1)
	cases := []struct {
		name string
		data []byte
	}{
		{"empty", nil},
		{"short header", make([]byte, _Z80HeaderLength-1)},
		{"v2 without extra length", make([]byte, _Z80HeaderLength)},
		{"extra header too short", func() []byte {
			d := make([]byte, _Z80HeaderLength+2+10)
			d[_Z80HeaderLength] = 10
			return d
		}()},
		{"truncated extra header", z80Image(0)[:_Z80HeaderLength+10]},
		{"truncated bank header", append(z80Image(0), 0x10, 0x00)},
		{"truncated bank data", append(z80Image(0), 0x10, 0x00, 8, 0x01)},
		{"unknown bank", z80Image(0, z80Bank(3, []byte{0x01}))},
		{"bank too big", append(z80Image(0), 0x01, 0x40, 8)},
		{"truncated uncompressed bank", append(z80Image(0), append([]byte{0xff, 0xff, 8}, make([]byte, 100)...)...)},
		{"decompresses to more than 16k", z80Image(0, z80Bank(8, overflow))},
		{"literal bytes after 16k", z80Image(0, z80Bank(8, append(bytes.Repeat([]byte{0xed, 0xed, 0x80, 0x00}, 128), 0x01)))},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if snap := LoadZ80(c.data); snap != nil {
				t.Error("LoadZ80 accepted invalid data")
			}
		})
	}
}

func FuzzLoadZ80(f *testing.F) {
	f.Add(z80V1Image(0x8000, 0x20, z80Compress(z80TestMemory()), z80V1EndMarker))
	f.Add(z80V1Image(0x8000, 0x00, z80TestMemory()))
	f.Add(z80Image(0x8000, z80Bank(8, []byte{0x01, 0xed, 0xed, 0x03, 0xaa, 0x02})))
	f.Add(z80Image(0x8000, z80Bank(4, make([]byte, _Z80BankSize))))
	f.Add(z80Image(0x8000, append([]byte{0xff, 0xff, 5}, make([]byte, _Z80BankSize)...)))
	f.Add(make([]byte, _Z80HeaderLength))
	f.Fuzz(func(t *testing.T, data []byte) {
		LoadZ80(data)
	})
}

func FuzzLoadSNA(f *testing.F) {
	f.Add(newTestSnapshot().SaveSNA())
	f.Fuzz(func(t *testing.T, data []byte) {
		LoadSNA(data)
	})
}

// z80V1Image builds a Z80 v1 snapshot : header (PC, SP, byte 12) and data
func z80V1Image(pc uint16, byte12 byte, parts ...[]byte) []byte {
	data := make([]byte, _Z80HeaderLength)
	data[6], data[7] = byte(pc), byte(pc>>8)
	data[8], data[9] = 0x00, 0xff // SP 0xff00
	data[12] = byte12
	for _, p := range parts {
		data = append(data, p...)
	}
	return data
}

// z80Compress compresses data as in the Z80 format: runs of 5 or more equal
// bytes, and runs of 2 or more ED bytes, as ED ED nn bb. The byte after a
// single ED is never compressed.
func z80Compress(data []byte) []byte {
	var out []byte
	for i := 0; i < len(data); {
		b, run := data[i], 1
		for i+run < len(data) && data[i+run] == b && run < 255 {
			run++
		}
		switch {
		case run >= 5 || (b == 0xed && run >= 2):
			out = append(out, 0xed, 0xed, byte(run), b)
			i += run
		case b == 0xed && i+1 < len(data):
			out = append(out, b, data[i+1])
			i += 2
		default:
			out = append(out, b)
			i++
		}
	}
	return out
}

// z80TestMemory is 48k of data with runs, single and repeated ED bytes
func z80TestMemory() []byte {
	mem := make([]byte, _Z80V1MemorySize)
	for i := range mem {
		switch {
		case i%997 < 300: // long runs
			mem[i] = byte(i / 997)
		case i%13 == 0:
			mem[i] = 0xed
		default:
			mem[i] = byte(i * 7)
		}
	}
	copy(mem[0x100:], []byte{0xed, 0xed, 0xed, 0x05, 0xed, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05})
	return mem
}

func TestZ80LoadV1(t *testing.T) {
	mem := z80TestMemory()
	compressed := z80Compress(mem)
	cases := []struct {
		name string
		data []byte
	}{
		{"uncompressed", z80V1Image(0x8123, 0x00, mem)},
		{"byte 12 is 255 (uncompressed)", z80V1Image(0x8123, 0xff, mem)},
		{"compressed with end marker", z80V1Image(0x8123, 0x20, compressed, z80V1EndMarker)},
		{"compressed without end marker", z80V1Image(0x8123, 0x20, compressed)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			snap := LoadZ80(c.data)
			if snap == nil {
				t.Fatal("LoadZ80 failed")
			}
			if snap.PC != 0x8123 || snap.SP != 0xff00 {
				t.Errorf("PC/SP = %#04x/%#04x, want 0x8123/0xff00", snap.PC, snap.SP)
			}
			if !bytes.Equal(snap.Memory[:], mem) {
				t.Error("memory does not match")
			}
		})
	}
	if len(compressed) >= len(mem) {
		t.Errorf("test data compressed to %d bytes, want less than %d", len(compressed), len(mem))
	}
}

func TestZ80LoadV1Validation(t *testing.T) {
	mem := z80TestMemory()
	cases := []struct {
		name string
		data []byte
	}{
		{"uncompressed too short", z80V1Image(0x8000, 0x00, mem[:len(mem)-1])},
		{"uncompressed too long", z80V1Image(0x8000, 0x00, mem, []byte{0x00})},
		{"compressed too short", z80V1Image(0x8000, 0x20, z80Compress(mem[:0x4000]), z80V1EndMarker)},
		{"compressed too long", z80V1Image(0x8000, 0x20, z80Compress(append(mem, mem[:16]...)), z80V1EndMarker)},
		{"no memory", z80V1Image(0x8000, 0x20)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if snap := LoadZ80(c.data); snap != nil {
				t.Error("LoadZ80 accepted invalid data")
			}
		})
	}
}
