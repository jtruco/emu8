// Package vfs contains the virtual file system components
package vfs

import (
	"time"
)

// File formats
const (
	FormatUnknown = iota
	FormatRom
	FormatSnapshot
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

// FileManager is the emulator file manager
type FileManager struct {
	vfs     FileSystem     // The underlying virtual file system
	formats map[string]int // The format extension mapping
}

// NewFileManager returns a new file manager
func NewFileManager() *FileManager {
	manager := new(FileManager)
	manager.vfs = GetFileSystem()
	manager.formats = make(map[string]int)
	manager.RegisterFormat(FormatRom, ExtRom)
	return manager
}

// SetFileSystem set virtual file system
func (manager *FileManager) SetFileSystem(vfs FileSystem) {
	manager.vfs = vfs
}

// Format management

// RegisterFormat adds a file extension format
func (manager *FileManager) RegisterFormat(format int, extension string) {
	manager.formats[extension] = format
}

// RegisterFormats adds a format and its extensions
func (manager *FileManager) RegisterFormats(format int, extensions []string) {
	for _, ext := range extensions {
		manager.RegisterFormat(format, ext)
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

// File information

// CreateFileInfo returns a the file information from filename
func (manager *FileManager) CreateFileInfo(filename string) *FileInfo {
	info := NewFileInfo(filename)
	// check extension format
	format, ok := manager.formats[info.Ext]
	if ok {
		info.Format = format
	}
	return info
}

// NewName helper funcion to obtain a new filename
func (manager *FileManager) NewName(prefix, ext string) string {
	now := time.Now().Format("20060102030405")
	return (prefix + "_" + now + "." + ext)
}
