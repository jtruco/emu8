package format

import "log"

// -----------------------------------------------------------------------------
// Z80 snapshot format
// Versions : 1..3. Only 48k mode suppported
// -----------------------------------------------------------------------------

const (
	_Z80HeaderLength = 30
	_Z80BankSize     = 0x4000
)

var (
	_Z80BankMap = map[byte]int{8: 0x0000, 4: 0x4000, 5: 0x8000}
)

// LoadZ80 loads snap from Z80 data format
func LoadZ80(data []byte) *Snapshot {
	// Check version (only v1 supported)
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

func z80LoadFileV1(data []byte, snap *Snapshot) bool {
	// not implemented
	log.Println("Z80 : Version 1 not implemented")
	return false
}

func z80LoadFileV23(data []byte, snap *Snapshot) bool {
	totalSize := len(data)
	extraSize := uint16(data[30]) | (uint16(data[31]) << 8)
	headerSize := int(_Z80HeaderLength + 2 + extraSize)
	if totalSize < headerSize {
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
		num := data[idx+2]
		address, ok := _Z80BankMap[num]
		if !ok {
			log.Println("Z80 : unsupported bank format")
			return false
		}
		size := int(uint16(data[idx+0]) | uint16(data[idx+1])<<8)
		if size > _Z80BankSize {
			log.Println("Z80 : wrong bank size")
			return false
		}
		idx += 3
		bankdata := data[idx : idx+size]
		if size < _Z80BankSize {
			bankdata = z80DecompressBlock(bankdata)
		}
		copy(snap.Memory[address:], bankdata[:_Z80BankSize])
		idx += size
	}
	return true
}

func z80DecompressBlock(data []byte) []byte {
	buffer := make([]byte, _Z80BankSize)
	sizeIn := len(data)
	sizeOut := 0
	for i := 0; i < sizeIn; {
		if data[i] == 0xED && i < (sizeIn-3) && data[i+1] == 0xED {
			for j := byte(0); j < data[i+2]; j++ {
				buffer[sizeOut] = data[i+3]
				sizeOut++
			}
			i += 4
		} else {
			buffer[sizeOut] = data[i]
			sizeOut++
			i++
		}
	}
	return buffer[0:sizeOut]
}
