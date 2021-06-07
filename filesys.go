package passport

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// This file contains utils for interacting with the os' file system.

var (
	ErrPathEmpty = errors.New("filesys: path can not be empty")
)

// Filesys is a high level interface used to interact with a filesystem.
type Filesys interface {
	// EnsureDirectory is used to ensure a directory exists.
	EnsureDirectory(path string) error

	// Write writes data to path. If the files does not exist,
	// it will be created, otherwise overwritten.
	Write(path string, data []byte) error

	// FileExists returns a boolean which determines if a file
	// at a given path exists.
	FileExists(path string) (bool, error)

	// Read reads all data from a file at path.
	Read(path string) ([]byte, error)
}

type osFilesys struct{}

// NewFilesys returns a new instance of Filesys, based on the host's OS.
func NewFilesys() Filesys {
	return &osFilesys{}
}

// EnsureDirectory attempts to ensure a given directory exists,
// by creating it if it does not. The directory will be created
// with os.ModePerm.
func (*osFilesys) EnsureDirectory(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return ErrPathEmpty
	}

	return os.MkdirAll(path, os.ModePerm)
}

// Write writes data to dir+file. If the file does not
// exist, it will be created, otherwise overwritten.
func (*osFilesys) Write(path string, data []byte) error {
	if path == "" {
		return ErrPathEmpty
	}

	f, err := os.Create(path)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}

		f, err = os.OpenFile(path, os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// FileExists returns a boolean which indicates whether a file
// at the give path exists.
func (*osFilesys) FileExists(path string) (bool, error) {
	if path == "" {
		return false, ErrPathEmpty
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Read reads all data from a file at path. This is a wrapper
// around ioutil.ReadFile.
func (*osFilesys) Read(path string) ([]byte, error) {
	if path == "" {
		return nil, ErrPathEmpty
	}

	return ioutil.ReadFile(path)
}
