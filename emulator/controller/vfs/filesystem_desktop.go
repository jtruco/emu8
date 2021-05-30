// +build darwin freebsd linux windows
// +build !android,!ios,!js

package vfs

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Default subpath constants
const (
	PathRom      = "roms"  // ROMs default subpath
	PathSnapshot = "snaps" // Snapshots default subpath
	PathTape     = "tapes" // Tapes default subpath
)

// -----------------------------------------------------------------------------
// DesktopFileSystem
// -----------------------------------------------------------------------------

// DesktopFileSystem implements the file system for desktop
type DesktopFileSystem struct {
	path     string
	subpaths [FormatMax]string // Subpaths by file format
}

// NewDesktopFileSystem creates a new desktop filesystem
func NewDesktopFileSystem(path string) *DesktopFileSystem {
	fs := new(DesktopFileSystem)
	fs.path = path
	fs.subpaths[FormatUnknown] = path
	fs.subpaths[FormatRom] = filepath.Join(path, PathRom)
	fs.subpaths[FormatSnapshot] = filepath.Join(path, PathSnapshot)
	fs.subpaths[FormatTape] = filepath.Join(path, PathTape)
	return fs
}

// InitDesktop initialices the desktop filesystem
func InitDesktop() {
	cwd, _ := os.Getwd()
	SetFileSystem(NewDesktopFileSystem(cwd))
}

// LoadFile loads the file data from it's storage location.
func (dfs *DesktopFileSystem) LoadFile(info *FileInfo) error {
	// check for file
	err := dfs.stat(info)
	if err != nil {
		return err
	}
	// read and unzip data
	data, err := ioutil.ReadFile(info.Path)
	if err == nil {
		info.Data = data
		if info.IsZip {
			err = dfs.unzip(info)
		}
	}
	return err
}

// SaveFile saves fhe file data to it's storage location.
func (dfs *DesktopFileSystem) SaveFile(info *FileInfo) error {
	const fileMode = 0664
	dfs.formatPath(info)
	return ioutil.WriteFile(info.Path, info.Data, fileMode)
}

// stat checks exists file in default folders
func (dfs *DesktopFileSystem) stat(info *FileInfo) error {
	// check for file
	_, err := os.Stat(info.Path)
	if os.IsNotExist(err) {
		// search in format folder
		dfs.formatPath(info)
		_, err = os.Stat(info.Path)
	}
	return err
}

// FormatPath gets filename from format path
func (dfs *DesktopFileSystem) formatPath(info *FileInfo) {
	info.Path = filepath.Join(dfs.subpaths[info.Format], filepath.Base(info.Name))
}

// Unzip unzips file data
func (dfs *DesktopFileSystem) unzip(file *FileInfo) error {
	zipdata := file.Data
	zr, err := zip.NewReader(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		name := strings.ToLower(f.Name)
		if strings.HasSuffix(name, file.Ext) {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			var buffer bytes.Buffer
			_, err = io.Copy(&buffer, rc)
			file.Data = buffer.Bytes()
			return err
		}
	}
	return nil // no contents in zip file
}
