// Package spectrum implements de ZX Spectrum machine
package spectrum

import (
	"github.com/jtruco/emu8/cpu"
	"github.com/jtruco/emu8/cpu/z80"
	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/audio"
	"github.com/jtruco/emu8/device/memory"
	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/machine"
	"github.com/jtruco/emu8/machine/spectrum/format"
)

// -----------------------------------------------------------------------------
// ZX Spectrum
// -----------------------------------------------------------------------------

// Default ZX Spectrum constants
const (
	fps            = 50    // 50 Hz (50.08 Hz)
	frameTStates   = 69888 // TStates per frame
	audioFrecuency = 48000 // 48 KHz
	romName        = "zxspectrum.rom"
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
	clock      *cpu.ClockDevice      // The system clock
	cpu        *z80.Z80              // The Zilog Z80A CPU
	memory     *memory.Memory        // The machine memory
	ula        *ULA                  // The spectrum ULA
	tv         *TVVideo              // The spectrum TV video output
	beeper     *audio.Beeper         // The spectrum Beeper
	keyboard   *Keyboard             // The spectrum Keyboard
}

// NewSpectrum returns a new ZX Spectrum
func NewSpectrum(model int) *Spectrum {
	spectrum := &Spectrum{}
	spectrum.config.Model = model
	spectrum.config.FrameTStates = frameTStates
	spectrum.config.SetFPS(fps)
	spectrum.buildMachine()
	return spectrum
}

// buildMachine create and connect machine components
func (spectrum *Spectrum) buildMachine() {

	// Build components
	spectrum.clock = cpu.NewClock()
	spectrum.buildMemory()
	spectrum.ula = NewULA(spectrum)
	spectrum.cpu = z80.New(spectrum.clock, spectrum.memory, spectrum.ula)
	spectrum.tv = NewTVVideo(spectrum)
	spectrum.beeper = audio.NewBeeper(audioFrecuency, fps, frameTStates)
	spectrum.beeper.SetMap(beeperMap)
	spectrum.keyboard = NewKeyboard()

	// register components
	spectrum.registerComponents()
}

// buildMemory builds memory mapping
func (spectrum *Spectrum) buildMemory() {
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
	spectrum.memory.SetMapper(&memory.BusMapper{Shift: 14, Mask: 0x3fff})
}

// register components
func (spectrum *Spectrum) registerComponents() {
	spectrum.components = device.NewComponents(7)
	spectrum.components.Add(spectrum.clock)
	spectrum.components.Add(spectrum.memory)
	spectrum.components.Add(spectrum.ula)
	spectrum.components.Add(spectrum.cpu)
	spectrum.components.Add(spectrum.keyboard)
	spectrum.components.Add(spectrum.tv)
	spectrum.components.Add(spectrum.beeper)
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

func (spectrum *Spectrum) initSpectrum() {
	// load ROM at bank 0
	data, err := spectrum.controller.File().LoadROM(romName)
	if err != nil {
		return
	}
	rom := spectrum.memory.GetBankMap(0).Bank()
	rom.Load(0, data[0:0x4000])
}

// Machine properties

// Clock gets the machine clock
func (spectrum *Spectrum) Clock() cpu.Clock {
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
func (spectrum *Spectrum) SetController(cntrlr controller.Controller) {
	spectrum.controller = cntrlr
	spectrum.controller.Video().SetVideo(spectrum.tv)
	spectrum.controller.Audio().SetAudio(spectrum.beeper)
	spectrum.controller.Keyboard().AddReceiver(spectrum.keyboard, zxKeyboardMap)
	spectrum.controller.File().RegisterFormat(controller.FormatSnap, snapFormats)
	spectrum.controller.File().RegisterFormat(controller.FormatTape, tapeFormats)
}

// VideoMemory gets the video memory bank
func (spectrum *Spectrum) VideoMemory() *memory.Bank {
	return spectrum.memory.GetBankMap(1).Bank()
}

// Emulation control

// BeginFrame begin emulation frame tasks
func (spectrum *Spectrum) BeginFrame() {
	// Emit cpu maskable interrupt
	spectrum.cpu.Interrupt()
}

// Emulate one machine step
func (spectrum *Spectrum) Emulate() {
	// Exetues a CPU instruction
	spectrum.cpu.Execute()
}

// EndFrame end emulation frame tasks
func (spectrum *Spectrum) EndFrame() {}

// Snapshots : load & save state

// LoadFile loads a file into machine
func (spectrum *Spectrum) LoadFile(name string) {
	filefmt, ext := spectrum.controller.File().FileFormat(name)
	if filefmt == controller.FormatUnknown {
		return
	}
	data, err := spectrum.controller.File().LoadFileFormat(name, filefmt)
	if err != nil {
		return
	}
	// load snapshop formats
	if filefmt == controller.FormatSnap {
		var snap *format.Snapshot
		switch ext {
		case formatSNA:
			snap = format.LoadSNA(data)
		case formatZ80:
			snap = format.LoadZ80(data)
		default:
			// not supported format
		}
		if snap != nil {
			spectrum.LoadState(snap)
		}
	} else if filefmt == controller.FormatTape {
		// TODO not implemented
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
	spectrum.memory.Load(0x4000, snap.Memory[0:0xc000])
}
