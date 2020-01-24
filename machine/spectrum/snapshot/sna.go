package snapshot

import "log"

// -----------------------------------------------------------------------------
// SNA format
// -----------------------------------------------------------------------------

const (
	_SNAFileLength = 49179
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
	snap.SP = uint16(data[23]) | (uint16(data[24]) << 8)
	snap.IM = data[25] & 0x03
	snap.PC = 0x72 // RETN at address 0x72
	snap.Border = data[26] & 0x07
	for i := 0; i < 0xc000; i++ {
		snap.Memory[i] = data[i+27]
	}
	snap.Tstates = 0
	return snap
}
