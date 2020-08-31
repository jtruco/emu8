// Package spectrum implements de ZX Spectrum machine
package spectrum

import (
	"log"

	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/audio"
	"github.com/jtruco/emu8/device/cpu"
	"github.com/jtruco/emu8/device/cpu/z80"
	"github.com/jtruco/emu8/device/memory"
	"github.com/jtruco/emu8/device/tape"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/machine"
	"github.com/jtruco/emu8/machine/spectrum/format"
)

// -----------------------------------------------------------------------------
// ZX Spectrum
// -----------------------------------------------------------------------------

// Default ZX Spectrum constants
const (
	zxFPS        = 50    // 50 Hz (50.08 Hz)
	zxTStates    = 69888 // TStates per frame
	zxIntTstates = 32    // ZX Spectrum 16k & 48k
	zxRomName    = "zxspectrum.rom"
)

// ZX Spectrum formats
const (
	formatSNA = "sna"
	formatZ80 = "z80"
	formatTAP = "tap"
	formatTZX = "tzx"
)

var (
	snapFormats = []string{formatSNA, formatZ80}
	tapeFormats = []string{formatTAP, formatTZX}
)

// Spectrum the ZX Spectrum
type Spectrum struct {
	config     machine.Config        // Machine information
	controller controller.Controller // The emulator controller
	components *device.Components    // Machine device components
	clock      *device.ClockDevice   // The system clock
	cpu        *z80.Z80              // The Zilog Z80A CPU
	memory     *memory.Memory        // The machine memory
	ula        *ULA                  // The spectrum ULA
	tv         *TVVideo              // The spectrum TV video output
	beeper     *audio.Beeper         // The spectrum Beeper
	keyboard   *Keyboard             // The spectrum Keyboard
	tape       *tape.Drive           // The spectrum Tape drive
	joystick   *Joystick             // The spectrum Joystick
}

