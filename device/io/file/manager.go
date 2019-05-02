package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// -----------------------------------------------------------------------------
// Manager
// -----------------------------------------------------------------------------

const (
	pathRom      = "roms"  // ROMs default subpath
	pathSnapshot = "snaps" // Snapshots default subpath
)

// Manager is the emulator file manager
type Manager struct {
	path         string // The base path
	romPath      string // ROMs path
	snapshotPath string // Snapshots path
}

// DefaultManager returns the default file manager
func DefaultManager() *Manager {
	dir, _ := os.Getwd()
	return NewManager(dir)
}

// NewManager returns a new file manager
func NewManager(path string) *Manager {
	manager := &Manager{}
	manager.SetPath(path)
	return manager
}

// Paths management

// FilenameROM gets ROM filename full path
func (manager *Manager) FilenameROM(rom string) string {
	return filepath.Join(manager.romPath, rom)
}

// FilenameSnapshot gets Snapshot filename full path
func (manager *Manager) FilenameSnapshot(snap string) string {
	return filepath.Join(manager.snapshotPath, snap)
}

// SetPath sets the base path of the file manager
func (manager *Manager) SetPath(path string) {
	manager.path = path
	manager.romPath = filepath.Join(path, pathRom)
	manager.snapshotPath = filepath.Join(path, pathSnapshot)
}

// Files management

// LoadROM loads a file from ROMs path
func (manager *Manager) LoadROM(name string) ([]byte, error) {
	return manager.LoadFile(manager.FilenameROM(name))
}

// LoadSnapshot loads a file from Snapshots path
func (manager *Manager) LoadSnapshot(name string) ([]byte, error) {
	return manager.LoadFile(manager.FilenameSnapshot(name))
}

// LoadFile loads a file and return its data bytes
func (manager *Manager) LoadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
