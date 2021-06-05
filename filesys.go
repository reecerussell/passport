package passport

import (
	"errors"
	"os"
	"strings"
)

// This file contains utils for interacting with the os' file system.

var (
	ErrPathEmpty = errors.New("filesys: path can not be empty")
)

// EnsureDirectory attempts to ensure a given directory exists,
// by creating it if it does not. The directory will be created
// with os.ModePerm.
func EnsureDirectory(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return ErrPathEmpty
	}

	return os.MkdirAll(path, os.ModePerm)
}

// Write writes data to dir+file. If the file does not
// exist, it will be created, otherwise overwritten.
func Write(path string, data []byte) error {
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
func FileExists(path string) (bool, error) {
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
