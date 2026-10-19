package controller

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/device"
	"github.com/jtruco/emu8/emulator/device/io/tape"
	"github.com/jtruco/emu8/emulator/machine"
)

// panicMachine is a machine whose snapshot and tape loaders panic
type panicMachine struct {
	config machine.Config
	loads  int
}

func (m *panicMachine) Init() error                    { return nil }
func (m *panicMachine) Reset()                         {}
func (m *panicMachine) Config() *machine.Config        { return &m.config }
func (m *panicMachine) Clock() device.Clock            { return device.NewClock() }
func (m *panicMachine) Components() *device.Components { return device.NewComponents() }
func (m *panicMachine) Emulate()                       {}
func (m *panicMachine) BeginFrame()                    {}
func (m *panicMachine) EndFrame()                      {}
func (m *panicMachine) SaveState() machine.State       { return machine.State{} }
func (m *panicMachine) LoadState(state machine.State)  { m.loads++; panic("malformed snapshot") }
func (m *panicMachine) InitControl(ctl machine.Control) {
	ctl.BindTapeDrive(tape.New(device.NewClock()))
	ctl.RegisterSnapshot("snp")
	ctl.RegisterTape("tpe", func() tape.Tape { return panicTape{} })
}

type panicTape struct{}

func (panicTape) Info() *tape.Info       { return new(tape.Info) }
func (panicTape) Blocks() []tape.Block   { return nil }
func (panicTape) Load(data []byte) bool  { panic("malformed tape") }
func (panicTape) Play(ctl *tape.Control) {}

func TestLoadFileRecoversFromPanics(t *testing.T) {
	fs := vfs.NewMemFileSystem()
	m := new(panicMachine)
	controller := New(m, fs)
	for _, name := range []string{"game.snp", "game.tpe"} {
		if err := fs.SaveFile(vfs.NewFileInfoData(name, []byte{0x00})); err != nil {
			t.Fatal(err)
		}
		controller.LoadFile(name) // must not panic
	}
	if m.loads != 1 {
		t.Errorf("snapshot loads = %d, want 1", m.loads)
	}
}
