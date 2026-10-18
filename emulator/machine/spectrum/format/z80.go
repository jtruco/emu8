package format

import (
	"bytes"
	"log"
)

// -----------------------------------------------------------------------------
// Z80 snapshot format
// Versions : 1..3. Only 48k mode suppported
// -----------------------------------------------------------------------------

// Z80 format extension
const Z80 = "z80"

const (
	_Z80HeaderLength     = 30
	_Z80ExtraMinLength   = 23 // v2 additional header length (v3 : 54 or 55)
	_Z80BankSize         = 0x4000
	_Z80BankUncompressed = 0xffff // v3 bank length : 16k not compressed
	_Z80V1MemorySize     = 0xc000 // v1 : 48k memory from 0x4000
)

var z80V1EndMarker = []byte{0x00, 0xed, 0xed, 0x00}

var z80BankAddresses = map[byte]int{8: 0x0000, 4: 0x4000, 5: 0x8000}

// LoadZ80 loads snap from Z80 data format
func LoadZ80(data []byte) *Snapshot {
	// Check header (versions 1, 2 and 3)
	if len(data) < _Z80HeaderLength {
		log.Println("Z80 : Invalid file format")
		return nil
	}
	snap := NewSnapshot()
	if !z80LoadHeaderV1(data, snap) {
		return nil
	}
	PC := uint16(data[6]) | (uint16(data[7]) << 8)
	if PC > 0 { // v1 format
		if !z80LoadFileV1(data, snap) {
			return nil
		}
	} else { // v2/v3 format
		if !z80LoadFileV23(data, snap) {
			return nil
		}
	}
	return snap
}

func z80LoadHeaderV1(data []byte, snap *Snapshot) bool {
	snap.A = data[0]
	snap.F = data[1]
	snap.C = data[2]
	snap.B = data[3]
	snap.L = data[4]
	snap.H = data[5]
	snap.PC = readWord(data, 6)
	snap.SP = readWord(data, 8)
	snap.I = data[10]
	data12 := data[12]
	if data12 == 255 {
		data12 = 1
	}
	snap.R = (data[11] & 0x7f) | ((data12 & 0x01) << 7)
	snap.Border = (data12 >> 1) & 0x07
	snap.E = data[13]
	snap.D = data[14]
	snap.Cx = data[15]
	snap.Bx = data[16]
	snap.Ex = data[17]
	snap.Dx = data[18]
	snap.Lx = data[19]
	snap.Hx = data[20]
	snap.Ax = data[21]
	snap.Fx = data[22]
	snap.IYl = data[23]
	snap.IYh = data[24]
	snap.IXl = data[25]
	snap.IXh = data[26]
	snap.IFF1 = (data[27] != 0)
	snap.IFF2 = (data[28] != 0)
	snap.IM = data[29] & 0x03
	snap.Tstates = 0
	return true
}

// z80LoadFileV1 loads the 48k memory of a version 1 snapshot. It is
// compressed if header byte 12 bit 5 is set, ending with a 00 ED ED 00 marker.
func z80LoadFileV1(data []byte, snap *Snapshot) bool {
	memory := data[_Z80HeaderLength:]
	compressed := data[12] != 255 && data[12]&0x20 != 0
	if compressed {
		if n := len(memory); n >= len(z80V1EndMarker) && bytes.Equal(memory[n-len(z80V1EndMarker):], z80V1EndMarker) {
			memory = memory[:n-len(z80V1EndMarker)]
		}
		var ok bool
		memory, ok = z80Decompress(memory, _Z80V1MemorySize)
		if !ok {
			log.Println("Z80 : wrong compressed memory data")
			return false
		}
	}
	if len(memory) != _Z80V1MemorySize {
		log.Println("Z80 : wrong memory size")
		return false
	}
	copy(snap.Memory[:], memory)
	return true
}

func z80LoadFileV23(data []byte, snap *Snapshot) bool {
	totalSize := len(data)
	if totalSize < _Z80HeaderLength+2 {
		log.Println("Z80 : Invalid file format")
		return false
	}
	extraSize := readInt(data, _Z80HeaderLength)
	headerSize := _Z80HeaderLength + 2 + extraSize
	if extraSize < _Z80ExtraMinLength || totalSize < headerSize {
		log.Println("Z80 : Invalid file format")
		return false
	}
	if data[34] != 0 { // hadware mode : only 48k mode supported
		log.Println("Z80 : Unsupported machine hardware mode")
	}
	if data[37]>>7 == 1 { // machine modification not supported
		log.Println("Z80 : Unsupported hardware modification")
	}
	// v2/v3 Z80 program counter
	snap.PC = uint16(data[32]) | (uint16(data[33]) << 8)
	// load 16k data banks. 48k model RAM banks
	for idx := headerSize; idx < totalSize; {
		if idx+3 > totalSize {
			log.Println("Z80 : truncated bank header")
			return false
		}
		num := data[idx+2]
		address, ok := z80BankAddresses[num]
		if !ok {
			log.Println("Z80 : unsupported bank format")
			return false
		}
		size := readInt(data, idx)
		compressed := true
		switch {
		case size == _Z80BankUncompressed: // v3 : 16k not compressed
			size, compressed = _Z80BankSize, false
		case size == _Z80BankSize:
			compressed = false
		case size > _Z80BankSize:
			log.Println("Z80 : wrong bank size")
			return false
		}
		idx += 3
		if idx+size > totalSize {
			log.Println("Z80 : truncated bank data")
			return false
		}
		bankdata := data[idx : idx+size]
		if compressed {
			bankdata, ok = z80DecompressBlock(bankdata)
			if !ok {
				log.Println("Z80 : wrong compressed bank data")
				return false
			}
		}
		copy(snap.Memory[address:address+_Z80BankSize], bankdata)
		idx += size
	}
	return true
}

// z80DecompressBlock decompresses a 16k bank. It fails if the data
// decompresses to more than 16k.
func z80DecompressBlock(data []byte) ([]byte, bool) {
	return z80Decompress(data, _Z80BankSize)
}

// z80Decompress decompresses data (ED ED nn bb : nn times bb). It fails if
// the data decompresses to more than size bytes.
func z80Decompress(data []byte, size int) ([]byte, bool) {
	buffer := make([]byte, size)
	sizeIn := len(data)
	sizeOut := 0
	for i := 0; i < sizeIn; {
		if data[i] == 0xED && i < (sizeIn-3) && data[i+1] == 0xED {
			count := int(data[i+2])
			if sizeOut+count > size {
				return nil, false
			}
			for j := 0; j < count; j++ {
				buffer[sizeOut] = data[i+3]
				sizeOut++
			}
			i += 4
		} else {
			if sizeOut == size {
				return nil, false
			}
			buffer[sizeOut] = data[i]
			sizeOut++
			i++
		}
	}
	return buffer[0:sizeOut], true
}
