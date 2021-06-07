package passport

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDirectory(t *testing.T) {
	t.Run("Where Given Path Is Empty", func(t *testing.T) {
		fs := &osFilesys{}
		const path = " "
		err := fs.EnsureDirectory(path)
		assert.Equal(t, ErrPathEmpty, err)
	})

	t.Run("Where Directory Does Not Exist", func(t *testing.T) {
		fs := &osFilesys{}
		const path = "TestEnsureDirectory-1"
		err := fs.EnsureDirectory(path)
		assert.NoError(t, err)

		_, err = os.ReadDir(path)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(path)
		})
	})

	t.Run("Where Directory Already Exists", func(t *testing.T) {
		fs := &osFilesys{}
		const path = "TestEnsureDirectory-2"
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}

		err = fs.EnsureDirectory(path)
		assert.NoError(t, err)

		_, err = os.ReadDir(path)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(path)
		})
	})
}

func TestOsFilesys_Read(t *testing.T) {
	t.Run("Given Empty Path", func(t *testing.T) {
		fs := NewFilesys()
		bytes, err := fs.Read("")
		assert.Nil(t, bytes)
		assert.Equal(t, ErrPathEmpty, err)
	})

	t.Run("Given Valid Path", func(t *testing.T) {
		testData := []byte("Hello World")
		testPath := "TestOsFilesys_Read1.txt"

		f, err := os.Create(testPath)
		if err != nil {
			panic(err)
		}

		f.Write(testData)
		f.Close()

		t.Cleanup(func() {
			os.Remove(testPath)
		})

		fs := NewFilesys()
		data, err := fs.Read(testPath)
		assert.NoError(t, err)
		assert.Equal(t, testData, data)
	})
}
