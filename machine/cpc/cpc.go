// Package cpc implements the Amstrad CPC machine
package cpc

import (
	"log"

	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/cpu"
	"github.com/jtruco/emu8/device/cpu/z80"
	"github.com/jtruco/emu8/device/memory"
	"github.com/jtruco/emu8/device/tape"
	"github.com/jtruco/emu8/device/video"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/machine"
	"github.com/jtruco/emu8/machine/cpc/format"
	zxformat "github.com/jtruco/emu8/machine/spectrum/format"
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

// Amstrad CPC formats
const (
	cpcFormatSNA = "sna"
	cpcFormatCDT = "cdt"
)

var (
	cpcSnapFormats = []string{cpcFormatSNA}
	cpcTapeFormats = []string{cpcFormatCDT}
)

// AmstradCPC the Amstrad CPC 464
type AmstradCPC struct {
	config     machine.Config        // Machine information
	controller controller.Controller // The emulator controller
	components *device.Components    // Machine device components
	clock      *device.ClockDevice   // The system clock
	cpu        *z80.Z80              // The Zilog Z80A CPU
	memory     *memory.Memory        // The machine memory
	lowerRom   *memory.BankMap       // The lower rom
	upperRom   *memory.BankMap       // The upper rom
	gatearray  *GateArray            // The Gate-Array
	crtc       *video.MC6845         // The Cathode Ray Tube Controller
	psg        *Psg                  // The Programmable Sound Generator
	ppi        *Ppi                  // The Parallel Peripheral Interface
	video      *VduVideo             // The VDU video
	keyboard   *Keyboard             // The matrix keyboard
	tape       *tape.Drive           // The spectrum Tape drive
}

// NewAmstradCPC returns a new Amstrad CPC
func NewAmstradCPC(model int) *AmstradCPC {
	cpc := new(AmstradCPC)
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
	cpc.clock = device.NewClock()
	cpc.cpu = z80.New(cpc.clock, cpc.memory, cpc)
	cpc.crtc = video.NewMC6845()
	cpc.gatearray = NewGateArray(cpc)
	cpc.video = NewVduVideo(cpc)
	cpc.keyboard = NewKeyboard()
	cpc.psg = NewPsg(cpc)
	cpc.ppi = NewPpi(cpc)
	cpc.tape = tape.New(cpc.clock)
	// register all components
	cpc.components = device.NewComponents(10)
	cpc.components.Add(cpc.clock)
	cpc.components.Add(cpc.cpu)
	cpc.components.Add(cpc.memory)
	cpc.components.Add(cpc.gatearray)
	cpc.components.Add(cpc.crtc)
	cpc.components.Add(cpc.video)
	cpc.components.Add(cpc.keyboard)
	cpc.components.Add(cpc.psg)
	cpc.components.Add(cpc.tape)
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
func (cpc *AmstradCPC) Clock() device.Clock {
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
func (cpc *AmstradCPC) SetController(control controller.Controller) {
	control.Video().SetVideo(cpc.video)
	control.Keyboard().AddReceiver(cpc.keyboard, cpcKeyboardMap)
	control.File().RegisterFormat(controller.FormatSnap, cpcSnapFormats)
	control.File().RegisterFormat(controller.FormatTape, cpcTapeFormats)
	control.Tape().SetDrive(cpc.tape)
	cpc.controller = control
}

// Emulation control
// -----------------------------------------------------------------------------

// BeginFrame begin emulation frame tasks
func (cpc *AmstradCPC) BeginFrame() {
	// nothing todo
}

// Emulate one machine step
func (cpc *AmstradCPC) Emulate() {
	// Tape emulation
	if cpc.tape.IsPlaying() {
		cpc.tape.Playback()
	}

	// z80 cpu emulation
	lapse := cpc.cpu.Execute()
	fix := lapse & 0x03 // CPC 4T instruction round
	if fix != 0 {
		fix = 0x04 - lapse
		cpc.clock.Add(fix)
		lapse += fix
	}

	// CPC bus emulation
	cpc.gatearray.Emulate(lapse)
}

// EndFrame end emulation frame tasks
func (cpc *AmstradCPC) EndFrame() {
	cpc.video.EndFrame()
}

// CPC IO bus
// -----------------------------------------------------------------------------

// Read bus at address
func (cpc *AmstradCPC) Read(address uint16) byte {
	var result byte = 0xff
	if address&0x4000 == 0 { // CRTC
		port := byte(address>>8) & 0x03
		result &= cpc.crtc.Read(port)
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
	if address&0x0800 == 0 { // PPI select
		port := byte(address>>8) & 0x3
		cpc.ppi.Write(port, data)
	}
}

// Files : load & save state / tape
// -----------------------------------------------------------------------------

// LoadFile loads a file into machine
func (cpc *AmstradCPC) LoadFile(filename string) {
	filefmt, ext := cpc.controller.File().FileFormat(filename)
	if filefmt == controller.FormatUnknown {
		log.Println("CPC : Not supported format:", ext)
		return
	}
	name := cpc.controller.File().BaseName(filename)
	data, err := cpc.controller.File().LoadFileFormat(filename, filefmt)
	if err != nil {
		log.Println("CPC : Error loading file:", name)
		return
	}
	// load snapshop formats
	if filefmt == controller.FormatSnap {
		var snap *format.Snapshot
		switch ext {
		case cpcFormatSNA:
			snap = format.LoadSNA(data)
		default:
			log.Println("Spectrum : Not implemented format:", ext)
		}
		if snap != nil {
			cpc.LoadState(snap)
		}
	} else if filefmt == controller.FormatTape {
		var tape tape.Tape
		loaded := false
		switch ext {
		case cpcFormatCDT:
			tape = zxformat.NewTzx()
			loaded = tape.Load(data)
		default:
			log.Println("CPC : Not implemented format:", ext)
		}
		if loaded {
			tape.Info().Name = name
			cpc.tape.Insert(tape)
		}
	}
}

// LoadState loads the Amstrad CPC snapshot
func (cpc *AmstradCPC) LoadState(snap *format.Snapshot) {
	// CPU
	cpc.cpu.State.Copy(&snap.State)
	// Memory
	cpc.memory.LoadRAM(0x00, snap.Memory[0:])
	// GateArray
	cpc.gatearray.SetPen(snap.GaSelectedPen)
	for i := 0; i < gaTotalPens; i++ {
		cpc.gatearray.Palette()[i] = int(snap.GaPenColours[i])
	}
	cpc.gatearray.Write(snap.GaMultiConfig)
	// Crtc
	cpc.crtc.SelectRegister(snap.CrtcSelected)
	for i := byte(0); i < 18; i++ {
		cpc.crtc.WriteRegister(i, snap.CrtcRegisters[i])
	}
	// Psg
	cpc.psg.selected = snap.PsgSelected
	for i := byte(0); i < 16; i++ {
		cpc.psg.registers[i] = snap.PsgRegisters[i]
	}
}
