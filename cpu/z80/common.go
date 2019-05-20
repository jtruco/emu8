package z80

import (
	"github.com/jtruco/emu8/cpu"
)

// -----------------------------------------------------------------------------
// Constants and tables
// -----------------------------------------------------------------------------

// Constants
const (
	// Flags constants
	FlagC byte = 0x01
	FlagN byte = 0x02
	FlagP byte = 0x04
	FlagV byte = FlagP
	Flag3 byte = 0x08
	FlagH byte = 0x10
	Flag5 byte = 0x20
	FlagZ byte = 0x40
	FlagS byte = 0x80
)

// Precalculated tables
var (
	// The half carry lookup table for adds
	halfCarryAddTable = []byte{0, FlagH, FlagH, FlagH, 0, 0, 0, FlagH}
	// The half carry lookup table for subs
	halfCarrySubTable = []byte{0, 0, FlagH, 0, FlagH, 0, FlagH, FlagH}
	// The overflow lookup table for adds
	overflowAddTable = []byte{0, 0, 0, FlagV, FlagV, 0, 0, 0}
	// The overflow lookup table for subs
	overflowSubTable = []byte{0, FlagV, 0, 0, 0, 0, FlagV, 0}
	// The S, Z, 5 and 3 bits table
	sz53Table [0x100]byte
	// The parity table
	parityTable [0x100]byte
	// The S, Z, 5 and 3 bits and parity table
	sz53pTable [0x100]byte
)

// init initializes tables and variables
func init() {
	var i int16
	var j, k, p byte

	for i = 0; i < 0x100; i++ {
		sz53Table[i] = byte(i) & (Flag3 | Flag5 | FlagS)
		j, p = byte(i), 0
		for k = 0; k < 8; k++ {
			p ^= j & 1
			j >>= 1
		}
		parityTable[i] = ifthen(p != 0, 0, FlagP)
		sz53pTable[i] = sz53Table[i] | parityTable[i]
	}
	sz53Table[0] |= FlagZ
	sz53pTable[0] |= FlagZ
}

// -----------------------------------------------------------------------------
// Common functions
// -----------------------------------------------------------------------------

// addPC increments the program counter by a signed byte offset
func (z80 *Z80) addPC(value uint16) {
	z80.PC = z80.PC + value
}

// decPC decrements the program counter
func (z80 *Z80) decPC() {
	z80.PC--
}

// decSP decrements the stack pointer
func (z80 *Z80) decSP() {
	z80.SP--
}

// incPC increments the program counter
func (z80 *Z80) incPC() {
	z80.PC++
}

// incR increments the refresh register ( 7 bit )
func (z80 *Z80) incR() {
	z80.R = (z80.R & 0x80) | ((z80.R + 1) & 0x7f)
}

// incSP increments the stack pointer
func (z80 *Z80) incSP() {
	z80.SP++
}

// setPC increments the program counter by a signed byte offset
func (z80 *Z80) setPC(address uint16) {
	z80.PC = address
}

// setSP increments the program counter by a signed byte offset
func (z80 *Z80) setSP(address uint16) {
	z80.SP = address
}

// -----------------------------------------------------------------------------
// Memory & IO functions
// -----------------------------------------------------------------------------

// readByte reads a byte from memory
func (z80 *Z80) readByte(address uint16) byte {
	data := z80.mem.Read(address)
	z80.clock.Add(3) // +3 tstates in memory access
	return data
}

// readBytePC reads a byte from pc address and increments PC
func (z80 *Z80) readBytePC() byte {
	data := z80.readByte(z80.PC)
	z80.incPC()
	return data
}

// readNoReq only puts address on bus (no MREQ) during n tstates
func (z80 *Z80) readNoReq(address uint16, n int) {
	for i := 0; i < n; i++ {
		z80.mem.Access(address)
		z80.clock.Inc() // +1 tstates in bus access
	}
}

// readPort writes a byte to a port
func (z80 *Z80) readPort(port uint16) byte {
	z80.readPortNoReq(port, 1) // pre
	z80.readPortNoReq(port, 3) // post
	value := z80.io.Read(port)
	return value
}

// readPortNoReq only puts address on bus (no IOREQ) during n tstates
func (z80 *Z80) readPortNoReq(address uint16, n int) {
	for i := 0; i < n; i++ {
		z80.io.Access(address)
		z80.clock.Inc() // +1 tstates in bus access
	}
}

