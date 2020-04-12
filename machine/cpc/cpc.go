// Package cpc implements the Amstrad CPC machine
package cpc

import (
	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/cpu"
	"github.com/jtruco/emu8/device/cpu/z80"
	"github.com/jtruco/emu8/device/memory"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/machine"
)

// -----------------------------------------------------------------------------
// Amstrad CPC
// -----------------------------------------------------------------------------

// Default Amstrad CPC
const (
	cpcFPS          = 50    // 50 Hz ( 50.08 Hz )
	cpcFrameTStates = 79872 // TStates per frame ( 312 sl * 256 Ts ) ~ 4 Mhz
	cpcRomName      = "cpc464.rom"
	cpcJumpers      = 0x1e
)

// AmstradCPC the Amstrad CPC 464
type AmstradCPC struct {
	config     machine.Config        // Machine information
	controller controller.Controller // The emulator controller
	components *device.Components    // Machine device components
	clock      *cpu.ClockDevice      // The system clock
	cpu        *z80.Z80              // The Zilog Z80A CPU
	memory     *memory.Memory        // The machine memory
	lowerRom   *memory.BankMap       // The lower rom
	upperRom   *memory.BankMap       // The upper rom
	gatearray  *GateArray            // The Gate-Array
	crtc       *Crtc                 // The Cathode Ray Tube Controller
	psg        *Psg                  // The Programmable Sound Generator
	ppi        *Ppi                  // The Parallel Peripheral Interface
	video      *VduVideo             // The VDU video
	keyboard   *Keyboard             // The matrix keyboard
}

// NewAmstradCPC returns a new Amstrad CPC
func NewAmstradCPC(model int) *AmstradCPC {
	cpc := &AmstradCPC{}
	cpc.config.Model = model
	cpc.config.FrameTStates = cpcFrameTStates
	cpc.config.SetFPS(cpcFPS)
	// memory map
	cpc.memory = memory.New(memory.Size64K, 6)
	cpc.memory.SetMap(0, memory.NewROM(0x0000, memory.Size16K)) // Lower ROM Bios
	cpc.memory.SetMap(1, memory.NewRAM(0x0000, memory.Size16K))
	cpc.memory.SetMap(2, memory.NewRAM(0x4000, memory.Size16K))
	cpc.memory.SetMap(3, memory.NewRAM(0x8000, memory.Size16K))
	cpc.memory.SetMap(4, memory.NewROM(0xC000, memory.Size16K)) // Upper ROM Basic
	cpc.memory.SetMap(5, memory.NewRAM(0xC000, memory.Size16K))
	cpc.lowerRom = cpc.memory.Map(0)
	cpc.upperRom = cpc.memory.Map(4)
	// devices
	cpc.clock = cpu.NewClock()
	cpc.cpu = z80.New(cpc.clock, cpc.memory, cpc)
	cpc.gatearray = NewGateArray(cpc)
	cpc.crtc = NewCrtc(cpc)
	cpc.video = NewVduVideo(cpc)
	cpc.keyboard = NewKeyboard()
	cpc.psg = NewPsg(cpc)
	cpc.ppi = NewPpi(cpc)
	// register all components
	cpc.components = device.NewComponents(9)
	cpc.components.Add(cpc.clock)
	cpc.components.Add(cpc.cpu)
	cpc.components.Add(cpc.memory)
	cpc.components.Add(cpc.gatearray)
	cpc.components.Add(cpc.crtc)
	cpc.components.Add(cpc.video)
	cpc.components.Add(cpc.keyboard)
	cpc.components.Add(cpc.psg)
	cpc.components.Add(cpc.ppi)
	return cpc
}

// Device interface
// -----------------------------------------------------------------------------

// Init initializes machine
func (cpc *AmstradCPC) Init() {
	cpc.components.Init()
	cpc.initAmstrad()
}

// Reset resets the machine
func (cpc *AmstradCPC) Reset() {
	cpc.components.Reset()
	cpc.initAmstrad()
}

// initAmstrad common init tasks
func (cpc *AmstradCPC) initAmstrad() {
	// rom load
	data, err := cpc.controller.File().LoadROM(cpcRomName)
	if err != nil {
		return
	}
	cpc.lowerRom.Bank().Load(0, data[:0x4000]) // lower rom
	cpc.upperRom.Bank().Load(0, data[0x4000:]) // upper rom
	// devices
	cpc.ppi.jumpers = cpcJumpers
}

// Machine interface
// -----------------------------------------------------------------------------

// Clock gets the machine clock
func (cpc *AmstradCPC) Clock() cpu.Clock {
	return cpc.clock
}

// Config gets the machine info
func (cpc *AmstradCPC) Config() *machine.Config {
	return &cpc.config
}

// CPU gets the machine CPU
func (cpc *AmstradCPC) CPU() cpu.CPU {
	return cpc.cpu
}

// Components gets the machine components
func (cpc *AmstradCPC) Components() *device.Components {
	return cpc.components
}

// SetController connect UI controllers & device components
func (cpc *AmstradCPC) SetController(cntrlr controller.Controller) {
	cpc.controller = cntrlr
	cpc.controller.Video().SetVideo(cpc.video)
	cpc.controller.Keyboard().AddReceiver(cpc.keyboard, cpcKeyboardMap)
}

// Emulation control
// -----------------------------------------------------------------------------

// BeginFrame begin emulation frame tasks
func (cpc *AmstradCPC) BeginFrame() {}

// Emulate one machine step
func (cpc *AmstradCPC) Emulate() {
	// z80 cpu emulation
	ts := cpc.clock.Tstates()
	cpc.cpu.Execute()
	lapse := cpc.cpu.Clock().Tstates() - ts

	// CPC 4T instruction halt
	fix := lapse & 0x03
	if fix != 0 {
		fix = 0x04 - lapse
		cpc.clock.Add(fix)
		lapse += fix
	}

	// amstrad bus emulation
	cpc.gatearray.Emulate(lapse)
}

// EndFrame end emulation frame tasks
func (cpc *AmstradCPC) EndFrame() {
	cpc.video.EndFrame()
}

// Files : load & save state / tape

// LoadFile loads a file into machine
func (cpc *AmstradCPC) LoadFile(filename string) {
	// TODO
}

// CPC IO bus
// -----------------------------------------------------------------------------

// Read bus at address
func (cpc *AmstradCPC) Read(address uint16) byte {
	var result byte = 0xff
	if address&0x4000 == 0 { // CRTC
		port := byte(address>>8) & 0x3
		result &= cpc.crtc.Read(port)
	}
	if address&0xC000 == 0x4000 { // Gate-Array
		result &= cpc.gatearray.Read()
	}
	if address&0x0800 == 0 { // PPI select
		port := byte(address>>8) & 0x3
		result &= cpc.ppi.Read(port)
	}
	return result
}

// Write bus at address
func (cpc *AmstradCPC) Write(address uint16, data byte) {
	if address&0x4000 == 0 { // CRTC
		port := byte(address>>8) & 0x3
		cpc.crtc.Write(port, data)
	}
	if address&0xC000 == 0x4000 { // Gate-Array
		cpc.gatearray.Write(data)
	}
	if address&0x2000 == 0 { // ROM select
		// TODO
	}
	if address&0x0800 == 0 { // PPI select
		port := byte(address>>8) & 0x3
		cpc.ppi.Write(port, data)
	}
}
