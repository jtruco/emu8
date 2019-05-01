package spectrum

import (
	"io/ioutil"

	"github.com/jtruco/emu8/cpu"
	"github.com/jtruco/emu8/cpu/z80"
	"github.com/jtruco/emu8/device"
	"github.com/jtruco/emu8/device/audio"
	"github.com/jtruco/emu8/device/memory"
	"github.com/jtruco/emu8/machine"
)

// -----------------------------------------------------------------------------
// ZX Spectrum
// -----------------------------------------------------------------------------

// Default ZX Spectrum constants
const (
	// fps          = 50.08 // 50.08 Hz
	fps            = 50    // 50 Hz
	frameTStates   = 69888 // TStates per frame
	audioFrecuency = 48000 // 48 KHz
)

// Spectrum the ZX Spectrum
type Spectrum struct {
	config     machine.Config     // Machine information
	controller machine.Controller // The machine controller
	components *device.Components // Machine device components
	clock      *cpu.ClockDevice   // The system clock
	cpu        *z80.Z80           // The Zilog Z80A CPU
	memory     *memory.Memory     // The machine memory
	ula        *ULA               // The spectrum ULA
	tv         *TVVideo           // The spectrum TV video output
	beeper     *audio.Beeper      // The spectrum Beeper
	keyboard   *Keyboard          // The spectrum Keyboard
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
	spectrum.tv = NewTVVideo(spectrum.memory.GetBankMap(1).Bank())
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
	spectrum.loadROM()
}

// Reset resets the machine
func (spectrum *Spectrum) Reset() {
	// reset components
	spectrum.components.Reset()
	// resets spectrum
	spectrum.loadROM()
}

// Machine properties

// Config gets the machine info
func (spectrum *Spectrum) Config() *machine.Config {
	return &spectrum.config
}

// CPU gets the machine clock
func (spectrum *Spectrum) CPU() cpu.CPU {
	return spectrum.cpu
}

// Components gets the machine components
func (spectrum *Spectrum) Components() *device.Components {
	return spectrum.components
}

// SetController connect controllers & components
func (spectrum *Spectrum) SetController(controller machine.Controller) {
	spectrum.controller = controller
	controller.Video().SetVideo(spectrum.tv)
	controller.Audio().SetAudio(spectrum.beeper)
	controller.Keyboard().AddReceiver(spectrum.keyboard, zxKeyboardMap)
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

// LoadState loads a ZX Spectrum snapshot
func (spectrum *Spectrum) LoadState(snap *Snapshot) {
	// CPU
	// FIXME implement Z80 State copy
	spectrum.cpu.A = snap.A
	spectrum.cpu.F = snap.F
	spectrum.cpu.B = snap.B
	spectrum.cpu.C = snap.C
	spectrum.cpu.D = snap.D
	spectrum.cpu.E = snap.E
	spectrum.cpu.H = snap.H
	spectrum.cpu.L = snap.L
	spectrum.cpu.Ax = snap.Ax
	spectrum.cpu.Fx = snap.Fx
	spectrum.cpu.Bx = snap.Bx
	spectrum.cpu.Cx = snap.Cx
	spectrum.cpu.Dx = snap.Dx
	spectrum.cpu.Ex = snap.Ex
	spectrum.cpu.Hx = snap.Hx
	spectrum.cpu.Lx = snap.Lx
	spectrum.cpu.IXl = snap.IXl
	spectrum.cpu.IXh = snap.IXh
	spectrum.cpu.IYl = snap.IYl
	spectrum.cpu.IYh = snap.IYh
	spectrum.cpu.I = snap.I
	spectrum.cpu.IFF1 = snap.IFF1
	spectrum.cpu.IFF2 = snap.IFF2
	spectrum.cpu.IM = snap.IM
	spectrum.cpu.R = snap.R
	spectrum.cpu.PC = snap.PC
	spectrum.cpu.SP = snap.SP
	// TStates
	spectrum.clock.SetTstates(int(snap.Tstates))
	// Border
	spectrum.tv.SetBorder(snap.Border)
	// Memory
	spectrum.memory.Load(0x4000, snap.Memory[0:0xc000])
}

// -----------------------------------------------------------------------------
// TODO : ROM & FILES
// -----------------------------------------------------------------------------

// LoadFile loads a snapshot file
func (spectrum *Spectrum) LoadFile(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	// currently only SNA format
	state := LoadSNA(data)
	if state != nil {
		spectrum.LoadState(state)
	}
}

func (spectrum *Spectrum) loadROM() {
	data, err := ioutil.ReadFile("48.rom")
	if err != nil {
		return
	}
	// spectrum.memory.Load(0, data[0:0x4000])
	rom := spectrum.memory.GetBankMap(0).Bank()
	rom.Load(0, data[0:0x4000])
}

// -----------------------------------------------------------------------------