// writeByte writes a byte into memory
func (z80 *Z80) writeByte(address uint16, value byte) {
	z80.mem.Write(address, value)
	z80.clock.Add(3) // +3 tstates in memory access
}

// writeNoReq only puts address on bus (no MREQ) during n tstates
func (z80 *Z80) writeNoReq(address uint16, n int) {
	for i := 0; i < n; i++ {
		z80.mem.Access(address)
		z80.clock.Inc() // +1 tstates in bus access
	}
}

// writePort writes a byte to a port
func (z80 *Z80) writePort(port uint16, value byte) {
	z80.readPortNoReq(port, 1) // pre
	z80.io.Write(port, value)
	z80.readPortNoReq(port, 3) // post
}

// -----------------------------------------------------------------------------
// Basic instructions
// -----------------------------------------------------------------------------

// adc acumulator add with carry
func (z80 *Z80) adc(value byte) {
	tmp := uint16(z80.A) + uint16(value) + uint16((z80.F & FlagC))
	lookup := ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte(((tmp & 0x88) >> 1))
	z80.A = byte(tmp)
	z80.F = ifthen((tmp&0x100) != 0, FlagC, 0) | halfCarryAddTable[lookup&0x07] |
		overflowAddTable[lookup>>4] | sz53Table[z80.A]
}

// add acumulator add
func (z80 *Z80) add(value byte) {
	tmp := uint16(z80.A) + uint16(value)
	lookup := ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte(((tmp & 0x88) >> 1))
	z80.A = byte(tmp)
	z80.F = ifthen((tmp&0x100) != 0, FlagC, 0) | halfCarryAddTable[lookup&0x07] |
		overflowAddTable[lookup>>4] | sz53Table[z80.A]
}

// and acumulator and
func (z80 *Z80) and(value byte) {
	z80.A &= value
	z80.F = FlagH | sz53pTable[z80.A]
}

// add16 16 bit add without carry
func (z80 *Z80) add16(reg1 *cpu.Register16, value uint16) {
	tmp := uint32(reg1.Get()) + uint32(value)
	lookup := byte((reg1.Get()&0x0800)>>11) | byte((value&0x0800)>>10) | byte((tmp&0x0800)>>9)
	z80.Memptr.Set(reg1.Get() + 1)
	reg1.Set(uint16(tmp))
	z80.F = (z80.F & (FlagV | FlagZ | FlagS)) | ifthen(tmp&0x10000 != 0, FlagC, 0) |
		(byte(tmp>>8) & (Flag3 | Flag5)) | halfCarryAddTable[lookup]
}

// call calls routine
func (z80 *Z80) call(condition bool) {
	z80.Z = z80.readBytePC()
	z80.W = z80.readByte(z80.PC)
	if condition {
		z80.readNoReq(z80.PC, 1)
		z80.incPC()
		z80.push8(highbyte(z80.PC))
		z80.push8(lowbyte(z80.PC))
		z80.setPC(z80.Memptr.Get())
	} else {
		z80.incPC()
	}
}

// cp acumulator sub
func (z80 *Z80) cp(value byte) {
	tmp := uint16(z80.A) - uint16(value)
	lookup := ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte(((tmp & 0x88) >> 1))
	z80.F = ifthen((tmp&0x100) != 0, FlagC, ifthen(tmp != 0, 0, FlagZ)) | FlagN |
		halfCarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] |
		(value & (Flag3 | Flag5)) | (lowbyte(tmp) & FlagS)
}

// daa decimal adjust acummulator
func (z80 *Z80) daa() {
	add, carry := byte(0), (z80.F & FlagC)
	if ((z80.F & FlagH) != 0) || ((z80.A & 0x0f) > 9) {
		add = 6
	}
	if (carry != 0) || (z80.A > 0x99) {
		add |= 0x60
	}
	if z80.A > 0x99 {
		carry = FlagC
	}
	if (z80.F & FlagN) != 0 {
		z80.sub(add)
	} else {
		z80.add(add)
	}
	z80.F = (z80.F & (^(FlagC | FlagP))) | carry | parityTable[z80.A]
}

