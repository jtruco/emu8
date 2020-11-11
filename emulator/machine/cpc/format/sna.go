package format

import "log"

// -----------------------------------------------------------------------------
// CPC SNA format
// -----------------------------------------------------------------------------

const (
	_SNAIdString = "MV - SNA"
	_SNAMemDump  = 0x100
)

// LoadSNA loads snap from SNA data format
func LoadSNA(data []byte) *Snapshot {
	// Check format
	idstr := string(data[:0x08])
	if idstr != _SNAIdString {
		log.Println("SNA : Invalid file format")
		return nil
	}
	// Check version
	version := data[0x10]
	if version != 1 {
		log.Println("SNA : Unsupported version: ", version)
	}
	// Check men size
	memsize := int(readWord(data, 0x6b))
	if memsize > 64 {
		log.Println("SNA : Only 64K snapshopts supported")
		return nil
	}
	memsize *= 0x400 // 1k

	// load SNA data
	snap := NewSnapshot()
	// Z80 state
	snap.F = data[0x11]
	snap.A = data[0x12]
	snap.C = data[0x13]
	snap.B = data[0x14]
	snap.E = data[0x15]
	snap.D = data[0x16]
	snap.L = data[0x17]
	snap.H = data[0x18]
	snap.R = data[0x19]
	snap.I = data[0x1a]
	snap.IFF1 = data[0x1b] != 0
	snap.IFF2 = data[0x1c] != 0
	snap.IXl = data[0x1d]
	snap.IXh = data[0x1e]
	snap.IYl = data[0x1f]
	snap.IYh = data[0x20]
	snap.SP = readWord(data, 0x21)
	snap.PC = readWord(data, 0x23)
	snap.IM = data[0x25] & 0x03
	snap.Fx = data[0x26]
	snap.Ax = data[0x27]
	snap.Cx = data[0x28]
	snap.Bx = data[0x29]
	snap.Ex = data[0x2a]
	snap.Dx = data[0x2b]
	snap.Lx = data[0x2c]
	snap.Hx = data[0x2d]

	// Gatearray
	snap.GaSelectedPen = data[0x2e]
	copy(snap.GaPenColours[:], data[0x2f:])
	snap.GaMultiConfig = data[0x40]
	snap.GaRAMSelect = data[0x41]

	// Crtc
	snap.CrtcSelected = data[0x42]
	copy(snap.CrtcRegisters[:], data[0x43:])

	// ROM select
	snap.RomSelect = data[0x55]

	// PPI
	snap.PpiPortA = data[0x56]
	snap.PpiPortB = data[0x57]
	snap.PpiPortC = data[0x58]
	snap.PpiControl = data[0x59]

	// PSG
	snap.PsgSelected = data[0x5a]
	copy(snap.PsgRegisters[:], data[0x5b:])

	// Memory dump
	copy(snap.Memory[0:memsize], data[_SNAMemDump:])

	return snap
}

// SaveSNA saves CPC464 snapshot to SNA data format
func (snap *Snapshot) SaveSNA() []byte {
	const memsize = 64 * 0x400 // 64k
	var data = make([]byte, _SNAMemDump+memsize)

	// Format string, version and size
	copy(data[:8], _SNAIdString)
	data[0x10] = 1  // version
	data[0x6b] = 64 // size

	// Z80 state
	data[0x11] = snap.F
	data[0x12] = snap.A
	data[0x13] = snap.C
	data[0x14] = snap.B
	data[0x15] = snap.E
	data[0x16] = snap.D
	data[0x17] = snap.L
	data[0x18] = snap.H
	data[0x19] = snap.R
	data[0x1a] = snap.I
	if snap.IFF1 {
		data[0x1b] |= 1
	}
	if snap.IFF2 {
		data[0x1c] |= 1
	}
	data[0x1d] = snap.IXl
	data[0x1e] = snap.IXh
	data[0x1f] = snap.IYl
	data[0x20] = snap.IYh
	writeWord(data, 0x21, snap.SP)
	writeWord(data, 0x23, snap.PC)
	data[0x25] = snap.IM & 0x03
	data[0x26] = snap.Fx
	data[0x27] = snap.Ax
	data[0x28] = snap.Cx
	data[0x29] = snap.Bx
	data[0x2a] = snap.Ex
	data[0x2b] = snap.Dx
	data[0x2c] = snap.Lx
	data[0x2d] = snap.Hx

	// Gatearray
	data[0x2e] = snap.GaSelectedPen
	copy(data[0x2f:], snap.GaPenColours[:])
	data[0x40] = snap.GaMultiConfig
	data[0x41] = snap.GaRAMSelect

	// Crtc
	data[0x42] = snap.CrtcSelected
	copy(data[0x43:], snap.CrtcRegisters[:])

	// ROM select
	data[0x55] = snap.RomSelect

	// PPI
	data[0x56] = snap.PpiPortA
	data[0x57] = snap.PpiPortB
	data[0x58] = snap.PpiPortC
	data[0x59] = snap.PpiControl

	// PSG
	data[0x5a] = snap.PsgSelected
	copy(data[0x5b:], snap.PsgRegisters[:])

	// Memory dump
	copy(data[_SNAMemDump:], snap.Memory[0:memsize])

	return data
}
