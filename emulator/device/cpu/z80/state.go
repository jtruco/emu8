package z80

import "github.com/jtruco/emu8/emulator/device/cpu"

// -----------------------------------------------------------------------------
// State - Z80 CPU State
// -----------------------------------------------------------------------------

// State is the Z80 cpu state
type State struct {

	// 8 bit registers, accumulator, flags, BC, DE, HL, Indexes and second bank
	A, F, B, C, D, E, H, L         cpu.Reg8
	Ax, Fx, Bx, Cx, Dx, Ex, Hx, Lx cpu.Reg8
	IXh, IXl, IYh, IYl             cpu.Reg8

	// Interrupt register, refresh and internal WZ register
	I    cpu.Reg8
	R    cpu.Reg8
	W, Z cpu.Reg8

	// Program and Stack 16bit registers
	SP, PC cpu.Reg16

	// 16 bit register pairs
	AF, BC, DE, HL     cpu.Register16
	AFx, BCx, DEx, HLx cpu.Register16
	IX, IY             cpu.Register16
	IR                 cpu.Register16
	Memptr             cpu.Register16

	// Control
	Halted     bool
	IM         byte
	IFF1, IFF2 bool
	ActiveEI   bool
	ReadIFF2   bool
	IntRq      bool
	NmiRq      bool
}

// NewState creates a new Z80 state
func NewState() *State {
	state := new(State)
	state.Init()
	return state
}

// Init initializes state
func (state *State) Init() {
	state.AF = cpu.NewRegister16(&state.A, &state.F)
	state.AFx = cpu.NewRegister16(&state.Ax, &state.Fx)
	state.BC = cpu.NewRegister16(&state.B, &state.C)
	state.BCx = cpu.NewRegister16(&state.Bx, &state.Cx)
	state.DE = cpu.NewRegister16(&state.D, &state.E)
	state.DEx = cpu.NewRegister16(&state.Dx, &state.Ex)
	state.HL = cpu.NewRegister16(&state.H, &state.L)
	state.HLx = cpu.NewRegister16(&state.Hx, &state.Lx)
	state.IX = cpu.NewRegister16(&state.IXh, &state.IXl)
	state.IY = cpu.NewRegister16(&state.IYh, &state.IYl)
	state.IR = cpu.NewRegister16(&state.I, &state.R)
	state.Memptr = cpu.NewRegister16(&state.W, &state.Z)
	state.initControl()
}

// SoftReset initializes state (soft)
func (state *State) SoftReset() {
	state.PC = 0
	state.SP = 0xffff
	state.AF.Set(0xffff)
	state.AFx.Set(0xffff)
	state.IR.Set(0x0)
	state.initControl()
}

// HardReset initializes state (hard)
func (state *State) HardReset() {
	state.SoftReset()
	state.BC.Set(0x0)
	state.BCx.Set(0x0)
	state.DE.Set(0x0)
	state.DEx.Set(0x0)
	state.HL.Set(0x0)
	state.HLx.Set(0x0)
	state.IX.Set(0x0)
	state.IY.Set(0x0)
	state.Memptr.Set(0x0)
}

// initControl initializes control variables
func (state *State) initControl() {
	state.Halted = false
	state.IM = 0
	state.IFF1 = false
	state.IFF2 = false
	state.ActiveEI = false
	state.ReadIFF2 = false
	state.IntRq = false
	state.NmiRq = false
}

// Copy copies state values
func (state *State) Copy(value *State) {
	state.A = value.A
	state.F = value.F
	state.B = value.B
	state.C = value.C
	state.D = value.D
	state.E = value.E
	state.H = value.H
	state.L = value.L
	state.Ax = value.Ax
	state.Fx = value.Fx
	state.Bx = value.Bx
	state.Cx = value.Cx
	state.Dx = value.Dx
	state.Ex = value.Ex
	state.Hx = value.Hx
	state.Lx = value.Lx
	state.IXl = value.IXl
	state.IXh = value.IXh
	state.IYl = value.IYl
	state.IYh = value.IYh
	state.I = value.I
	state.R = value.R
	state.W = value.W
	state.Z = value.Z
	state.SP = value.SP
	state.PC = value.PC
	state.Halted = value.Halted
	state.IM = value.IM
	state.IFF1 = value.IFF1
	state.IFF2 = value.IFF2
	state.ActiveEI = value.ActiveEI
	state.ReadIFF2 = value.ActiveEI
}