// dec8 decreases a 8 bit register updating Flags
func (z80 *Z80) dec8(reg *cpu.Reg8) {
	z80.F = (z80.F & FlagC) | ifthen(*reg&0x0f != 0, 0, FlagH) | FlagN
	*reg--
	z80.F |= ifthen(*reg == 0x7f, FlagV, 0) | sz53Table[*reg]
}

// inc8 increases a 8 bit register updating Flags
func (z80 *Z80) inc8(reg *cpu.Reg8) {
	*reg++
	z80.F = (z80.F & FlagC) | ifthen(*reg == 0x80, FlagV, 0) |
		ifthen(*reg&0x0f != 0, 0, FlagH) | sz53Table[*reg]
}

// jp jump
func (z80 *Z80) jp(condition bool) {
	z80.Z = z80.readBytePC()
	z80.W = z80.readByte(z80.PC)
	if condition {
		z80.setPC(z80.Memptr.Get())
	} else {
		z80.incPC()
	}
}

// jr relative jump
func (z80 *Z80) jr(condition bool) {
	if condition {
		tmp := expandsign(z80.readByte(z80.PC))
		z80.readNoReq(z80.PC, 5)
		z80.addPC(tmp + 1)
		z80.Memptr.Set(z80.PC)
	} else {
		z80.readBytePC()
	}
}

// ld16nnrr writes 16 bit register to extended address
func (z80 *Z80) ld16nnrr(regL, regH cpu.Reg8) {
	tmp := toword(z80.readBytePC(), z80.readBytePC())
	z80.writeByte(tmp, regL)
	tmp++
	z80.Memptr.Set(tmp)
	z80.writeByte(tmp, regH)
}

// ld16rr loads 16 bit inmediate value
func (z80 *Z80) ld16rr(reg *cpu.Register16) {
	reg.SetL(z80.readBytePC())
	reg.SetH(z80.readBytePC())
}

// ld16rrnn loads 16 bit register from extended address
func (z80 *Z80) ld16rrnn(regL, regH *cpu.Reg8) {
	tmp := toword(z80.readBytePC(), z80.readBytePC())
	*regL = z80.readByte(tmp)
	tmp++
	z80.Memptr.Set(tmp)
	*regH = z80.readByte(tmp)
}

// or acumulator or
func (z80 *Z80) or(value byte) {
	z80.A |= value
	z80.F = sz53pTable[z80.A]
}

// pop8 pops a byte from stack
func (z80 *Z80) pop8() byte {
	value := z80.readByte(z80.SP)
	z80.incSP()
	return value
}

// pop16 pops a 16 bit register from stack
func (z80 *Z80) pop16(reg *cpu.Register16) {
	reg.SetL(z80.pop8())
	reg.SetH(z80.pop8())
}

// push8 push a byte to stack
func (z80 *Z80) push8(value byte) {
	z80.decSP()
	z80.writeByte(z80.SP, value)
}

// push16 push a 16 bit register to stack
func (z80 *Z80) push16(reg *cpu.Register16) {
	z80.push8(reg.GetH())
	z80.push8(reg.GetL())
}

// ret return from subrutine
func (z80 *Z80) ret(condition bool) {
	if condition {
		z80.setPC(toword(z80.pop8(), z80.pop8()))
		z80.Memptr.Set(z80.PC)
	}
}

// rst restart at address
func (z80 *Z80) rst(address uint16) {
	z80.push8(highbyte(z80.PC))
	z80.push8(lowbyte(z80.PC))
	z80.setPC(address)
	z80.Memptr.Set(z80.PC)
}

// sbc acumulator sub with carry
func (z80 *Z80) sbc(value byte) {
	tmp := uint16(z80.A) - uint16(value) - uint16((z80.F & FlagC))
	lookup := ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte(((tmp & 0x88) >> 1))
	z80.A = uint8(tmp)
	z80.F = ifthen((tmp&0x100) != 0, FlagC, 0) | FlagN |
		halfCarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] | sz53Table[z80.A]
}

// sub acumulator sub
func (z80 *Z80) sub(value byte) {
	tmp := uint16(z80.A) - uint16(value)
	lookup := ((z80.A & 0x88) >> 3) | ((value & 0x88) >> 2) | byte(((tmp & 0x88) >> 1))
	z80.A = uint8(tmp)
	z80.F = ifthen((tmp&0x100) != 0, FlagC, 0) | FlagN |
		halfCarrySubTable[lookup&0x07] | overflowSubTable[lookup>>4] | sz53Table[z80.A]
}

