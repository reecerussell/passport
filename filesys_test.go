package passport

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDirectory(t *testing.T) {
	t.Run("Where Given Path Is Empty", func(t *testing.T) {
		const path = " "
		err := EnsureDirectory(path)
		assert.Equal(t, ErrPathEmpty, err)
	})

	t.Run("Where Directory Does Not Exist", func(t *testing.T) {
		const path = "TestEnsureDirectory-1"
		err := EnsureDirectory(path)
		assert.NoError(t, err)

		_, err = os.ReadDir(path)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(path)
		})
	})

	t.Run("Where Directory Already Exists", func(t *testing.T) {
		const path = "TestEnsureDirectory-2"
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}

		err = EnsureDirectory(path)
		assert.NoError(t, err)

		_, err = os.ReadDir(path)
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove(path)
		})
	})
}
