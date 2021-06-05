package passport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_AddSecret(t *testing.T) {
	cnf := &Config{}

	t.Run("Given Empty Name", func(t *testing.T) {
		err := cnf.AddSecret("", "myValue", false)
		assert.Equal(t, ErrSecretNameEmpty, err)
	})

	t.Run("Given Empty Value", func(t *testing.T) {
		err := cnf.AddSecret("myName", "", false)
		assert.Equal(t, ErrSecretValueEmpty, err)
	})
}
