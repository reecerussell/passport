package passport

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostCryptoProvider(t *testing.T) {
	const testValue = "Hello World"

	cp := NewCryptoProvider()

	encValue, err := cp.EncryptString(testValue)
	assert.NoError(t, err)

	plainValue, err := cp.DecryptString(encValue)
	assert.NoError(t, err)

	assert.Equal(t, testValue, plainValue)
}

func TestHostCryptoProvider_EncryptString(t *testing.T) {
	cp := NewCryptoProvider()

	t.Run("Error Should Be Nil", func(t *testing.T) {
		const testValue = "Hello World"
		result, err := cp.EncryptString(testValue)
		assert.NotEqual(t, testValue, result)
		assert.NotEqual(t, "", testValue)
		assert.NoError(t, err)
	})
}

func TestHostCryptoProvider_DecryptString(t *testing.T) {
	cp := NewCryptoProvider()

	t.Run("Where Key Is Invalid", func(t *testing.T) {
		key := make([]byte, 32)
		rand.Read(key)

		c, _ := aes.NewCipher(key)
		gcm, _ := cipher.NewGCM(c)
		nonce := make([]byte, gcm.NonceSize())
		rand.Read(nonce)

		data := gcm.Seal(nonce, nonce, []byte("Hello World"), nil)
		encValue := base64.StdEncoding.EncodeToString(data)

		result, err := cp.DecryptString(encValue)
		assert.Equal(t, "", result)
		assert.Equal(t, ErrDecryptFailed, err)
	})
}