// xor acumulator xor
func (z80 *Z80) xor(value byte) {
	z80.A ^= value
	z80.F = sz53pTable[z80.A]
}

// -----------------------------------------------------------------------------
// Extended CB instructions
// -----------------------------------------------------------------------------

// bit test bit
func (z80 *Z80) bit(bit byte, reg cpu.Reg8) {
	z80.F = (z80.F & FlagC) | FlagH | (reg & (Flag3 | Flag5))
	if (reg & (0x01 << bit)) == 0 {
		z80.F |= FlagP | FlagZ
	}
	if bit == 7 && ((reg & 0x80) != 0) {
		z80.F |= FlagS
	}
}

// bitmemptr test bit memptr
func (z80 *Z80) bitmemptr(bit, value byte) {
	z80.F = (z80.F & FlagC) | FlagH | (z80.W & (Flag3 | Flag5))
	if (value & (0x01 << bit)) == 0 {
		z80.F |= FlagP | FlagZ
	}
	if bit == 7 && ((value & 0x80) != 0) {
		z80.F |= FlagS
	}
}

// rl rotate left with carry
func (z80 *Z80) rl(value *cpu.Reg8) {
	tmp := *value
	*value = (*value << 1) | (z80.F & FlagC)
	z80.F = (tmp >> 7) | sz53pTable[*value]
}

// rlc rotate left branch carry
func (z80 *Z80) rlc(value *cpu.Reg8) {
	*value = (*value << 1) | (*value >> 7)
	z80.F = (*value & FlagC) | sz53pTable[*value]
}

// rr rotate right with carry
func (z80 *Z80) rr(value *cpu.Reg8) {
	tmp := *value
	*value = (*value >> 1) | (z80.F << 7)
	z80.F = (tmp & FlagC) | sz53pTable[*value]
}

// rrc rotate right branch carry
func (z80 *Z80) rrc(value *cpu.Reg8) {
	z80.F = *value & FlagC
	*value = (*value >> 1) | (*value << 7)
	z80.F |= sz53pTable[*value]
}

// sla arithmetic left shift
func (z80 *Z80) sla(value *cpu.Reg8) {
	z80.F = *value >> 7
	*value <<= 1
	z80.F |= sz53pTable[*value]
}

// sll logical left shift
func (z80 *Z80) sll(value *cpu.Reg8) {
	z80.F = *value >> 7
	*value = (*value << 1) | 0x01
	z80.F |= sz53pTable[*value]
}

// sra arithmetic right shift
func (z80 *Z80) sra(value *cpu.Reg8) {
	z80.F = *value & FlagC
	*value = (*value & 0x80) | (*value >> 1)
	z80.F |= sz53pTable[*value]
}

// srl logical right shift
func (z80 *Z80) srl(value *cpu.Reg8) {
	z80.F = *value & FlagC
	*value >>= 1
	z80.F |= sz53pTable[*value]
}

// -----------------------------------------------------------------------------
// Extended ED instructions
// -----------------------------------------------------------------------------

// adc16 add with borrow HL
func (z80 *Z80) adc16(value uint16) {
	tmp := uint32(z80.HL.Get()) + uint32(value) + uint32((z80.F & FlagC))
	lookup := byte((z80.HL.Get()&0x8800)>>11) | byte(((value & 0x8800) >> 10)) | byte(((tmp & 0x8800) >> 9))
	z80.Memptr.Set(z80.HL.Get() + 1)
	z80.HL.Set(uint16(tmp))
	z80.F = ifthen((tmp&0x10000) != 0, FlagC, 0) | overflowAddTable[lookup>>4] |
		(z80.H & (Flag3 | Flag5 | FlagS)) | halfCarryAddTable[lookup&0x07] |
		ifthen(z80.HL.IsZero(), FlagZ, 0)
}

