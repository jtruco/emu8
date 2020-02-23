package z80

import "github.com/jtruco/emu8/device/cpu"

// executeDD decodes and executes Z80 DD extended opcodes
func (z80 *Z80) executeDD(opcode byte) {
	z80.executeDDFD(opcode, &z80.IX)
}

// executeFD decodes and executes Z80 FD extended opcodes
func (z80 *Z80) executeFD(opcode byte) {
	z80.executeDDFD(opcode, &z80.IY)
}

// executeDDFD decodes and executes Z80 DD/FD extended opcodes
func (z80 *Z80) executeDDFD(opcode byte, reg *cpu.Register16) {

	switch opcode {

	case 0x09: // ADD REGISTER,BC
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(reg, z80.BC.Get())

	case 0x19: // ADD REGISTER,DE
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(reg, z80.DE.Get())

	case 0x21: // LD REGISTER,nnnn
		reg.SetL(z80.readBytePC())
		reg.SetH(z80.readBytePC())

	case 0x22: // LD (nnnn),REGISTER
		z80.ld16nnrr(reg.GetL(), reg.GetH())

	case 0x23: // INC REGISTER
		z80.readNoReq(z80.IR.Get(), 2)
		reg.Inc()

	case 0x24: // INC REGISTERH
		z80.inc8(reg.H())

	case 0x25: // DEC REGISTERH
		z80.dec8(reg.H())

	case 0x26: // LD REGISTERH,nn
		reg.SetH(z80.readBytePC())

	case 0x29: // ADD REGISTER,REGISTER
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(reg, reg.Get())

	case 0x2a: // LD REGISTER,(nnnn)
		z80.ld16rrnn(reg.L(), reg.H())

	case 0x2b: // DEC REGISTER
		z80.readNoReq(z80.IR.Get(), 2)
		reg.Dec()

	case 0x2c: // INC REGISTERL
		z80.inc8(reg.L())

	case 0x2d: // DEC REGISTERL
		z80.dec8(reg.L())

	case 0x2e: // LD REGISTERL,nn
		reg.SetL(z80.readBytePC())

	case 0x34: // INC (REGISTER+dd)
		offset := expandsign(z80.readByte(z80.PC))
		z80.readNoReq(z80.PC, 5)
		z80.incPC()
		z80.Memptr.Set(reg.Get() + offset)
		tmp := z80.readByte(z80.Memptr.Get())
		z80.readNoReq(z80.Memptr.Get(), 1)
		z80.inc8(&tmp)
		z80.writeByte(z80.Memptr.Get(), tmp)

	case 0x35: // DEC (REGISTER+dd)
		offset := expandsign(z80.readByte(z80.PC))
		z80.readNoReq(z80.PC, 5)
		z80.incPC()
		z80.Memptr.Set(reg.Get() + offset)
		tmp := z80.readByte(z80.Memptr.Get())
		z80.readNoReq(z80.Memptr.Get(), 1)
		z80.dec8(&tmp)
		z80.writeByte(z80.Memptr.Get(), tmp)

	case 0x36: // LD (REGISTER+dd),nn
		offset := expandsign(z80.readBytePC())
		tmp := z80.readByte(z80.PC)
		z80.readNoReq(z80.PC, 2)
		z80.incPC()
		z80.Memptr.Set(reg.Get() + offset)
		z80.writeByte(z80.Memptr.Get(), tmp)

	case 0x39: // ADD REGISTER,SP
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(reg, z80.SP)

	case 0x44: // LD B,REGISTERH
		z80.B = reg.GetH()

	case 0x45: // LD B,REGISTERL
		z80.B = reg.GetL()

	case 0x46: // LD B,(REGISTER+dd)
		z80.ld8rixdd(&z80.B, reg)

	case 0x4c: // LD C,REGISTERH
		z80.C = reg.GetH()

	case 0x4d: // LD C,REGISTERL
		z80.C = reg.GetL()

	case 0x4e: // LD C,(REGISTER+dd)
		z80.ld8rixdd(&z80.C, reg)

	case 0x54: // LD D,REGISTERH
		z80.D = reg.GetH()

	case 0x55: // LD D,REGISTERL
		z80.D = reg.GetL()

	case 0x56: // LD D,(REGISTER+dd)
		z80.ld8rixdd(&z80.D, reg)

	case 0x5c: // LD E,REGISTERH
		z80.E = reg.GetH()

	case 0x5d: // LD E,REGISTERL
		z80.E = reg.GetL()

	case 0x5e: // LD E,(REGISTER+dd)
		z80.ld8rixdd(&z80.E, reg)

	case 0x60: // LD REGISTERH,B
		reg.SetH(z80.B)

	case 0x61: // LD REGISTERH,C
		reg.SetH(z80.C)

	case 0x62: // LD REGISTERH,D
		reg.SetH(z80.D)

	case 0x63: // LD REGISTERH,E
		reg.SetH(z80.E)

	case 0x64: // LD REGISTERH,REGISTERH
		// NOP

	case 0x65: // LD REGISTERH,REGISTERL
		reg.SetH(reg.GetL())

	case 0x66: // LD H,(REGISTER+dd)
		z80.ld8rixdd(&z80.H, reg)

	case 0x67: // LD REGISTERH,A
		reg.SetH(z80.A)

	case 0x68: // LD REGISTERL,B
		reg.SetL(z80.B)

	case 0x69: // LD REGISTERL,C
		reg.SetL(z80.C)

	case 0x6a: // LD REGISTERL,D
		reg.SetL(z80.D)

	case 0x6b: // LD REGISTERL,E
		reg.SetL(z80.E)

	case 0x6c: // LD REGISTERL,REGISTERH
		reg.SetL(reg.GetH())

	case 0x6d: // LD REGISTERL,REGISTERL
		// NOP

	case 0x6e: // LD L,(REGISTER+dd)
		z80.ld8rixdd(&z80.L, reg)

	case 0x6f: // LD REGISTERL,A
		reg.SetL(z80.A)

	case 0x70: // LD (REGISTER+dd),B
		z80.ld8ixddr(reg, z80.B)

	case 0x71: // LD (REGISTER+dd),C
		z80.ld8ixddr(reg, z80.C)

	case 0x72: // LD (REGISTER+dd),D
		z80.ld8ixddr(reg, z80.D)

	case 0x73: // LD (REGISTER+dd),E
		z80.ld8ixddr(reg, z80.E)

	case 0x74: // LD (REGISTER+dd),H
		z80.ld8ixddr(reg, z80.H)

	case 0x75: // LD (REGISTER+dd),L
		z80.ld8ixddr(reg, z80.L)

	case 0x77: // LD (REGISTER+dd),A
		z80.ld8ixddr(reg, z80.A)

	case 0x7c: // LD A,REGISTERH
		z80.A = reg.GetH()

	case 0x7d: // LD A,REGISTERL
		z80.A = reg.GetL()

	case 0x7e: // LD A,(REGISTER+dd)
		z80.ld8rixdd(&z80.A, reg)

	case 0x84: // ADD A,REGISTERH
		z80.add(reg.GetH())

	case 0x85: // ADD A,REGISTERL
		z80.add(reg.GetL())

	case 0x86: // ADD A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.add)

	case 0x8c: // ADC A,REGISTERH
		z80.adc(reg.GetH())

	case 0x8d: // ADC A,REGISTERL
		z80.adc(reg.GetL())

	case 0x8e: // ADC A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.adc)

	case 0x94: // SUB A,REGISTERH
		z80.sub(reg.GetH())

	case 0x95: // SUB A,REGISTERL
		z80.sub(reg.GetL())

	case 0x96: // SUB A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.sub)

	case 0x9c: // SBC A,REGISTERH
		z80.sbc(reg.GetH())

	case 0x9d: // SBC A,REGISTERL
		z80.sbc(reg.GetL())

	case 0x9e: // SBC A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.sbc)

	case 0xa4: // AND A,REGISTERH
		z80.and(reg.GetH())

	case 0xa5: // AND A,REGISTERL
		z80.and(reg.GetL())

	case 0xa6: // AND A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.and)

	case 0xac: // XOR A,REGISTERH
		z80.xor(reg.GetH())

	case 0xad: // XOR A,REGISTERL
		z80.xor(reg.GetL())

	case 0xae: // XOR A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.xor)

	case 0xb4: // OR A,REGISTERH
		z80.or(reg.GetH())

	case 0xb5: // OR A,REGISTERL
		z80.or(reg.GetL())

	case 0xb6: // OR A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.or)

	case 0xbc: // CP A,REGISTERH
		z80.cp(reg.GetH())

	case 0xbd: // CP A,REGISTERL
		z80.cp(reg.GetL())

	case 0xbe: // CP A,(REGISTER+dd)
		z80.ac8ixdd(reg, z80.cp)

	case 0xcb: // shift DDFDCB
		offset := expandsign(z80.readByte(z80.PC))
		z80.Memptr.Set(reg.Get() + offset)
		z80.incPC()
		shiftopcode := z80.readByte(z80.PC)
		z80.readNoReq(z80.PC, 2)
		z80.incPC()
		z80.executeDDFDCB(shiftopcode, reg)

	case 0xe1: /* POP REGISTER */
		z80.pop16(reg)

	case 0xe3: // EX (SP),REGISTER
		tmpl := z80.readByte(z80.SP)
		tmph := z80.readByte(z80.SP + 1)
		z80.readNoReq(z80.SP+1, 1)
		z80.writeByte(z80.SP+1, reg.GetH())
		z80.writeByte(z80.SP, reg.GetL())
		z80.writeNoReq(z80.SP, 2)
		*reg.L(), z80.Z = tmpl, tmpl
		*reg.H(), z80.W = tmph, tmph

	case 0xe5: // PUSH REGISTER
		z80.readNoReq(z80.IR.Get(), 1)
		z80.push16(reg)

	case 0xe9: // JP REGISTER
		z80.setPC(reg.Get())

	case 0xf9: // LD SP,REGISTER
		z80.readNoReq(z80.IR.Get(), 2)
		z80.SP = reg.Get()

	default:
		// No IX/IY instruction : parse normal mode opcode
		z80.execute(opcode)
	}
}
