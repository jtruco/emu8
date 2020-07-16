package controller

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// File Manager
// -----------------------------------------------------------------------------

// Path constants
const (
	defaultRomPath  = "roms"  // ROMs default subpath
	defaultSnapPath = "snaps" // Snapshots default subpath
	defaultTapePath = "tapes" // Tapes default subpath
)

// File formats
const (
	FormatUnknown = iota
	FormatRom
	FormatSnap
	FormatTape
	_FormatMax // limit count
)

const defaultFileMode = 0664

// FileManager is the emulator files manager
type FileManager struct {
	path     string             // The file manager base path
	subpaths [_FormatMax]string // Subpaths by file format
	formats  map[string]int     // The file type extension mapping
}

// DefaultFileManager returns the default file manager
func DefaultFileManager() *FileManager {
	dir, _ := os.Getwd()
	return NewFileManager(dir)
}

// NewFileManager returns a new file manager
func NewFileManager(path string) *FileManager {
	manager := new(FileManager)
	manager.formats = make(map[string]int)
	manager.SetPath(path)
	return manager
}

// Path management

// SetPath sets the base path of the file manager
func (manager *FileManager) SetPath(path string) {
	manager.path = path
	manager.subpaths[FormatUnknown] = path
	manager.subpaths[FormatRom] = filepath.Join(path, defaultRomPath)
	manager.subpaths[FormatSnap] = filepath.Join(path, defaultSnapPath)
	manager.subpaths[FormatTape] = filepath.Join(path, defaultTapePath)
}

// File management

// LoadROM loads a file from ROMs path
func (manager *FileManager) LoadROM(name string) ([]byte, error) {
	return manager.LoadFileFormat(name, FormatRom)
}

// LoadFileFormat loads a base filename from its format default location
func (manager *FileManager) LoadFileFormat(filename string, format int) ([]byte, error) {
	data, err := manager.LoadFile(filename)
	if err == nil || format == FormatUnknown {
		return data, err
	}
	// find base file in standar location
	filename = filepath.Join(manager.subpaths[format], filepath.Base(filename))
	return manager.LoadFile(filename)
}

// LoadFile loads a file and return its data bytes
func (manager *FileManager) LoadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// SaveFileFormat saves data to a new file
func (manager *FileManager) SaveFileFormat(filename string, format int, data []byte) error {
	// find base file in standar location
	filename = filepath.Join(manager.subpaths[format], filepath.Base(filename))
	return manager.SaveFile(filename, data)
}

// SaveFile saves data to a new file
func (manager *FileManager) SaveFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, defaultFileMode)
}

// BaseName helper funcion to obtain base filename
func (manager *FileManager) BaseName(filename string) string {
	return filepath.Base(filename)
}

// NewName helper funcion to obtain a new filename
func (manager *FileManager) NewName(prefix, ext string) string {
	now := time.Now().Format("20060102030405")
	return (prefix + "_" + now + "." + ext)
}

// File extension type management

// RegisterFormat adds a format and its extensions
func (manager *FileManager) RegisterFormat(format int, extensions []string) {
	for _, ext := range extensions {
		manager.AddFormat(format, ext)
	}
}

// AddFormat adds a file extension format
func (manager *FileManager) AddFormat(format int, extension string) {
	manager.formats[extension] = format
}

// FileFormat detects and returns the file machine format and supported extension
func (manager *FileManager) FileFormat(filename string) (int, string) {
	extension := filepath.Ext(filename)
	if len(extension) > 0 {
		extension = strings.ToLower(extension[1:])
		format, ok := manager.formats[extension]
		if ok {
			return format, extension
		}
	}
	return FormatUnknown, extension
}
