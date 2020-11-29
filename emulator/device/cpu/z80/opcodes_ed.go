package z80

// executeED decodes and executes Z80 ED extended opcodes
func (z80 *Z80) executeED(opcode byte) {

	switch opcode {

	case 0x40: // IN B,(C)
		z80.in(&z80.B, z80.BC.Get())

	case 0x41: // OUT (C),B
		z80.out(z80.BC.Get(), z80.B)

	case 0x42: // SBC HL,BC
		z80.readNoReq(z80.IR.Get(), 7)
		z80.sbc16(z80.BC.Get())

	case 0x43: // LD (nnnn),BC
		z80.ld16nnrr(z80.C, z80.B)

	case 0x44, 0x4c, 0x54, 0x5c, 0x64, 0x6c, 0x74, 0x7c: // NEG
		tmp := z80.A
		z80.A = 0
		z80.sub(tmp)

	case 0x45, 0x4d, 0x55, 0x5d, 0x65, 0x6d, 0x75, 0x7d: // RETN
		z80.IFF1 = z80.IFF2
		z80.ret(true)

	case 0x46, 0x4e, 0x66, 0x6e: // IM 0
		z80.IM = 0

	case 0x47: // LD I,A
		z80.readNoReq(z80.IR.Get(), 1)
		z80.I = z80.A

	case 0x48: // IN C,(C)
		z80.in(&z80.C, z80.BC.Get())

	case 0x49: // OUT (C),C
		z80.out(z80.BC.Get(), z80.C)

	case 0x4a: // ADC HL,BC
		z80.readNoReq(z80.IR.Get(), 7)
		z80.adc16(z80.BC.Get())

	case 0x4b: // LD BC,(nnnn)
		z80.ld16rrnn(&z80.C, &z80.B)

	case 0x4f: // LD R,A
		z80.readNoReq(z80.IR.Get(), 1)
		z80.R = z80.A

	case 0x50: // IN D,(C)
		z80.in(&z80.D, z80.BC.Get())

	case 0x51: // OUT (C),D
		z80.out(z80.BC.Get(), z80.D)

	case 0x52: // SBC HL,DE
		z80.readNoReq(z80.IR.Get(), 7)
		z80.sbc16(z80.DE.Get())

	case 0x53: // LD (nnnn),DE
		z80.ld16nnrr(z80.E, z80.D)

	case 0x56, 0x76: // IM 1
		z80.IM = 1

	case 0x57: // LD A,I
		z80.readNoReq(z80.IR.Get(), 1)
		z80.A = z80.I
		z80.F = (z80.F & FlagC) | sz53Table[z80.A] | ifthen(z80.IFF2, FlagV, 0)
		z80.ReadIFF2 = true

	case 0x58: // IN E,(C)
		z80.in(&z80.E, z80.BC.Get())

	case 0x59: // OUT (C),E
		z80.out(z80.BC.Get(), z80.E)

	case 0x5a: // ADC HL,DE
		z80.readNoReq(z80.IR.Get(), 7)
		z80.adc16(z80.DE.Get())

	case 0x5b: // LD DE,(nnnn)
		z80.ld16rrnn(&z80.E, &z80.D)

	case 0x5e:
	case 0x7e: // IM 2
		z80.IM = 2

	case 0x5f: // LD A,R
		z80.readNoReq(z80.IR.Get(), 1)
		z80.A = z80.R
		z80.F = (z80.F & FlagC) | sz53Table[z80.A] | ifthen(z80.IFF2, FlagV, 0)
		z80.ReadIFF2 = true

	case 0x60: // IN H,(C)
		z80.in(&z80.H, z80.BC.Get())

	case 0x61: // OUT (C),H
		z80.out(z80.BC.Get(), z80.H)

	case 0x62: // SBC HL,HL
		z80.readNoReq(z80.IR.Get(), 7)
		z80.sbc16(z80.HL.Get())

	case 0x63: // LD (nnnn),HL
		z80.ld16nnrr(z80.L, z80.H)

	case 0x67: // RRD
		tmp := z80.readByte(z80.HL.Get())
		z80.readNoReq(z80.HL.Get(), 4)
		z80.writeByte(z80.HL.Get(), ((z80.A << 4) | (tmp >> 4)))
		z80.A = (z80.A & 0xf0) | (tmp & 0x0f)
		z80.F = (z80.F & FlagC) | sz53pTable[z80.A]
		z80.Memptr.Set(z80.HL.Get() + 1)

	case 0x68: // IN L,(C)
		z80.in(&z80.L, z80.BC.Get())

	case 0x69: // OUT (C),L
		z80.out(z80.BC.Get(), z80.L)

	case 0x6a: // ADC HL,HL
		z80.readNoReq(z80.IR.Get(), 7)
		z80.adc16(z80.HL.Get())

	case 0x6b: // LD HL,(nnnn)
		z80.ld16rrnn(&z80.L, &z80.H)

	case 0x6f: // RLD
		tmp := z80.readByte(z80.HL.Get())
		z80.readNoReq(z80.HL.Get(), 4)
		z80.writeByte(z80.HL.Get(), ((tmp << 4) | (z80.A & 0x0f)))
		z80.A = (z80.A & 0xf0) | (tmp >> 4)
		z80.F = (z80.F & FlagC) | sz53pTable[z80.A]
		z80.Memptr.Set(z80.HL.Get() + 1)

	case 0x70: // IN F,(C)
		var tmp byte
		z80.in(&tmp, z80.BC.Get())

	case 0x71: // OUT (C),0
		z80.out(z80.BC.Get(), 0) // NMOS

	case 0x72: // SBC HL,SP
		z80.readNoReq(z80.IR.Get(), 7)
		z80.sbc16(z80.SP)

	case 0x73: // LD (nnnn),SP
		z80.ld16nnrr(lowbyte(z80.SP), highbyte(z80.SP))

	case 0x78: // IN A,(C)
		z80.in(&z80.A, z80.BC.Get())

	case 0x79: // OUT (C),A
		z80.out(z80.BC.Get(), z80.A)

	case 0x7a: // ADC HL,SP
		z80.readNoReq(z80.IR.Get(), 7)
		z80.adc16(z80.SP)

	case 0x7b: // LD SP,(nnnn)
		var tmpl, tmph byte
		z80.ld16rrnn(&tmpl, &tmph)
		z80.SP = toword(tmpl, tmph)

	case 0xa0: // LDI
		z80.ldi()

	case 0xa1: // CPI
		z80.cpi()

	case 0xa2: // INI
		z80.ini()

	case 0xa3: // OUTI
		z80.outi()

	case 0xa8: // LDD
		z80.ldd()

	case 0xa9: // CPD
		z80.cpd()

	case 0xaa: // IND
		z80.ind()

	case 0xab: // OUTD
		z80.outd()

	case 0xb0: // LDIR
		z80.ldi()
		if !z80.BC.IsZero() {
			z80.writeNoReq(z80.DE.Get()-1, 5)
			z80.addPC(0xfffe) // PC -= 2
			z80.Memptr.Set(z80.PC + 1)
		}

	case 0xb1: // CPIR
		z80.cpi()
		if (z80.F & (FlagV | FlagZ)) == FlagV {
			z80.readNoReq(z80.HL.Get()-1, 5)
			z80.addPC(0xfffe) // PC -= 2
			z80.Memptr.Set(z80.PC + 1)
		}

	case 0xb2: // INIR
		z80.ini()
		if z80.B != 0 {
			z80.writeNoReq(z80.HL.Get()-1, 5)
			z80.addPC(0xfffe) // PC -= 2
		}

	case 0xb3: // OTIR
		z80.outi()
		if z80.B != 0 {
			z80.readNoReq(z80.BC.Get(), 5)
			z80.addPC(0xfffe) // PC -= 2
		}

	case 0xb8: // LDDR
		z80.ldd()
		if !z80.BC.IsZero() {
			z80.writeNoReq(z80.DE.Get()+1, 5)
			z80.addPC(0xfffe) // PC -= 2
			z80.Memptr.Set(z80.PC + 1)
		}

	case 0xb9: // CPDR
		z80.cpd()
		if (z80.F & (FlagV | FlagZ)) == FlagV {
			z80.readNoReq(z80.HL.Get()+1, 5)
			z80.addPC(0xfffe) // PC -= 2
			z80.Memptr.Set(z80.PC + 1)
		}

	case 0xba: // INDR
		z80.ind()
		if z80.B != 0 {
			z80.writeNoReq(z80.HL.Get()-1, 5)
			z80.addPC(0xfffe) // PC -= 2
		}

	case 0xbb: // OTDR
		z80.outd()
		if z80.B != 0 {
			z80.readNoReq(z80.BC.Get(), 5)
			z80.addPC(0xfffe) // PC -= 2
		}

	}
}
