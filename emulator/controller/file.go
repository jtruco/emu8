package controller

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// -----------------------------------------------------------------------------
// File Manager
// -----------------------------------------------------------------------------

const (
	pathRom      = "roms"  // ROMs default subpath
	pathSnapshot = "snaps" // Snapshots default subpath
)

// FileManager is the emulator files manager
type FileManager struct {
	path         string // The base path
	romPath      string // ROMs path
	snapshotPath string // Snapshots path
}

// DefaultFileManager returns the default file manager
func DefaultFileManager() *FileManager {
	dir, _ := os.Getwd()
	return NewFileManager(dir)
}

// NewFileManager returns a new file manager
func NewFileManager(path string) *FileManager {
	manager := &FileManager{}
	manager.SetPath(path)
	return manager
}

// Paths management

// FilenameROM gets ROM filename full path
func (manager *FileManager) FilenameROM(rom string) string {
	return filepath.Join(manager.romPath, rom)
}

// FilenameSnapshot gets Snapshot filename full path
func (manager *FileManager) FilenameSnapshot(snap string) string {
	return filepath.Join(manager.snapshotPath, snap)
}

// SetPath sets the base path of the file manager
func (manager *FileManager) SetPath(path string) {
	manager.path = path
	manager.romPath = filepath.Join(path, pathRom)
	manager.snapshotPath = filepath.Join(path, pathSnapshot)
}

// Files management

// LoadROM loads a file from ROMs path
func (manager *FileManager) LoadROM(name string) ([]byte, error) {
	return manager.LoadFile(manager.FilenameROM(name))
}

// LoadSnapshot loads a file from Snapshots path
func (manager *FileManager) LoadSnapshot(name string) ([]byte, error) {
	return manager.LoadFile(manager.FilenameSnapshot(name))
}

// LoadFile loads a file and return its data bytes
func (manager *FileManager) LoadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