// cpd compare with decrement
func (z80 *Z80) cpd() {
	value := z80.readByte(z80.HL.Get())
	tmp := z80.A - value
	lookup := ((z80.A & 0x08) >> 3) | ((value & 0x08) >> 2) | ((tmp & 0x08) >> 1)
	z80.readNoReq(z80.HL.Get(), 5)
	z80.HL.Dec()
	z80.BC.Dec()
	z80.F = (z80.F & FlagC) | ifthen(z80.BC.IsZero(), FlagN, (FlagV|FlagN)) |
		halfCarrySubTable[lookup] | ifthen(tmp != 0, 0, FlagZ) | (tmp & FlagS)
	if (z80.F & FlagH) != 0 {
		tmp--
	}
	z80.F |= (tmp & Flag3) | ifthen((tmp&0x02) != 0, Flag5, 0)
	z80.Memptr.Dec()
}

// cpi compare with increment
func (z80 *Z80) cpi() {
	value := z80.readByte(z80.HL.Get())
	tmp := z80.A - value
	lookup := ((z80.A & 0x08) >> 3) | ((value & 0x08) >> 2) | ((tmp & 0x08) >> 1)
	z80.readNoReq(z80.HL.Get(), 5)
	z80.HL.Inc()
	z80.BC.Dec()
	z80.F = (z80.F & FlagC) | ifthen(z80.BC.IsZero(), FlagN, (FlagV|FlagN)) |
		halfCarrySubTable[lookup] | ifthen(tmp != 0, 0, FlagZ) | (tmp & FlagS)
	if (z80.F & FlagH) != 0 {
		tmp--
	}
	z80.F |= (tmp & Flag3) | ifthen((tmp&0x02) != 0, Flag5, 0)
	z80.Memptr.Inc()
}

// in input data from port to reg
func (z80 *Z80) in(reg *cpu.Reg8, port uint16) {
	z80.Memptr.Set(port + 1)
	*reg = z80.readPort(port)
	z80.F = (z80.F & FlagC) | sz53pTable[*reg]
}

// ind input with decrement
func (z80 *Z80) ind() {
	z80.readNoReq(z80.IR.Get(), 1)
	tmp := z80.readPort(z80.BC.Get())
	z80.writeByte(z80.HL.Get(), tmp)
	z80.Memptr.Set(z80.BC.Get() - 1)
	z80.B--
	z80.HL.Dec()
	tmp2 := tmp + z80.C - 1
	z80.F = ifthen((tmp&0x80) != 0, FlagN, 0) |
		ifthen(tmp2 < tmp, FlagH|FlagC, 0) |
		ifthen(parityTable[(tmp2&0x07)^z80.B] != 0, FlagP, 0) |
		sz53Table[z80.B]
}

// ini input with increment
func (z80 *Z80) ini() {
	z80.readNoReq(z80.IR.Get(), 1)
	tmp := z80.readPort(z80.BC.Get())
	z80.writeByte(z80.HL.Get(), tmp)
	z80.Memptr.Set(z80.BC.Get() + 1)
	z80.B--
	z80.HL.Inc()
	tmp2 := tmp + z80.C + 1
	z80.F = ifthen((tmp&0x80) != 0, FlagN, 0) |
		ifthen(tmp2 < tmp, FlagH|FlagC, 0) |
		ifthen(parityTable[(tmp2&0x07)^z80.B] != 0, FlagP, 0) |
		sz53Table[z80.B]
}

// ldd block load with decrement
func (z80 *Z80) ldd() {
	tmp := z80.readByte(z80.HL.Get())
	z80.BC.Dec()
	z80.writeByte(z80.DE.Get(), tmp)
	z80.writeNoReq(z80.DE.Get(), 2)
	z80.DE.Dec()
	z80.HL.Dec()
	tmp += z80.A
	z80.F = (z80.F & (FlagC | FlagZ | FlagS)) | ifthen(z80.BC.IsZero(), 0, FlagV) |
		(tmp & Flag3) | ifthen((tmp&0x02) != 0, Flag5, 0)
}

// ldi block load with increment
func (z80 *Z80) ldi() {
	tmp := z80.readByte(z80.HL.Get())
	z80.BC.Dec()
	z80.writeByte(z80.DE.Get(), tmp)
	z80.writeNoReq(z80.DE.Get(), 2)
	z80.DE.Inc()
	z80.HL.Inc()
	tmp += z80.A
	z80.F = (z80.F & (FlagC | FlagZ | FlagS)) | ifthen(z80.BC.IsZero(), 0, FlagV) |
		(tmp & Flag3) | ifthen((tmp&0x02) != 0, Flag5, 0)
}

// out output data from reg to port
func (z80 *Z80) out(port uint16, reg cpu.Reg8) {
	z80.writePort(port, reg)
	z80.Memptr.Set(port + 1)
}

