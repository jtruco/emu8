package io

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	vfs "github.com/jtruco/emu8/emulator/controller"
)

// Path constants
const (
	PathRom  = "roms"  // ROMs default subpath
	PathSnap = "snaps" // Snapshots default subpath
	PathTape = "tapes" // Tapes default subpath
)

// FileSystem implements the virtual file system
type FileSystem struct {
	path     string
	subpaths [vfs.FormatMax]string // Subpaths by file format
}

// NewFileSystem creates a new filesystem
func NewFileSystem(path string) *FileSystem {
	fs := new(FileSystem)
	fs.path = path
	fs.subpaths[vfs.FormatUnknown] = path
	fs.subpaths[vfs.FormatRom] = filepath.Join(path, PathRom)
	fs.subpaths[vfs.FormatSnap] = filepath.Join(path, PathSnap)
	fs.subpaths[vfs.FormatTape] = filepath.Join(path, PathTape)
	return fs
}

// DefaultFileSystem creates the default filesystem
func DefaultFileSystem() {
	cwd, _ := os.Getwd()
	vfs.DefaultFileSystem = NewFileSystem(cwd)
}

// LoadFile loads the file data from it's storage location.
func (fs *FileSystem) LoadFile(info *vfs.FileInfo) error {
	// check for file
	err := fs.stat(info)
	if err != nil {
		return err
	}
	// read and unzip data
	data, err := ioutil.ReadFile(info.Path)
	if err == nil {
		info.Data = data
		if info.IsZip {
			err = fs.unzip(info)
		}
	}
	return err
}

// SaveFile saves fhe file data to it's storage location.
func (fs *FileSystem) SaveFile(info *vfs.FileInfo) error {
	return nil
}

// stat checks exists file in default folders
func (fs *FileSystem) stat(info *vfs.FileInfo) error {
	// check for file
	_, err := os.Stat(info.Path)
	if os.IsNotExist(err) {
		// search in format folder
		fs.formatPath(info)
		_, err = os.Stat(info.Path)
	}
	return err
}

// FormatPath gets filename from format path
func (fs *FileSystem) formatPath(info *vfs.FileInfo) {
	info.Path = filepath.Join(fs.subpaths[info.Format], filepath.Base(info.Name))
}

// Unzip unzips file data
func (fs *FileSystem) unzip(file *vfs.FileInfo) error {
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
