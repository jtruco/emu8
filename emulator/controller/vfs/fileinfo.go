package vfs

import (
	"path/filepath"
	"strings"
)

const extIsZip = "." + ExtZip

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

// NewFileInfo returns new file info
func NewFileInfo(path string) *FileInfo {
	info := new(FileInfo)
	info.Path = path
	info.Name = filepath.Base(path)
	info.Format = FormatUnknown
	info.extension()
	return info
}

// NewFileInfoData returns new file info with data
func NewFileInfoData(path string, data []byte) *FileInfo {
	info := NewFileInfo(path)
	info.Data = data
	return info
}

// extension validates file extension and zip
func (info *FileInfo) extension() {
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
