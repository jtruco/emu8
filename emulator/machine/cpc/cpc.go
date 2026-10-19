// Package cpc implements the Amstrad CPC machine
package cpc

import (
	"fmt"

	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/jtruco/emu8/emulator/device/bus"
	"github.com/jtruco/emu8/emulator/device/cpu"
	"github.com/jtruco/emu8/emulator/device/cpu/z80"
	"github.com/jtruco/emu8/emulator/device/io/tape"
	"github.com/jtruco/emu8/emulator/device/memory"
	"github.com/jtruco/emu8/emulator/device/video"
	"github.com/jtruco/emu8/emulator/machine"
	"github.com/jtruco/emu8/emulator/machine/cpc/format"
)

// -----------------------------------------------------------------------------
// Amstrad CPC
// -----------------------------------------------------------------------------

// Amstrad CPC models
const (
	AmstradCPC464 = iota
)

// Default Amstrad CPC
const (
	cpcFPS          = 50              // 50 Hz ( 50.08 Hz )
	cpcTStates      = 79872           // TStates per frame ( 312 sl * 256 Ts ) ~ 4 Mhz
	cpcAudioTStates = cpcTStates >> 5 // Audio TStates (~ 1MHz / 8)
	cpcOsRomName    = "cpc464_os.rom"
	cpcOsRomNameES  = "cpc464_os_es.rom"
	cpcOsRomNameFR  = "cpc464_os_fr.rom"
	cpcBasicRomName = "cpc464_basic.rom"
	cpcJumpers      = 0x1e
	cpcLowerROM     = 0
	cpcUpperROM     = 4
)

// AmstradCPC the Amstrad CPC 464
type AmstradCPC struct {
	config     machine.Config      // Machine information
	control    machine.Control     // The emulator controller
	components *device.Components  // Machine device components
	clock      *device.ClockDevice // The system clock
	cpu        *z80.Z80            // The Zilog Z80A CPU
	memory     *memory.Memory      // The machine memory
	gatearray  *GateArray          // The Gate-Array
	crtc       *video.MC6845       // The Cathode Ray Tube Controller
	psg        *audio.AY38910      // The Programmable Sound Generator
	ppi        *Ppi                // The Parallel Peripheral Interface
	video      *VduVideo           // The VDU video
	keyboard   *Keyboard           // The matrix keyboard
	tape       *tape.Drive         // The tape drive
	joystick   *Joystick           // The CPC Joystick
}

// New returns a new Amstrad CPC
func New(model int, settings machine.Settings) machine.Machine {
	cpc := new(AmstradCPC)
	cpc.config.Settings = settings
	cpc.config.Model = model
	cpc.config.SetTimings(cpcTStates, cpcFPS)
	// memory map
	cpc.memory = memory.New(6)
	cpc.memory.SetMap(0, memory.NewROM(0x0000, memory.Size16K)) // Lower ROM Bios
	cpc.memory.SetMap(1, memory.NewRAM(0x0000, memory.Size16K))
	cpc.memory.SetMap(2, memory.NewRAM(0x4000, memory.Size16K))
	cpc.memory.SetMap(3, memory.NewRAM(0x8000, memory.Size16K))
	cpc.memory.SetMap(4, memory.NewROM(0xC000, memory.Size16K)) // Upper ROM Basic
	cpc.memory.SetMap(5, memory.NewRAM(0xC000, memory.Size16K))
	cpc.memory.SetMapper(bus.NewPageMapper(14)) // ROMs over RAM : 16K pages
	// devices
	cpc.clock = device.NewClock()
	cpc.cpu = z80.New(cpc.clock, cpc.memory, cpc)
	cpc.crtc = video.NewMC6845()
	cpc.gatearray = NewGateArray(cpc)
	cpc.video = NewVduVideo(cpc)
	cpc.keyboard = NewKeyboard()
	cpc.psg = audio.NewAY38910(
		audio.NewConfig(cpc.config.Frequency(), cpcFPS, cpcAudioTStates))
	cpc.psg.OnReadPortA = cpc.onPsgReadPortA
	cpc.ppi = NewPpi(cpc)
	cpc.tape = tape.New(cpc.clock)
	cpc.joystick = NewJoystick(cpc.keyboard)
	// register all components
	cpc.components = device.NewComponents()
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
	cpc.components.Add(cpc.joystick)
	return cpc
}

// Device interface
// -----------------------------------------------------------------------------

// Init initializes machine and loads its ROMs
func (cpc *AmstradCPC) Init() error {
	cpc.components.Init()
	osrom := cpcOsRomName
	switch cpc.config.Options { // OS ROM language
	case "es":
		osrom = cpcOsRomNameES
	case "fr":
		osrom = cpcOsRomNameFR
	}
	if err := machine.LoadROM(cpc.control, osrom, cpc.memory.Bank(cpcLowerROM)); err != nil {
		return err
	}
	return machine.LoadROM(cpc.control, cpcBasicRomName, cpc.memory.Bank(cpcUpperROM))
}

