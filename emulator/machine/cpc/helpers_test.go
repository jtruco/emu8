package cpc

import (
	"testing"

	"github.com/jtruco/emu8/emulator/controller"
	"github.com/jtruco/emu8/emulator/controller/vfs"
	"github.com/jtruco/emu8/emulator/machine"
)

// romFileSystem returns a memory file system with zeroed ROM files
func romFileSystem() *vfs.MemFileSystem {
	fs := vfs.NewMemFileSystem()
	fs.SaveFile(vfs.NewFileInfoData(cpcOsRomName, make([]byte, 0x4000)))
	fs.SaveFile(vfs.NewFileInfoData(cpcBasicRomName, make([]byte, 0x4000)))
	return fs
}

// newTestCPC builds and initializes a CPC 464 with zeroed ROMs
func newTestCPC(t *testing.T) *AmstradCPC {
	t.Helper()
	m := New(AmstradCPC464, machine.Settings{}).(*AmstradCPC)
	controller.New(m, romFileSystem())
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	return m
}
