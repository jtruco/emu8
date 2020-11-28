package vfs

import "errors"

// -----------------------------------------------------------------------------
// Virtual File System
// -----------------------------------------------------------------------------

// fileSystem is the current filesystem
var fileSystem FileSystem = NewMemFileSystem()

// GetFileSystem returns the current filesystem
func GetFileSystem() FileSystem {
	return fileSystem
}

// SetFileSystem returns the current filesystem
func SetFileSystem(fs FileSystem) {
	fileSystem = fs
}

// FileSystem is a virtual filesystem
type FileSystem interface {
	// LoadFile loads the file data from the storage location.
	LoadFile(info *FileInfo) error
	// SaveFile saves fhe file data to the storage location.
	SaveFile(info *FileInfo) error
}

// MemFileSystem is a simple in-memory filesystem
type MemFileSystem struct {
	files map[string][]byte
}

// NewMemFileSystem creates a new in-memory filesystem
func NewMemFileSystem() *MemFileSystem {
	mfs := new(MemFileSystem)
	mfs.files = make(map[string][]byte)
	return mfs
}

// LoadFile loads the file data from memory
func (mfs *MemFileSystem) LoadFile(info *FileInfo) error {
	data := mfs.files[info.Path]
	if data == nil {
		return errors.New("MemoryFileSystem : file not found")
	}
	info.Data = data
	return nil
}

// SaveFile stores fhe file data into memory
func (mfs *MemFileSystem) SaveFile(info *FileInfo) error {
	mfs.files[info.Path] = info.Data
	return nil
}
