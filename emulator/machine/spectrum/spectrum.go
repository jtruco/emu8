// Package spectrum implements de ZX Spectrum machine
package spectrum

import (
	"log"

	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/audio"
	"github.com/jtruco/emu8/emulator/device/bus"
	"github.com/jtruco/emu8/emulator/device/cpu"
	"github.com/jtruco/emu8/emulator/device/cpu/z80"
	"github.com/jtruco/emu8/emulator/device/io/tape"
	"github.com/jtruco/emu8/emulator/device/memory"
	"github.com/jtruco/emu8/emulator/machine"
	"github.com/jtruco/emu8/emulator/machine/spectrum/format"
)

// -----------------------------------------------------------------------------
// ZX Spectrum
// -----------------------------------------------------------------------------

// ZX Spectrum model constants
const (
	ZXSpectrum16K = iota
	ZXSpectrum48K
)

// Default ZX Spectrum constants
const (
	zxFPS         = 50    // 50 Hz (50.08 Hz)
	zxTStates     = 69888 // TStates per frame
	zxIntTstates  = 32    // ZX Spectrum 16k & 48k
	zxVideoMemory = 1     // Video memory bank
	zxRomName     = "zxspectrum.rom"
)

// Spectrum the ZX Spectrum
type Spectrum struct {
	config     machine.Config      // Machine information
	control    machine.Control     // The emulator controller
	components *device.Components  // Machine device components
	clock      *device.ClockDevice // The system clock
	cpu        *z80.Z80            // The Zilog Z80A CPU
	memory     *memory.Memory      // The machine memory
	ula        *ULA                // The spectrum ULA
	tv         *TvVideo            // The spectrum TV video output
	beeper     *audio.Beeper       // The spectrum Beeper
	keyboard   *Keyboard           // The spectrum Keyboard
	tape       *tape.Drive         // The spectrum Tape drive
	joystick   *Joystick           // The spectrum Joystick
}

// New returns a new ZX Spectrum
func New(model int) machine.Machine {
	spectrum := new(Spectrum)
	spectrum.config.Model = model
	spectrum.config.SetTimings(zxTStates, zxFPS)
	// memory mapping
	if spectrum.config.Model == ZXSpectrum16K {
		spectrum.memory = memory.New(2)
		spectrum.memory.SetMap(0, memory.NewROM(0x0000, memory.Size16K))
		spectrum.memory.SetMap(1, memory.NewRAM(0x4000, memory.Size16K))
	} else {
		spectrum.memory = memory.New(4)
		spectrum.memory.SetMap(0, memory.NewROM(0x0000, memory.Size16K))
		spectrum.memory.SetMap(1, memory.NewRAM(0x4000, memory.Size16K))
		spectrum.memory.SetMap(2, memory.NewRAM(0x8000, memory.Size16K))
		spectrum.memory.SetMap(3, memory.NewRAM(0xC000, memory.Size16K))
	}
	spectrum.memory.SetMapper(bus.NewMaskMapper(14))
	// build device components
	spectrum.clock = device.NewClock()
	spectrum.ula = NewULA(spectrum)
	spectrum.cpu = z80.New(spectrum.clock, spectrum.memory, spectrum.ula)
	spectrum.cpu.OnIntAck = spectrum.onInterruptAck
	spectrum.tv = NewTVVideo(spectrum)
	spectrum.beeper = audio.NewBeeper(audio.NewConfig(zxFPS, zxTStates))
	spectrum.beeper.SetMap(zxBeeperMap)
	spectrum.keyboard = NewKeyboard()
	spectrum.tape = tape.New(spectrum.clock)
	spectrum.joystick = NewJoystick()
	// register all components
	spectrum.components = device.NewComponents()
	spectrum.components.Add(spectrum.clock)
	spectrum.components.Add(spectrum.memory)
	spectrum.components.Add(spectrum.ula)
	spectrum.components.Add(spectrum.cpu)
	spectrum.components.Add(spectrum.tv)
	spectrum.components.Add(spectrum.beeper)
	spectrum.components.Add(spectrum.keyboard)
	spectrum.components.Add(spectrum.tape)
	spectrum.components.Add(spectrum.joystick)

	return spectrum
}

// Device interface

// Init initializes machine
func (spectrum *Spectrum) Init() {
	// initialize components
	spectrum.components.Init()
	// initialize spectrum
	spectrum.initSpectrum()
}