// Reset resets the machine
func (cpc *AmstradCPC) Reset() {
	cpc.components.Reset()
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

// InitControl connect UI controllers & device components
func (cpc *AmstradCPC) InitControl(control machine.Control) {
	// Bind devices
	control.BindVideo(cpc.video)
	control.BindAudio(cpc.psg)
	control.BindKeyboard(cpc.keyboard)
	control.BindJoystick(cpc.joystick)
	control.BindTapeDrive(cpc.tape)
	// Register formats
	control.RegisterSnapshot(format.SNA)
	control.RegisterTape(format.CDT, format.NewCdt)
	cpc.control = control
}

// Emulation control
// -----------------------------------------------------------------------------

// BeginFrame begin emulation frame tasks
func (cpc *AmstradCPC) BeginFrame() {
	cpc.psg.BeginFrame()
}

// Emulate one machine step
func (cpc *AmstradCPC) Emulate() {
	// Executes a CPU instruction
	tstates := cpc.cpu.Execute()

	// 4-tstate rounding fix
	fix := (0x04 - tstates) & 0x03
	if fix != 0 {
		cpc.clock.Add(fix)
		tstates += fix
	}

	// CPC bus emulation
	cpc.gatearray.Emulate(tstates)

	// Tape emulation
	cpc.tape.Emulate(tstates)
}

// EndFrame end emulation frame tasks
func (cpc *AmstradCPC) EndFrame() {
	cpc.video.EndFrame()
	cpc.psg.EndFrame()
}

// CPC IO bus
// -----------------------------------------------------------------------------

// Read bus at address. The Z80 core adds the 4 Ts of the I/O cycle; the Gate
// Array wait states are applied per instruction in Emulate (known D7, B17 for
// OUT (C),r and IN r,(C)).
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

// Write bus at address. The Z80 core adds the 4 Ts of the I/O cycle; the Gate
// Array wait states are applied per instruction in Emulate (known D7, B17 for
// OUT (C),r and IN r,(C)).
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

// onPsgReadPortA
func (cpc *AmstradCPC) onPsgReadPortA() byte {
	// Keyboard connected to PSG Port A
	return cpc.keyboard.State()
}

// Files : load & save state / tape
// -----------------------------------------------------------------------------

// LoadState loads an Amstrad CPC snapshot
func (cpc *AmstradCPC) LoadState(state machine.State) error {
	var snap *format.Snapshot
	switch state.Format {
	case format.SNA:
		snap = format.LoadSNA(state.Data)
	default:
		return fmt.Errorf("not supported snapshot format %q", state.Format)
	}
	if snap == nil {
		return fmt.Errorf("invalid %s snapshot", state.Format)
	}
	cpc.loadSnapshot(snap)
	return nil
}

func (cpc *AmstradCPC) loadSnapshot(snap *format.Snapshot) {
	// CPU
	cpc.cpu.State.Copy(&snap.State)
	// Memory
	cpc.memory.LoadRAM(0x00, snap.Memory[0:])
	// GateArray
	cpc.gatearray.SetPen(snap.GaSelectedPen)
	palette := cpc.gatearray.Palette()
	for i := 0; i < gaTotalPens; i++ {
		palette[i] = int(snap.GaPenColours[i])
	}
	cpc.gatearray.Write(snap.GaMultiConfig)
	// Crtc
	cpc.crtc.SelectRegister(snap.CrtcSelected)
	for i := byte(0); i < 18; i++ {
		cpc.crtc.WriteRegister(i, snap.CrtcRegisters[i])
	}
	// Psg
	cpc.psg.SelectRegister(snap.PsgSelected)
	for i := byte(0); i < 16; i++ {
		cpc.psg.WriteRegister(i, snap.PsgRegisters[i])
	}
}

// SaveState saves an Amstrad CPC snapshot
func (cpc *AmstradCPC) SaveState() machine.State {
	return machine.State{
		Format: format.SNA,
		Data:   cpc.saveSnapshot().SaveSNA()}
}

func (cpc *AmstradCPC) saveSnapshot() *format.Snapshot {
	var snap = new(format.Snapshot)
	// CPU
	snap.State.Copy(&cpc.cpu.State)
	// Memory banks (64k)
	cpc.memory.Bank(1).Save(snap.Memory[0x0000:])
	cpc.memory.Bank(2).Save(snap.Memory[0x4000:])
	cpc.memory.Bank(3).Save(snap.Memory[0x8000:])
	cpc.memory.Bank(5).Save(snap.Memory[0xC000:])
	// GateArray
	snap.GaSelectedPen = cpc.gatearray.Pen()
	palette := cpc.gatearray.Palette()
	for i := 0; i < gaTotalPens; i++ {
		snap.GaPenColours[i] = byte(palette[i])
	}
	snap.GaMultiConfig = cpc.gatearray.Config()
	// Crtc
	snap.CrtcSelected = cpc.crtc.Selected()
	for i := byte(0); i < 18; i++ {
		snap.CrtcRegisters[i] = cpc.crtc.Register(i)
	}
	// Psg
	snap.PsgSelected = cpc.psg.Selected()
	for i := byte(0); i < 16; i++ {
		snap.PsgRegisters[i] = cpc.psg.Register(i)
	}
	return snap
}
