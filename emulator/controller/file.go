package controller

import (
	"path/filepath"
	"strings"
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

// Common extensions
const (
	ExtRom = "rom"
	ExtZip = "zip"
)

// -----------------------------------------------------------------------------
// File information
// -----------------------------------------------------------------------------

// FileInfo contains file information
type FileInfo struct {
	Path   string // File path
	Name   string // File name
	Ext    string // File extension
	Format int    // File format
	IsZip  bool   // Is a zip file
	Data   []byte // File data
}

// NewFileInfo returns new FileInfo
func NewFileInfo(path string) *FileInfo {
	info := new(FileInfo)
	info.Path = path
	info.Name = filepath.Base(path)
	info.Format = FormatUnknown
	info.extension()
	return info
}

// extension validates file extension and zip
func (info *FileInfo) extension() {
	const extIsZip = "." + ExtZip
	name := strings.ToLower(info.Name)
	ext := filepath.Ext(name)
	if ext == extIsZip { // check if is a zip file
		info.IsZip = true
		ext = filepath.Ext(name[:len(name)-4])
	}
	if ext != "" {
		info.Ext = ext[1:]
	}
}

// -----------------------------------------------------------------------------
// Virtual File System
// -----------------------------------------------------------------------------

// VFileSystem the virtual filesystem interface
type VFileSystem interface {
	// LoadFile loads the file data from it's storage location.
	LoadFile(info *FileInfo) error
	// SaveFile saves fhe file data to it's storage location.
	SaveFile(info *FileInfo) error
}

// DefaultFileSystem the current filesystem
var DefaultFileSystem VFileSystem

// -----------------------------------------------------------------------------
// File Manager
// -----------------------------------------------------------------------------

// FileManager is the emulator files manager
type FileManager struct {
	vfs     VFileSystem    // The underlying virtual file system
	formats map[string]int // The format extension mapping
}

// NewFileManager returns a new file manager
func NewFileManager() *FileManager {
	manager := new(FileManager)
	manager.vfs = DefaultFileSystem
	manager.formats = make(map[string]int)
	manager.AddFormat(FormatRom, ExtRom)
	return manager
}

// SetFileSystem set virtual file system
func (manager *FileManager) SetFileSystem(vfs VFileSystem) {
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