// Reset resets the machine
func (spectrum *Spectrum) Reset() {
	// reset components
	spectrum.components.Reset()
	// resets spectrum
	spectrum.initSpectrum()
}

// initSpectrum commont init tasks
func (spectrum *Spectrum) initSpectrum() {
	// load ROM at bank 0
	data, err := spectrum.control.LoadROM(zxRomName)
	if err != nil {
		return
	}
	rom := spectrum.memory.Bank(0)
	rom.Load(0, data[0:0x4000])
}

// Machine properties

// Clock gets the machine clock
func (spectrum *Spectrum) Clock() device.Clock {
	return spectrum.clock
}

// Config gets the machine info
func (spectrum *Spectrum) Config() *machine.Config {
	return &spectrum.config
}

// CPU gets the machine CPU
func (spectrum *Spectrum) CPU() cpu.CPU {
	return spectrum.cpu
}

// Components gets the machine components
func (spectrum *Spectrum) Components() *device.Components {
	return spectrum.components
}

// InitControl connect controllers & components
func (spectrum *Spectrum) InitControl(control machine.Control) {
	// Bind devices
	control.BindVideo(spectrum.tv)
	control.BindAudio(spectrum.beeper)
	control.BindKeyboard(spectrum.keyboard)
	control.BindJoystick(spectrum.joystick)
	control.BindTapeDrive(spectrum.tape)
	// Register formats
	control.RegisterSnapshot(format.SNA)
	control.RegisterSnapshot(format.Z80)
	control.RegisterTape(format.TAP, format.NewTap)
	control.RegisterTape(format.TZX, format.NewTzx)
	spectrum.control = control
}

// Emulation control

// BeginFrame begin emulation frame tasks
func (spectrum *Spectrum) BeginFrame() {
	// Request cpu maskable interrupt
	spectrum.cpu.InterruptRequest(true)
}

// Emulate one machine step
func (spectrum *Spectrum) Emulate() {
	// Tape emulation
	spectrum.tape.Playback()

	// Executes a CPU instruction
	spectrum.cpu.Execute()

	// Maskable interrupt request length
	if spectrum.cpu.IntRq && spectrum.clock.Tstates() >= zxIntTstates {
		spectrum.cpu.InterruptRequest(false)
	}
}

// EndFrame end emulation frame tasks
func (spectrum *Spectrum) EndFrame() {} // nothing to do

// onInterruptAck
func (spectrum *Spectrum) onInterruptAck() bool {
	return true
}

// Snapshots : load & save state

// LoadState loads a ZX Spectrum snapshot
func (spectrum *Spectrum) LoadState(state machine.State) {
	var snap *format.Snapshot
	switch state.Format {
	case format.SNA:
		snap = format.LoadSNA(state.Data)
	case format.Z80:
		snap = format.LoadZ80(state.Data)
	default:
		log.Println("Spectrum : Not implemented snap format:", state.Format)
	}
	if snap != nil {
		spectrum.loadSnapshot(snap)
	}
}

func (spectrum *Spectrum) loadSnapshot(snap *format.Snapshot) {
	spectrum.cpu.State.Copy(&snap.State)    // CPU
	spectrum.clock.SetTstates(snap.Tstates) // TStates
	spectrum.tv.SetBorder(snap.Border)      // Border
	// Memory banks (16k, 48k)
	if spectrum.config.Model == ZXSpectrum16K {
		spectrum.memory.LoadRAM(0x4000, snap.Memory[0:0x4000])
	} else {
		spectrum.memory.LoadRAM(0x4000, snap.Memory[0:0xC000])
	}
}

// SaveState loads a ZX Spectrum snapshot
func (spectrum *Spectrum) SaveState() machine.State {
	return machine.State{
		Format: format.SNA,
		Data:   spectrum.saveSnapshot().SaveSNA()}
}

func (spectrum *Spectrum) saveSnapshot() *format.Snapshot {
	var snap = new(format.Snapshot)
	snap.State.Copy(&spectrum.cpu.State)    // CPU
	snap.Tstates = spectrum.clock.Tstates() // Clock
	snap.Border = spectrum.tv.border        // Border
	// Memory banks (16k, 48k)
	spectrum.memory.Bank(1).Save(snap.Memory[0x0000:])
	if spectrum.config.Model == ZXSpectrum48K {
		spectrum.memory.Bank(2).Save(snap.Memory[0x4000:])
		spectrum.memory.Bank(3).Save(snap.Memory[0x8000:])
	}
	return snap
}