// New returns a new ZX Spectrum
func New(model int) machine.Machine {
	spectrum := new(Spectrum)
	spectrum.config.Model = model
	spectrum.config.FrameTStates = zxTStates
	spectrum.config.SetFPS(zxFPS)
	// memory mapping
	if spectrum.config.Model == machine.ZXSpectrum16k {
		spectrum.memory = memory.New(memory.Size32K, 2)
		spectrum.memory.SetMap(0, memory.NewROM(0x0000, memory.Size16K))
		spectrum.memory.SetMap(1, memory.NewRAM(0x4000, memory.Size16K))
	} else {
		spectrum.memory = memory.New(memory.Size64K, 4)
		spectrum.memory.SetMap(0, memory.NewROM(0x0000, memory.Size16K))
		spectrum.memory.SetMap(1, memory.NewRAM(0x4000, memory.Size16K))
		spectrum.memory.SetMap(2, memory.NewRAM(0x8000, memory.Size16K))
		spectrum.memory.SetMap(3, memory.NewRAM(0xC000, memory.Size16K))
	}
	mapper := &memory.MaskMapper{Shift: 14, Mask: 0x3fff}
	spectrum.memory.SetMapper(mapper)
	// build device components
	spectrum.clock = device.NewClock()
	spectrum.ula = NewULA(spectrum)
	spectrum.cpu = z80.New(spectrum.clock, spectrum.memory, spectrum.ula)
	spectrum.cpu.OnIntAck = spectrum.onInterruptAck
	spectrum.tv = NewTVVideo(spectrum)
	spectrum.beeper = audio.NewBeeper(audio.NewConfig(zxFPS, zxTStates))
	spectrum.beeper.SetMap(beeperMap)
	spectrum.keyboard = NewKeyboard()
	spectrum.tape = tape.New(spectrum.clock)
	spectrum.joystick = NewJoystick()
	// register all components
	spectrum.components = device.NewComponents(9)
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

// Register spectrum machine
func Register() { machine.Register(machine.ZXSpectrum, New) }

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
	data, err := spectrum.controller.File().LoadROM(zxRomName)
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

// SetController connect controllers & components
func (spectrum *Spectrum) SetController(control controller.Controller) {
	control.Video().SetVideo(spectrum.tv)
	control.Audio().SetAudio(spectrum.beeper)
	control.Keyboard().AddReceiver(spectrum.keyboard, zxKeyboardMap)
	control.File().RegisterFormat(controller.FormatSnap, snapFormats)
	control.File().RegisterFormat(controller.FormatTape, tapeFormats)
	control.Tape().SetDrive(spectrum.tape)
	control.Joystick().AddReceiver(spectrum.joystick, 0)
	spectrum.controller = control
}

// VideoMemory gets the video memory bank
func (spectrum *Spectrum) VideoMemory() *memory.Bank {
	return spectrum.memory.Bank(1)
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
	if spectrum.tape.IsPlaying() {
		spectrum.tape.Playback()
	}
	// Exetues a CPU instruction
	spectrum.cpu.Execute()
	// Maskable interrupt request lenght
	if spectrum.cpu.IntRq && spectrum.clock.Tstates() >= zxIntTstates {
		spectrum.cpu.InterruptRequest(false)
	}
}

// EndFrame end emulation frame tasks
func (spectrum *Spectrum) EndFrame() {
	// nothing to do
}

// onInterruptAck
func (spectrum *Spectrum) onInterruptAck() bool {
	return true
}

// Snapshots : load & save state

// LoadFile loads a file into machine
func (spectrum *Spectrum) LoadFile(filename string) {
	info := spectrum.controller.File().FileInfo(filename)
	if info.Format == controller.FormatUnknown {
		log.Println("Spectrum : Not supported format:", info.Ext)
		return
	}
	err := spectrum.controller.File().LoadFile(info)
	if err != nil {
		log.Println("Spectrum : Error loading file:", info.Name)
		return
	}
	// load snapshop formats
	if info.Format == controller.FormatSnap {
		var snap *format.Snapshot
		switch info.Ext {
		case formatSNA:
			snap = format.LoadSNA(info.Data)
		case formatZ80:
			snap = format.LoadZ80(info.Data)
		default:
			log.Println("Spectrum : Not implemented snap format:", info.Ext)
		}
		if snap != nil {
			spectrum.LoadState(snap)
		}
	} else if info.Format == controller.FormatTape {
		var tape tape.Tape
		loaded := false
		switch info.Ext {
		case formatTAP:
			tape = format.NewTap()
			loaded = tape.Load(info.Data)
		case formatTZX:
			tape = format.NewTzx()
			loaded = tape.Load(info.Data)
		default:
			log.Println("Spectrum : Not implemented tape format:", info.Ext)
		}
		if loaded {
			tape.Info().Name = info.Name
			spectrum.tape.Insert(tape)
		}
	}
}

// TakeSnapshot takes and saves snapshop of the machine state
func (spectrum *Spectrum) TakeSnapshot() {
	snap := spectrum.SaveState()
	data := snap.SaveSNA()
	name := spectrum.controller.File().NewName("speccy", formatSNA)
	err := spectrum.controller.File().SaveFile(name, controller.FormatSnap, data)
	if err == nil {
		log.Println("Spectrum : Snapshot saved: ", name)
	} else {
		log.Println("Spectrum : Error saving snapshot: ", name)
	}
}

// LoadState loads a ZX Spectrum snapshot
func (spectrum *Spectrum) LoadState(snap *format.Snapshot) {
	// CPU
	spectrum.cpu.State.Copy(&snap.State)
	// TStates
	spectrum.clock.SetTstates(snap.Tstates)
	// Border
	spectrum.tv.SetBorder(snap.Border)
	// Memory
	spectrum.memory.LoadRAM(0x4000, snap.Memory[0:0xC000])
}

// SaveState save ZX Spectrum state
func (spectrum *Spectrum) SaveState() *format.Snapshot {
	var snap = new(format.Snapshot)
	// CPU
	snap.State.Copy(&spectrum.cpu.State)
	// Clock
	snap.Tstates = spectrum.clock.Tstates()
	// Border
	snap.Border = spectrum.tv.border
	// Memory banks (16k, 48k)
	copy(snap.Memory[0x0000:], spectrum.memory.Bank(1).Data())
	if spectrum.config.Model == machine.ZXSpectrum48k {
		copy(snap.Memory[0x4000:], spectrum.memory.Bank(2).Data())
		copy(snap.Memory[0x8000:], spectrum.memory.Bank(3).Data())
	}
	return snap
}
