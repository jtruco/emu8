package format

import "log"

// -----------------------------------------------------------------------------
// SNA format
// -----------------------------------------------------------------------------

// SNA format extension
const SNA = "sna"

const (
	_SNAFileLength = 49179
	_SNAHeaderSize = 27
	_SNARamAddress = 0x4000 // SNA memory dump starts at RAM address 0x4000
)

// LoadSNA loads snap from SNA data format
func LoadSNA(data []byte) *Snapshot {
	// Check format
	if len(data) != _SNAFileLength {
		log.Println("SNA : Invalid file format")
		return nil
	}

	// load SNA data
	snap := NewSnapshot()
	snap.I = data[0]
	snap.Lx = data[1]
	snap.Hx = data[2]
	snap.Ex = data[3]
	snap.Dx = data[4]
	snap.Cx = data[5]
	snap.Bx = data[6]
	snap.Fx = data[7]
	snap.Ax = data[8]
	snap.L = data[9]
	snap.H = data[10]
	snap.E = data[11]
	snap.D = data[12]
	snap.C = data[13]
	snap.B = data[14]
	snap.IYl = data[15]
	snap.IYh = data[16]
	snap.IXl = data[17]
	snap.IXh = data[18]
	intEnabled := (data[19] & 0x04) != 0
	snap.IFF1 = intEnabled
	snap.IFF2 = intEnabled
	snap.R = data[20]
	snap.F = data[21]
	snap.A = data[22]
	snap.SP = readWord(data, 23)
	snap.IM = data[25] & 0x03
	snap.PC = 0x72 // RETN at address 0x72
	snap.Border = data[26] & 0x07
	copy(snap.Memory[0:0xc000], data[27:])
	snap.Tstates = 0
	return snap
}

// SaveSNA saves snap to SNA data format
func (snap *Snapshot) SaveSNA() []byte {
	var data = make([]byte, _SNAFileLength)

	// save SNA data
	data[0] = snap.I
	data[1] = snap.Lx
	data[2] = snap.Hx
	data[3] = snap.Ex
	data[4] = snap.Dx
	data[5] = snap.Cx
	data[6] = snap.Bx
	data[7] = snap.Fx
	data[8] = snap.Ax
	data[9] = snap.L
	data[10] = snap.H
	data[11] = snap.E
	data[12] = snap.D
	data[13] = snap.C
	data[14] = snap.B
	data[15] = snap.IYl
	data[16] = snap.IYh
	data[17] = snap.IXl
	data[18] = snap.IXh
	if snap.IFF1 {
		data[19] |= 0x04
	}
	data[20] = snap.R
	data[21] = snap.F
	data[22] = snap.A
	data[25] = snap.IM & 0x03
	data[26] = snap.Border & 0x07
	copy(data[_SNAHeaderSize:], snap.Memory[0:0xc000])
	// SNA stores PC on the stack (LoadSNA resumes with RETN): push PC
	sp := snap.SP - 2
	writeWord(data, 23, sp)
	snaWriteRAM(data, sp, byte(snap.PC))
	snaWriteRAM(data, sp+1, byte(snap.PC>>8))
	return data
}

// snaWriteRAM writes a byte of the SNA memory dump at a machine address
func snaWriteRAM(data []byte, address uint16, value byte) {
	if address < _SNARamAddress {
		log.Println("SNA : Stack in ROM, PC can not be saved")
		return
	}
	data[_SNAHeaderSize+int(address-_SNARamAddress)] = value
}
