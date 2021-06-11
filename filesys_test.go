package passport

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gofrs/flock"
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

func TestOsFilesys_Write(t *testing.T) {
	t.Run("Given Empty Path", func(t *testing.T) {
		fs := NewFilesys()
		err := fs.Write("", []byte("Hello World"))
		assert.Equal(t, ErrPathEmpty, err)
	})

	t.Run("Given Valid Path", func(t *testing.T) {
		testPath := "TestOsFilesys_Write1"
		testData := []byte("Hello World")

		fs := NewFilesys()
		err := fs.Write(testPath, testData)
		assert.Nil(t, err)

		t.Cleanup(func() {
			os.Remove(testPath)
		})

		data, err := ioutil.ReadFile(testPath)
		assert.NoError(t, err)
		assert.Equal(t, testData, data)
	})

	t.Run("Where File Already Exists", func(t *testing.T) {
		testPath := "TestOsFilesys_Write2"
		testData := []byte("Hello World")

		f, err := os.Create(testPath)
		if err != nil {
			panic(err)
		}

		f.Write([]byte("My File"))
		f.Close()

		fs := NewFilesys()
		err = fs.Write(testPath, testData)
		assert.Nil(t, err)

		t.Cleanup(func() {
			os.Remove(testPath)
		})

		data, err := ioutil.ReadFile(testPath)
		assert.NoError(t, err)
		assert.Equal(t, testData, data)
	})

	t.Run("Where Directory Does Not Exist", func(t *testing.T) {
		testPath := "TestOsFilesys_Write3/myfile.txt"
		testData := []byte("Hello World")

		fs := NewFilesys()
		err := fs.Write(testPath, testData)
		assert.Equal(t, ErrDirNotExists, err)
	})

	t.Run("Where File Is In Use", func(t *testing.T) {
		testPath := "TestOsFilesys_Write4"
		testData := []byte("Hello World")

		f, err := os.Create(testPath)
		if err != nil {
			panic(err)
		}

		f.Write([]byte("My File"))
		f.Close()
		fl := flock.New(testPath)
		fl.Lock()

		t.Cleanup(func() {
			fl.Unlock()
			os.Remove(testPath)
		})

		fs := NewFilesys()
		err = fs.Write(testPath, testData)
		assert.Equal(t, ErrFileInUse, err)
	})
}

func TestOsFilesys_FileExists(t *testing.T) {
	t.Run("Given Empty Path", func(t *testing.T) {
		fs := NewFilesys()
		ok, err := fs.FileExists("")
		assert.False(t, ok)
		assert.Equal(t, ErrPathEmpty, err)
	})

	t.Run("Where File Exists", func(t *testing.T) {
		const testPath = "TestOsFilesys_FileExists1"
		err := os.WriteFile(testPath, []byte("hello"), 0006)
		if err != nil {
			panic(err)
		}

		t.Cleanup(func() {
			os.Remove(testPath)
		})

		fs := NewFilesys()
		ok, err := fs.FileExists(testPath)
		assert.True(t, ok)
		assert.NoError(t, err)
	})

	t.Run("Where File Does Not Exists", func(t *testing.T) {
		const testPath = "TestOsFilesys_FileExists1"

		fs := NewFilesys()
		ok, err := fs.FileExists(testPath)
		assert.False(t, ok)
		assert.NoError(t, err)
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
