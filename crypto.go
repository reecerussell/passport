package passport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var (
	// ErrDecryptFailed is an error indicating that a
	// value could not be decryped.
	ErrDecryptFailed = errors.New("decrypt: failed to decrypt data")
)

// EncryptString encrypts a string value using AES256, with a
// key generated from a host machine's unique identifier.
func EncryptString(value string) (string, error) {
	key, err := generateEncryptionKey()
	if err != nil {
		return "", nil
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	data := gcm.Seal(nonce, nonce, []byte(value), nil)
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecryptString decrypts a string value using AES256, with a
// key generated from a host machine's unique identifier. If
// value is invalid or cannot be decrypted, ErrDecryptFailed
// will be returned.
func DecryptString(value string) (string, error) {
	key, err := generateEncryptionKey()
	if err != nil {
		return "", nil
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(value) < nonceSize {
		return "", ErrDecryptFailed
	}

	cipherText, _ := base64.StdEncoding.DecodeString(value)
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}

	return string(plainText), nil
}

func generateEncryptionKey() ([]byte, error) {
	mid, err := getMachineID()
	if err != nil {
		return nil, err
	}

	sha := sha256.New()
	sha.Write([]byte(mid))

	return sha.Sum(nil), nil
}
