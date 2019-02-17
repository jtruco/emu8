package z80

// execute decodes and executes normal Z80 opcodes
func (z80 *Z80) execute(opcode byte) {

	switch opcode {

	case 0x00: // NOP

	case 0x01: // LD BC,nnnn
		z80.C = z80.readBytePC()
		z80.B = z80.readBytePC()

	case 0x02: // LD (BC),A
		z80.Z = lowbyte(z80.BC.Get() + 1)
		z80.W = z80.A
		z80.writeByte(z80.BC.Get(), z80.A)

	case 0x03: // INC BC
		z80.readNoReq(z80.IR.Get(), 2)
		z80.BC.Inc()

	case 0x04: // INC B
		z80.inc8(&z80.B)

	case 0x05: // DEC B
		z80.dec8(&z80.B)

	case 0x06: // LD B,nn
		z80.B = z80.readBytePC()

	case 0x07: // RLCA
		z80.A = (z80.A << 1) | (z80.A >> 7)
		z80.F = (z80.F & (FlagP | FlagZ | FlagS)) | (z80.A & (FlagC | Flag3 | Flag5))

	case 0x08: // EX AF,AF'
		z80.AF.Swap(&z80.AFx)

	case 0x09: // ADD HL,BC
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(&z80.HL, z80.BC.Get())

	case 0x0a: // LD A,(BC)
		z80.Memptr.Set(z80.BC.Get() + 1)
		z80.A = z80.readByte(z80.BC.Get())

	case 0x0b: // DEC BC
		z80.readNoReq(z80.IR.Get(), 2)
		z80.BC.Dec()

	case 0x0c: // INC C
		z80.inc8(&z80.C)

	case 0x0d: // DEC C
		z80.dec8(&z80.C)

	case 0x0e: // LD C,nn
		z80.C = z80.readBytePC()

	case 0x0f: // RRCA
		z80.F = (z80.F & (FlagP | FlagZ | FlagS)) | (z80.A & FlagC)
		z80.A = (z80.A >> 1) | (z80.A << 7)
		z80.F |= (z80.A & (Flag3 | Flag5))

	case 0x10: // DJNZ offset
		z80.readNoReq(z80.IR.Get(), 1)
		z80.B--
		z80.jr(z80.B != 0)

	case 0x11: // LD DE,nnnn
		z80.E = z80.readBytePC()
		z80.D = z80.readBytePC()

	case 0x12: // LD (DE),A
		z80.Z = lowbyte(z80.DE.Get() + 1)
		z80.W = z80.A
		z80.writeByte(z80.DE.Get(), z80.A)

	case 0x13: // INC DE
		z80.readNoReq(z80.IR.Get(), 2)
		z80.DE.Inc()

	case 0x14: // INC D
		z80.inc8(&z80.D)

	case 0x15: // DEC D
		z80.dec8(&z80.D)

	case 0x16: // LD D,nn
		z80.D = z80.readBytePC()

	case 0x17: // RLA
		tmp := z80.A
		z80.A = (z80.A << 1) | (z80.F & FlagC)
		z80.F = (z80.F & (FlagP | FlagZ | FlagS)) | (z80.A & (Flag3 | Flag5)) | (tmp >> 7)

	case 0x18: // JR offset
		z80.jr(true)

	case 0x19: // ADD HL,DE
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(&z80.HL, z80.DE.Get())

	case 0x1a: // LD A,(DE)
		z80.Memptr.Set(z80.DE.Get() + 1)
		z80.A = z80.readByte(z80.DE.Get())

	case 0x1b: // DEC DE
		z80.readNoReq(z80.IR.Get(), 2)
		z80.DE.Dec()

	case 0x1c: // INC E
		z80.inc8(&z80.E)

	case 0x1d: // DEC E
		z80.dec8(&z80.E)

	case 0x1e: // LD E,nn
		z80.E = z80.readBytePC()

	case 0x1f: // RRA
		tmp := z80.A
		z80.A = (z80.A >> 1) | (z80.F << 7)
		z80.F = (z80.F & (FlagP | FlagZ | FlagS)) | (z80.A & (Flag3 | Flag5)) | (tmp & FlagC)

	case 0x20: // JR NZ,offset
		z80.jr((z80.F & FlagZ) == 0)

	case 0x21: // LD HL,nnnn
		z80.L = z80.readBytePC()
		z80.H = z80.readBytePC()

	case 0x22: // LD (nnnn),HL
		z80.ld16nnrr(z80.L, z80.H)

	case 0x23: // INC HL
		z80.readNoReq(z80.IR.Get(), 2)
		z80.HL.Inc()

	case 0x24: // INC H
		z80.inc8(&z80.H)

	case 0x25: // DEC H
		z80.dec8(&z80.H)

	case 0x26: // LD H,nn
		z80.H = z80.readBytePC()

	case 0x27: // DAA
		z80.daa()

	case 0x28: // JR Z,offset
		z80.jr((z80.F & FlagZ) != 0)

	case 0x29: // ADD HL,HL
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(&z80.HL, z80.HL.Get())

	case 0x2a: // LD HL,(nnnn)
		z80.ld16rrnn(&z80.L, &z80.H)

	case 0x2b: // DEC HL
		z80.readNoReq(z80.IR.Get(), 2)
		z80.HL.Dec()

	case 0x2c: // INC L
		z80.inc8(&z80.L)

	case 0x2d: // DEC L
		z80.dec8(&z80.L)

	case 0x2e: // LD L,nn
		z80.L = z80.readBytePC()

	case 0x2f: // CPL
		z80.A ^= 0xff
		z80.F = (z80.F & (FlagC | FlagP | FlagZ | FlagS)) |
			(z80.A & (Flag3 | Flag5)) | (FlagN | FlagH)

	case 0x30: // JR NC,offset
		z80.jr((z80.F & FlagC) == 0)

	case 0x31: // LD SP,nnnn
		z80.setSP(toword(z80.readBytePC(), z80.readBytePC()))

	case 0x32: // LD (nnnn),A
		tmp := toword(z80.readBytePC(), z80.readBytePC())
		z80.Z = lowbyte(tmp + 1)
		z80.W = z80.A
		z80.writeByte(tmp, z80.A)

	case 0x33: // INC SP
		z80.readNoReq(z80.IR.Get(), 2)
		z80.incSP()

	case 0x34: // INC (HL)
		tmp := z80.readByte(z80.HL.Get())
		z80.readNoReq(z80.HL.Get(), 1)
		z80.inc8(&tmp)
		z80.writeByte(z80.HL.Get(), tmp)

	case 0x35: // DEC (HL)
		tmp := z80.readByte(z80.HL.Get())
		z80.readNoReq(z80.HL.Get(), 1)
		z80.dec8(&tmp)
		z80.writeByte(z80.HL.Get(), tmp)

	case 0x36: // LD (HL),nn
		z80.writeByte(z80.HL.Get(), z80.readBytePC())

	case 0x37: // SCF
		z80.F = (z80.F & (FlagP | FlagZ | FlagS)) | (z80.A & (Flag3 | Flag5)) | FlagC

	case 0x38: // JR C,offset
		z80.jr((z80.F & FlagC) != 0)

	case 0x39: // ADD HL,SP
		z80.readNoReq(z80.IR.Get(), 7)
		z80.add16(&z80.HL, z80.SP)

	case 0x3a: // LD A,(nnnn)
		z80.Z = z80.readBytePC()
		z80.W = z80.readBytePC()
		z80.A = z80.readByte(z80.Memptr.Get())
		z80.Memptr.Inc()

	case 0x3b: // DEC SP
		z80.readNoReq(z80.IR.Get(), 2)
		z80.decSP()

	case 0x3c: // INC A
		z80.inc8(&z80.A)

	case 0x3d: // DEC A
		z80.dec8(&z80.A)

	case 0x3e: // LD A,nn
		z80.A = z80.readBytePC()

	case 0x3f: // CCF
		z80.F = (z80.F & (FlagP | FlagZ | FlagS)) |
			ifthen((z80.F&FlagC) != 0, FlagH, FlagC) |
			(z80.A & (Flag3 | Flag5))

	case 0x40: // LD B,B

	case 0x41: // LD B,C
		z80.B = z80.C

	case 0x42: // LD B,D
		z80.B = z80.D

	case 0x43: // LD B,E
		z80.B = z80.E

	case 0x44: // LD B,H
		z80.B = z80.H

	case 0x45: // LD B,L
		z80.B = z80.L

	case 0x46: // LD B,(HL)
		z80.B = z80.readByte(z80.HL.Get())

	case 0x47: // LD B,A
		z80.B = z80.A

	case 0x48: // LD C,B
		z80.C = z80.B

	case 0x49: // LD C,C

	case 0x4a: // LD C,D
		z80.C = z80.D

	case 0x4b: // LD C,E
		z80.C = z80.E

	case 0x4c: // LD C,H
		z80.C = z80.H

	case 0x4d: // LD C,L
		z80.C = z80.L

	case 0x4e: // LD C,(HL)
		z80.C = z80.readByte(z80.HL.Get())

	case 0x4f: // LD C,A
		z80.C = z80.A

	case 0x50: // LD D,B
		z80.D = z80.B

	case 0x51: // LD D,C
		z80.D = z80.C

	case 0x52: // LD D,D

	case 0x53: // LD D,E
		z80.D = z80.E

	case 0x54: // LD D,H
		z80.D = z80.H

	case 0x55: // LD D,L
		z80.D = z80.L

	case 0x56: // LD D,(HL)
		z80.D = z80.readByte(z80.HL.Get())

	case 0x57: // LD D,A
		z80.D = z80.A

	case 0x58: // LD E,B
		z80.E = z80.B

	case 0x59: // LD E,C
		z80.E = z80.C

	case 0x5a: // LD E,D
		z80.E = z80.D

	case 0x5b: // LD E,E

	case 0x5c: // LD E,H
		z80.E = z80.H

	case 0x5d: // LD E,L
		z80.E = z80.L

	case 0x5e: // LD E,(HL)
		z80.E = z80.readByte(z80.HL.Get())

	case 0x5f: // LD E,A
		z80.E = z80.A

	case 0x60: // LD H,B
		z80.H = z80.B

	case 0x61: // LD H,C
		z80.H = z80.C

	case 0x62: // LD H,D
		z80.H = z80.D

	case 0x63: // LD H,E
		z80.H = z80.E

	case 0x64: // LD H,H

	case 0x65: // LD H,L
		z80.H = z80.L

	case 0x66: // LD H,(HL)
		z80.H = z80.readByte(z80.HL.Get())

	case 0x67: // LD H,A
		z80.H = z80.A

	case 0x68: // LD L,B
		z80.L = z80.B

	case 0x69: // LD L,C
		z80.L = z80.C

	case 0x6a: // LD L,D
		z80.L = z80.D

	case 0x6b: // LD L,E
		z80.L = z80.E

	case 0x6c: // LD L,H
		z80.L = z80.H

	case 0x6d: // LD L,L

	case 0x6e: // LD L,(HL)
		z80.L = z80.readByte(z80.HL.Get())

	case 0x6f: // LD L,A
		z80.L = z80.A

	case 0x70: // LD (HL),B
		z80.writeByte(z80.HL.Get(), z80.B)

	case 0x71: // LD (HL),C
		z80.writeByte(z80.HL.Get(), z80.C)

	case 0x72: // LD (HL),D
		z80.writeByte(z80.HL.Get(), z80.D)

	case 0x73: // LD (HL),E
		z80.writeByte(z80.HL.Get(), z80.E)

	case 0x74: // LD (HL),H
		z80.writeByte(z80.HL.Get(), z80.H)

	case 0x75: // LD (HL),L
		z80.writeByte(z80.HL.Get(), z80.L)

	case 0x76: // HALT
		z80.Halted = true
		z80.decPC()

	case 0x77: // LD (HL),A
		z80.writeByte(z80.HL.Get(), z80.A)

	case 0x78: // LD A,B
		z80.A = z80.B

	case 0x79: // LD A,C
		z80.A = z80.C

	case 0x7a: // LD A,D
		z80.A = z80.D

	case 0x7b: // LD A,E
		z80.A = z80.E

	case 0x7c: // LD A,H
		z80.A = z80.H

	case 0x7d: // LD A,L
		z80.A = z80.L

	case 0x7e: // LD A,(HL)
		z80.A = z80.readByte(z80.HL.Get())

	case 0x7f: // LD A,A

	case 0x80: // ADD A,B
		z80.add(z80.B)

	case 0x81: // ADD A,C
		z80.add(z80.C)

	case 0x82: // ADD A,D
		z80.add(z80.D)

	case 0x83: // ADD A,E
		z80.add(z80.E)

	case 0x84: // ADD A,H
		z80.add(z80.H)

	case 0x85: // ADD A,L
		z80.add(z80.L)

	case 0x86: // ADD A,(HL)
		z80.add(z80.readByte(z80.HL.Get()))

	case 0x87: // ADD A,A
		z80.add(z80.A)

	case 0x88: // ADC A,B
		z80.adc(z80.B)

	case 0x89: // ADC A,C
		z80.adc(z80.C)

	case 0x8a: // ADC A,D
		z80.adc(z80.D)

	case 0x8b: // ADC A,E
		z80.adc(z80.E)

	case 0x8c: // ADC A,H
		z80.adc(z80.H)

	case 0x8d: // ADC A,L
		z80.adc(z80.L)

	case 0x8e: // ADC A,(HL)
		z80.adc(z80.readByte(z80.HL.Get()))

	case 0x8f: // ADC A,A
		z80.adc(z80.A)

	case 0x90: // SUB A,B
		z80.sub(z80.B)

	case 0x91: // SUB A,C
		z80.sub(z80.C)

	case 0x92: // SUB A,D
		z80.sub(z80.D)

	case 0x93: // SUB A,E
		z80.sub(z80.E)

	case 0x94: // SUB A,H
		z80.sub(z80.H)

	case 0x95: // SUB A,L
		z80.sub(z80.L)

	case 0x96: // SUB A,(HL)
		z80.sub(z80.readByte(z80.HL.Get()))

	case 0x97: // SUB A,A
		z80.sub(z80.A)

	case 0x98: // SBC A,B
		z80.sbc(z80.B)

	case 0x99: // SBC A,C
		z80.sbc(z80.C)

	case 0x9a: // SBC A,D
		z80.sbc(z80.D)

	case 0x9b: // SBC A,E
		z80.sbc(z80.E)

	case 0x9c: // SBC A,H
		z80.sbc(z80.H)

	case 0x9d: // SBC A,L
		z80.sbc(z80.L)

	case 0x9e: // SBC A,(HL)
		z80.sbc(z80.readByte(z80.HL.Get()))

	case 0x9f: // SBC A,A
		z80.sbc(z80.A)

	case 0xa0: // AND A,B
		z80.and(z80.B)

	case 0xa1: // AND A,C
		z80.and(z80.C)

	case 0xa2: // AND A,D
		z80.and(z80.D)

	case 0xa3: // AND A,E
		z80.and(z80.E)

	case 0xa4: // AND A,H
		z80.and(z80.H)

	case 0xa5: // AND A,L
		z80.and(z80.L)

	case 0xa6: // AND A,(HL)
		z80.and(z80.readByte(z80.HL.Get()))

	case 0xa7: // AND A,A
		z80.and(z80.A)

	case 0xa8: // XOR A,B
		z80.xor(z80.B)

	case 0xa9: // XOR A,C
		z80.xor(z80.C)

	case 0xaa: // XOR A,D
		z80.xor(z80.D)

	case 0xab: // XOR A,E
		z80.xor(z80.E)

	case 0xac: // XOR A,H
		z80.xor(z80.H)

	case 0xad: // XOR A,L
		z80.xor(z80.L)

	case 0xae: // XOR A,(HL)
		z80.xor(z80.readByte(z80.HL.Get()))

	case 0xaf: // XOR A,A
		z80.xor(z80.A)

	case 0xb0: // OR A,B
		z80.or(z80.B)

	case 0xb1: // OR A,C
		z80.or(z80.C)

	case 0xb2: // OR A,D
		z80.or(z80.D)

	case 0xb3: // OR A,E
		z80.or(z80.E)

	case 0xb4: // OR A,H
		z80.or(z80.H)

	case 0xb5: // OR A,L
		z80.or(z80.L)

	case 0xb6: // OR A,(HL)
		z80.or(z80.readByte(z80.HL.Get()))

	case 0xb7: // OR A,A
		z80.or(z80.A)

	case 0xb8: // CP B
		z80.cp(z80.B)

	case 0xb9: // CP C
		z80.cp(z80.C)

	case 0xba: // CP D
		z80.cp(z80.D)

	case 0xbb: // CP E
		z80.cp(z80.E)

	case 0xbc: // CP H
		z80.cp(z80.H)

	case 0xbd: // CP L
		z80.cp(z80.L)

	case 0xbe: // CP (HL)
		z80.cp(z80.readByte(z80.HL.Get()))

	case 0xbf: // CP A
		z80.cp(z80.A)

	case 0xc0: // RET NZ
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagZ) == 0)

	case 0xc1: // POP BC
		z80.pop16(&z80.BC)

	case 0xc2: // JP NZ,nnnn
		z80.jp((z80.F & FlagZ) == 0)

	case 0xc3: // JP nnnn
		z80.jp(true)

	case 0xc4: // CALL NZ,nnnn
		z80.call((z80.F & FlagZ) == 0)

	case 0xc5: // PUSH BC
		z80.readNoReq(z80.IR.Get(), 1)
		z80.push16(&z80.BC)

	case 0xc6: // ADD A,nn
		z80.add(z80.readBytePC())

	case 0xc7: // RST 00
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x00)

	case 0xc8: // RET Z
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagZ) != 0)

	case 0xc9: // RET
		z80.ret(true)

	case 0xca: // JP Z,nnnn
		z80.jp((z80.F & FlagZ) != 0)

	case 0xcb: // shift CB
		z80.fetchAndExecute(z80.executeCB)

	case 0xcc: // CALL Z,nnnn
		z80.call((z80.F & FlagZ) != 0)

	case 0xcd: // CALL nnnn
		z80.call(true)

	case 0xce: // ADC A,nn
		z80.adc(z80.readBytePC())

	case 0xcf: // RST 8
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x08)

	case 0xd0: // RET NC
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagC) == 0)

	case 0xd1: // POP DE
		z80.pop16(&z80.DE)

	case 0xd2: // JP NC,nnnn
		z80.jp((z80.F & FlagC) == 0)

	case 0xd3: // OUT (nn),A
		nn := z80.readBytePC()
		tmp := toword(nn, z80.A)
		z80.W = z80.A
		z80.Z = nn + 1
		z80.writePort(tmp, z80.A)

	case 0xd4: // CALL NC,nnnn
		z80.call((z80.F & FlagC) == 0)

	case 0xd5: // PUSH DE
		z80.readNoReq(z80.IR.Get(), 1)
		z80.push16(&z80.DE)

	case 0xd6: // SUB nn
		z80.sub(z80.readBytePC())

	case 0xd7: // RST 10
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x10)

	case 0xd8: // RET C
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagC) != 0)

	case 0xd9: // EXX
		z80.BC.Swap(&z80.BCx)
		z80.DE.Swap(&z80.DEx)
		z80.HL.Swap(&z80.HLx)

	case 0xda: // JP C,nnnn
		z80.jp((z80.F & FlagC) != 0)

	case 0xdb: // IN A,(nn)
		tmp := toword(z80.readBytePC(), z80.A)
		z80.A = z80.readPort(tmp)
		z80.Memptr.Set(tmp + 1)

	case 0xdc: // CALL C,nnnn
		z80.call((z80.F & FlagC) != 0)

	case 0xdd: // shift DD
		z80.fetchAndExecute(z80.executeDD)

	case 0xde: // SBC A,nn
		z80.sbc(z80.readBytePC())

	case 0xdf: // RST 18
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x18)

	case 0xe0: // RET PO
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagP) == 0)

	case 0xe1: // POP HL
		z80.pop16(&z80.HL)

	case 0xe2: // JP PO,nnnn
		z80.jp((z80.F & FlagP) == 0)

	case 0xe3: // EX (SP),HL
		tmpl := z80.readByte(z80.SP)
		tmph := z80.readByte(z80.SP + 1)
		z80.readNoReq(z80.SP+1, 1)
		z80.writeByte(z80.SP+1, z80.H)
		z80.writeByte(z80.SP, z80.L)
		z80.writeNoReq(z80.SP, 2)
		z80.L, z80.Z = tmpl, tmpl
		z80.H, z80.W = tmph, tmph

	case 0xe4: // CALL PO,nnnn
		z80.call((z80.F & FlagP) == 0)

	case 0xe5: // PUSH HL
		z80.readNoReq(z80.IR.Get(), 1)
		z80.push16(&z80.HL)

	case 0xe6: // AND nn
		z80.and(z80.readBytePC())

	case 0xe7: // RST 20
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x20)

	case 0xe8: // RET PE
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagP) != 0)

	case 0xe9: // JP HL
		z80.setPC(z80.HL.Get())

	case 0xea: // JP PE,nnnn
		z80.jp((z80.F & FlagP) != 0)

	case 0xeb: // EX DE,HL
		z80.DE.Swap(&z80.HL)

	case 0xec: // CALL PE,nnnn
		z80.call((z80.F & FlagP) != 0)

	case 0xed: // shift ED
		z80.fetchAndExecute(z80.executeED)

	case 0xee: // XOR A,nn
		z80.xor(z80.readBytePC())

	case 0xef: // RST 28
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x28)

	case 0xf0: // RET P
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagS) == 0)

	case 0xf1: // POP AF
		z80.pop16(&z80.AF)

	case 0xf2: // JP P,nnnn
		z80.jp((z80.F & FlagS) == 0)

	case 0xf3: // DI
		z80.IFF1, z80.IFF2 = false, false

	case 0xf4: // CALL P,nnnn
		z80.call((z80.F & FlagS) == 0)

	case 0xf5: // PUSH AF
		z80.readNoReq(z80.IR.Get(), 1)
		z80.push16(&z80.AF)

	case 0xf6: // OR nn
		z80.or(z80.readBytePC())

	case 0xf7: // RST 30
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x30)

	case 0xf8: // RET M
		z80.readNoReq(z80.IR.Get(), 1)
		z80.ret((z80.F & FlagS) != 0)

	case 0xf9: // LD SP,HL
		z80.readNoReq(z80.IR.Get(), 2)
		z80.setSP(z80.HL.Get())

	case 0xfa: // JP M,nnnn
		z80.jp((z80.F & FlagS) != 0)

	case 0xfb: // EI
		z80.IFF1, z80.IFF2 = true, true
		z80.ActiveEI = true

	case 0xfc: // CALL M,nnnn
		z80.call((z80.F & FlagS) != 0)

	case 0xfd: // shift FD
		z80.fetchAndExecute(z80.executeFD)

	case 0xfe: // CP nn
		z80.cp(z80.readBytePC())

	case 0xff: // RST 38
		z80.readNoReq(z80.IR.Get(), 1)
		z80.rst(0x38)
	}
}
