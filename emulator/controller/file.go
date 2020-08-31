package controller

import (
	"archive/zip"
	"bytes"
	"io"
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
	PathRom  = "roms"  // ROMs default subpath
	PathSnap = "snaps" // Snapshots default subpath
	PathTape = "tapes" // Tapes default subpath
)

// File formats
const (
	FormatUnknown = iota
	FormatRom
	FormatSnap
	FormatTape
	_FormatMax // limit count
)

// Other constants
const (
	ExtRom    = "rom"
	ExtZip    = "zip"
	_FileMode = 0664
)

// FileManager is the emulator files manager
type FileManager struct {
	path     string             // The file manager base path
	subpaths [_FormatMax]string // Subpaths by file format
	formats  map[string]int     // The format extension mapping
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
	manager.AddFormat(FormatRom, ExtRom)
	manager.SetPath(path)
	return manager
}

// SetPath sets the base path of the file manager
func (manager *FileManager) SetPath(path string) {
	manager.path = path
	manager.subpaths[FormatUnknown] = path
	manager.subpaths[FormatRom] = filepath.Join(path, PathRom)
	manager.subpaths[FormatSnap] = filepath.Join(path, PathSnap)
	manager.subpaths[FormatTape] = filepath.Join(path, PathTape)
}

// Format management

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

// FormatPath gets filename from format path
func (manager *FileManager) FormatPath(format int, filename string) string {
	return filepath.Join(manager.subpaths[format], filepath.Base(filename))
}

// Load files

// LoadROM loads a file from ROMs path
func (manager *FileManager) LoadROM(filename string) ([]byte, error) {
	return ioutil.ReadFile(manager.FormatPath(FormatRom, filename))
}

// LoadFile loads a base filename from its format default location
func (manager *FileManager) LoadFile(file *FileInfo) error {
	data, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return err
	}
	if file.IsZip {
		return file.Unzip(data)
	}
	file.Data = data
	return err
}

// Save files

// SaveFile saves data to a new file
func (manager *FileManager) SaveFile(filename string, format int, data []byte) error {
	// find base file in standar location
	filename = filepath.Join(manager.subpaths[format], filepath.Base(filename))
	return ioutil.WriteFile(filename, data, _FileMode)
}

// NewName helper funcion to obtain a new filename
func (manager *FileManager) NewName(prefix, ext string) string {
	now := time.Now().Format("20060102030405")
	return (prefix + "_" + now + "." + ext)
}

// -----------------------------------------------------------------------------
// File information
// -----------------------------------------------------------------------------

// FileInfo contains file information
type FileInfo struct {
	Name   string // File name
	Path   string // File path
	Format int    // File format
	Ext    string // File format extension
	Data   []byte // File data
	IsZip  bool   // Is a zip file
}

// NewFileInfo returns new FileInfo
func NewFileInfo(filename string) *FileInfo {
	info := new(FileInfo)
	info.Name = filepath.Base(filename)
	info.checkExtension()
	info.Path = filename
	info.Format = FormatUnknown
	return info
}

// checkExtension validates file extension and zip
func (info *FileInfo) checkExtension() {
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

// Unzip unzips file data
func (info *FileInfo) Unzip(zipdata []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		name := strings.ToLower(f.Name)
		if strings.HasSuffix(name, info.Ext) {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			var buffer bytes.Buffer
			_, err = io.Copy(&buffer, rc)
			info.Data = buffer.Bytes()
			return err
		}
	}
	return nil // no contents in zip file
}

// FileInfo returns the file information
func (manager *FileManager) FileInfo(filename string) *FileInfo {
	info := NewFileInfo(filename)
	// check extension format
	format, ok := manager.formats[info.Ext]
	if ok {
		info.Format = format
		// check file path
		_, err := os.Stat(info.Path)
		if os.IsNotExist(err) {
			info.Path = manager.FormatPath(info.Format, info.Name)
		}
	}
	return info
}
