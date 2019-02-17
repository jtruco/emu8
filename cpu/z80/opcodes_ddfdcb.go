package z80

import (
	"github.com/jtruco/emu8/cpu"
)

// executeDDFDCB decodes and executes Z80 DDFFD/CB extended opcodes
func (z80 *Z80) executeDDFDCB(opcode byte, reg *cpu.Register16) {

	// memptr : reg + dd (address)
	memptr := z80.Memptr.Get()

	switch opcode {

	case 0x00: // LD B,RLC (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x01: // LD C,RLC (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x02: // LD D,RLC (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x03: // LD E,RLC (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x04: // LD H,RLC (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x05: // LD L,RLC (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x06: // RLC (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x07: // LD A,RLC (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rlc(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x08: // LD B,RRC (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x09: // LD C,RRC (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x0a: // LD D,RRC (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x0b: // LD E,RRC (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x0c: // LD H,RRC (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x0d: // LD L,RRC (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x0e: // RRC (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x0f: // LD A,RRC (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rrc(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x10: // LD B,RL (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x11: // LD C,RL (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x12: // LD D,RL (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x13: // LD E,RL (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x14: // LD H,RL (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x15: // LD L,RL (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x16: // RL (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x17: // LD A,RL (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rl(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x18: // LD B,RR (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x19: // LD C,RR (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x1a: // LD D,RR (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x1b: // LD E,RR (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x1c: // LD H,RR (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x1d: // LD L,RR (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x1e: // RR (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x1f: // LD A,RR (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.rr(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x20: // LD B,SLA (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x21: // LD C,SLA (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x22: // LD D,SLA (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x23: // LD E,SLA (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x24: // LD H,SLA (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x25: // LD L,SLA (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x26: // SLA (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x27: // LD A,SLA (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sla(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x28: // LD B,SRA (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x29: // LD C,SRA (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x2a: // LD D,SRA (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x2b: // LD E,SRA (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x2c: // LD H,SRA (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x2d: // LD L,SRA (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x2e: // SRA (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x2f: // LD A,SRA (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sra(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x30: // LD B,SLL (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x31: // LD C,SLL (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x32: // LD D,SLL (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x33: // LD E,SLL (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x34: // LD H,SLL (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x35: // LD L,SLL (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x36: // SLL (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x37: // LD A,SLL (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.sll(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x38: // LD B,SRL (REGISTER+dd)
		z80.B = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&z80.B)
		z80.writeByte(memptr, z80.B)

	case 0x39: // LD C,SRL (REGISTER+dd)
		z80.C = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&z80.C)
		z80.writeByte(memptr, z80.C)

	case 0x3a: // LD D,SRL (REGISTER+dd)
		z80.D = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&z80.D)
		z80.writeByte(memptr, z80.D)

	case 0x3b: // LD E,SRL (REGISTER+dd)
		z80.E = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&z80.E)
		z80.writeByte(memptr, z80.E)

	case 0x3c: // LD H,SRL (REGISTER+dd)
		z80.H = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&z80.H)
		z80.writeByte(memptr, z80.H)

	case 0x3d: // LD L,SRL (REGISTER+dd)
		z80.L = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&z80.L)
		z80.writeByte(memptr, z80.L)

	case 0x3e: // SRL (REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&tmp)
		z80.writeByte(memptr, tmp)

	case 0x3f: // LD A,SRL (REGISTER+dd)
		z80.A = z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.srl(&z80.A)
		z80.writeByte(memptr, z80.A)

	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47: // BIT 0,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(0, tmp)

	case 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f: // BIT 1,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(1, tmp)

	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57: // BIT 2,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(2, tmp)

	case 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f: // BIT 3,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(3, tmp)

	case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67: // BIT 4,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(4, tmp)

	case 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f: // BIT 5,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(5, tmp)

	case 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77: // BIT 6,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(6, tmp)

	case 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f: // BIT 7,(REGISTER+dd)
		tmp := z80.readByte(memptr)
		z80.readNoReq(memptr, 1)
		z80.bitmemptr(7, tmp)

	case 0x80: // LD B,RES 0,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0x81: // LD C,RES 0,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0x82: // LD D,RES 0,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0x83: // LD E,RES 0,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0x84: // LD H,RES 0,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0x85: // LD L,RES 0,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0x86: // RES 0,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0x87: // LD A,RES 0,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0xfe
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0x88: // LD B,RES 1,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0x89: // LD C,RES 1,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0x8a: // LD D,RES 1,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0x8b: // LD E,RES 1,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0x8c: // LD H,RES 1,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0x8d: // LD L,RES 1,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0x8e: // RES 1,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0x8f: // LD A,RES 1,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0xfd
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0x90: // LD B,RES 2,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0x91: // LD C,RES 2,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0x92: // LD D,RES 2,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0x93: // LD E,RES 2,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0x94: // LD H,RES 2,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0x95: // LD L,RES 2,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0x96: // RES 2,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0x97: // LD A,RES 2,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0xfb
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0x98: // LD B,RES 3,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0x99: // LD C,RES 3,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0x9a: // LD D,RES 3,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0x9b: // LD E,RES 3,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0x9c: // LD H,RES 3,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0x9d: // LD L,RES 3,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0x9e: // RES 3,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0x9f: // LD A,RES 3,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0xf7
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xa0: // LD B,RES 4,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xa1: // LD C,RES 4,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xa2: // LD D,RES 4,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xa3: // LD E,RES 4,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xa4: // LD H,RES 4,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xa5: // LD L,RES 4,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xa6: // RES 4,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xa7: // LD A,RES 4,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0xef
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xa8: // LD B,RES 5,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xa9: // LD C,RES 5,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xaa: // LD D,RES 5,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xab: // LD E,RES 5,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xac: // LD H,RES 5,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xad: // LD L,RES 5,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xae: // RES 5,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xaf: // LD A,RES 5,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0xdf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xb0: // LD B,RES 6,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xb1: // LD C,RES 6,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xb2: // LD D,RES 6,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xb3: // LD E,RES 6,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xb4: // LD H,RES 6,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xb5: // LD L,RES 6,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xb6: // RES 6,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xb7: // LD A,RES 6,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0xbf
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xb8: // LD B,RES 7,(REGISTER+dd)
		z80.B = z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xb9: // LD C,RES 7,(REGISTER+dd)
		z80.C = z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xba: // LD D,RES 7,(REGISTER+dd)
		z80.D = z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xbb: // LD E,RES 7,(REGISTER+dd)
		z80.E = z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xbc: // LD H,RES 7,(REGISTER+dd)
		z80.H = z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xbd: // LD L,RES 7,(REGISTER+dd)
		z80.L = z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xbe: // RES 7,(REGISTER+dd)
		tmp := z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xbf: // LD A,RES 7,(REGISTER+dd)
		z80.A = z80.readByte(memptr) & 0x7f
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xc0: // LD B,SET 0,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xc1: // LD C,SET 0,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xc2: // LD D,SET 0,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xc3: // LD E,SET 0,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xc4: // LD H,SET 0,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xc5: // LD L,SET 0,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xc6: // SET 0,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xc7: // LD A,SET 0,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x01
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xc8: // LD B,SET 1,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xc9: // LD C,SET 1,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xca: // LD D,SET 1,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xcb: // LD E,SET 1,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xcc: // LD H,SET 1,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xcd: // LD L,SET 1,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xce: // SET 1,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xcf: // LD A,SET 1,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x02
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xd0: // LD B,SET 2,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xd1: // LD C,SET 2,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xd2: // LD D,SET 2,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xd3: // LD E,SET 2,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xd4: // LD H,SET 2,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xd5: // LD L,SET 2,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xd6: // SET 2,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xd7: // LD A,SET 2,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x04
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xd8: // LD B,SET 3,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xd9: // LD C,SET 3,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xda: // LD D,SET 3,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xdb: // LD E,SET 3,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xdc: // LD H,SET 3,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xdd: // LD L,SET 3,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xde: // SET 3,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xdf: // LD A,SET 3,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x08
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xe0: // LD B,SET 4,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xe1: // LD C,SET 4,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xe2: // LD D,SET 4,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xe3: // LD E,SET 4,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xe4: // LD H,SET 4,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xe5: // LD L,SET 4,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xe6: // SET 4,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xe7: // LD A,SET 4,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x10
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xe8: // LD B,SET 5,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xe9: // LD C,SET 5,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xea: // LD D,SET 5,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xeb: // LD E,SET 5,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xec: // LD H,SET 5,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xed: // LD L,SET 5,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xee: // SET 5,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xef: // LD A,SET 5,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x20
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xf0: // LD B,SET 6,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xf1: // LD C,SET 6,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xf2: // LD D,SET 6,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xf3: // LD E,SET 6,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xf4: // LD H,SET 6,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xf5: // LD L,SET 6,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xf6: // SET 6,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xf7: // LD A,SET 6,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x40
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	case 0xf8: // LD B,SET 7,(REGISTER+dd)
		z80.B = z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.B)

	case 0xf9: // LD C,SET 7,(REGISTER+dd)
		z80.C = z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.C)

	case 0xfa: // LD D,SET 7,(REGISTER+dd)
		z80.D = z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.D)

	case 0xfb: // LD E,SET 7,(REGISTER+dd)
		z80.E = z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.E)

	case 0xfc: // LD H,SET 7,(REGISTER+dd)
		z80.H = z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.H)

	case 0xfd: // LD L,SET 7,(REGISTER+dd)
		z80.L = z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.L)

	case 0xfe: // SET 7,(REGISTER+dd)
		tmp := z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, tmp)

	case 0xff: // LD A,SET 7,(REGISTER+dd)
		z80.A = z80.readByte(memptr) | 0x80
		z80.readNoReq(memptr, 1)
		z80.writeByte(memptr, z80.A)

	}
}