// outd output with decrement
func (z80 *Z80) outd() {
	z80.readNoReq(z80.IR.Get(), 1)
	tmp := z80.readByte(z80.HL.Get())
	z80.B--
	z80.Memptr.Set(z80.BC.Get() - 1)
	z80.writePort(z80.BC.Get(), tmp)
	z80.HL.Dec()
	tmp2 := tmp + z80.L
	z80.F = ifthen((tmp&0x80) != 0, FlagN, 0) |
		ifthen(tmp2 < tmp, FlagH|FlagC, 0) |
		ifthen(parityTable[(tmp2&0x07)^z80.B] != 0, FlagP, 0) |
		sz53Table[z80.B]
}

// outd output with increment
func (z80 *Z80) outi() {
	z80.readNoReq(z80.IR.Get(), 1)
	tmp := z80.readByte(z80.HL.Get())
	z80.B--
	z80.Memptr.Set(z80.BC.Get() + 1)
	z80.writePort(z80.BC.Get(), tmp)
	z80.HL.Inc()
	tmp2 := tmp + z80.L
	z80.F = ifthen((tmp&0x80) != 0, FlagN, 0) |
		ifthen(tmp2 < tmp, FlagH|FlagC, 0) |
		ifthen(parityTable[(tmp2&0x07)^z80.B] != 0, FlagP, 0) |
		sz53Table[z80.B]
}

// sbc16 sub with borrow HL
func (z80 *Z80) sbc16(value uint16) {
	tmp := uint32(z80.HL.Get()) - uint32(value) - uint32((z80.F & FlagC))
	lookup := byte((z80.HL.Get()&0x8800)>>11) | byte(((value & 0x8800) >> 10)) | byte(((tmp & 0x8800) >> 9))
	z80.Memptr.Set(z80.HL.Get() + 1)
	z80.HL.Set(uint16(tmp))
	z80.F = ifthen((tmp&0x10000) != 0, FlagC, 0) | FlagN | overflowSubTable[lookup>>4] |
		(z80.H & (Flag3 | Flag5 | FlagS)) | halfCarrySubTable[lookup&0x07] |
		ifthen(z80.HL.IsZero(), FlagZ, 0)
}

// -----------------------------------------------------------------------------
// Register IX/IY extended instructions
// -----------------------------------------------------------------------------

// ld8rrixdd loads a value from an indexed address to a register
func (z80 *Z80) ld8rixdd(to *cpu.Reg8, reg *cpu.Register16) {
	offset := expandsign(z80.readByte(z80.PC))
	z80.readNoReq(z80.PC, 5)
	z80.incPC()
	address := reg.Get() + offset
	z80.Memptr.Set(address)
	*to = z80.readByte(reg.Get() + offset)
}

// ld8ixddrr loads a value to an indexed address
func (z80 *Z80) ld8ixddr(reg *cpu.Register16, value cpu.Reg8) {
	offset := expandsign(z80.readByte(z80.PC))
	z80.readNoReq(z80.PC, 5)
	z80.incPC()
	address := reg.Get() + offset
	z80.Memptr.Set(address)
	z80.writeByte(address, value)
}

// ac8ixdd acumulator op from indexed address value
func (z80 *Z80) ac8ixdd(reg *cpu.Register16, op8 func(byte)) {
	offset := expandsign(z80.readByte(z80.PC))
	z80.readNoReq(z80.PC, 5)
	z80.incPC()
	address := reg.Get() + offset
	z80.Memptr.Set(address)
	tmp := z80.readByte(address)
	op8(tmp)
}

// -----------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------

// expandsign the sign of a byte to an uint16
func expandsign(value byte) uint16 {
	return uint16(int16(int8(value)))
}

// highbyte gets the hight byte of a word
func highbyte(word uint16) byte {
	return byte(word >> 8)
}

// ifthen ternary operator for byte operands
func ifthen(condition bool, trueval, falseval byte) byte {
	if condition {
		return trueval
	}
	return falseval
}

// lowbyte gets the lower byte of a word
func lowbyte(word uint16) byte {
	return byte(word & 0xff)
}

// toword gets a word from low and high bytes
func toword(low, high byte) uint16 {
	return uint16(low) | (uint16(high) << 8)
}
