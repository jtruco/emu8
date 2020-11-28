package vfs

import (
	"time"
)

// -----------------------------------------------------------------------------
// File Constants
// -----------------------------------------------------------------------------

// File formats
const (
	FormatUnknown = iota
	FormatRom
	FormatSnap
	FormatTape
	FormatMax // limit count
)

// Commont extensions
const (
	ExtRom = "rom"
	ExtZip = "zip"
)

// -----------------------------------------------------------------------------
// File Manager
// -----------------------------------------------------------------------------

// FileManager is the emulator files manager
type FileManager struct {
	vfs     FileSystem     // The underlying virtual file system
	formats map[string]int // The format extension mapping
}

// NewFileManager returns a new file manager
func NewFileManager() *FileManager {
	manager := new(FileManager)
	manager.vfs = GetFileSystem()
	manager.formats = make(map[string]int)
	manager.AddFormat(FormatRom, ExtRom)
	return manager
}

// SetFileSystem set virtual file system
func (manager *FileManager) SetFileSystem(vfs FileSystem) {
	manager.vfs = vfs
}

// Format management

// AddFormat adds a file extension format
func (manager *FileManager) AddFormat(format int, extension string) {
	manager.formats[extension] = format
}

// RegisterFormat adds a format and its extensions
func (manager *FileManager) RegisterFormat(format int, extensions []string) {
	for _, ext := range extensions {
		manager.AddFormat(format, ext)
	}
}

// Load & Save Files

// LoadROM loads a file from ROMs path
func (manager *FileManager) LoadROM(filename string) ([]byte, error) {
	info := NewFileInfo(filename)
	info.Format = FormatRom
	return info.Data, manager.vfs.LoadFile(info)
}

// LoadFile loads a base filename from its format default location
func (manager *FileManager) LoadFile(info *FileInfo) error {
	return manager.vfs.LoadFile(info)
}

// SaveFile saves data to a new file
func (manager *FileManager) SaveFile(filename string, format int, data []byte) error {
	info := NewFileInfo(filename)
	info.Format = format
	info.Data = data
	return manager.vfs.SaveFile(info)
}

// FileInfo

// NewName helper funcion to obtain a new filename
func (manager *FileManager) NewName(prefix, ext string) string {
	now := time.Now().Format("20060102030405")
	return (prefix + "_" + now + "." + ext)
}

// FileInfo returns the file information
func (manager *FileManager) FileInfo(filename string) *FileInfo {
	info := NewFileInfo(filename)
	// check extension format
	format, ok := manager.formats[info.Ext]
	if ok {
		info.Format = format
	}
	return info
}
