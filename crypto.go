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

// CryptoProvider is an interface used to abstract encryption/decryption logic.
type CryptoProvider interface {
	EncryptString(value string) (string, error)
	DecryptString(value string) (string, error)
}

type hostCryptoProvider struct{}

// NewCryptoProvider returns a new instance of CryptoProvider.
func NewCryptoProvider() CryptoProvider {
	return &hostCryptoProvider{}
}

// EncryptString encrypts a string value using AES256, with a
// key generated from a host machine's unique identifier.
func (p *hostCryptoProvider) EncryptString(value string) (string, error) {
	key, err := p.generateEncryptionKey()
	if err != nil {
		return "", err
	}

	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	data := gcm.Seal(nonce, nonce, []byte(value), nil)
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecryptString decrypts a string value using AES256, with a
// key generated from a host machine's unique identifier. If
// value is invalid or cannot be decrypted, ErrDecryptFailed
// will be returned.
func (p *hostCryptoProvider) DecryptString(value string) (string, error) {
	key, err := p.generateEncryptionKey()
	if err != nil {
		return "", err
	}

	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)
	nonceSize := gcm.NonceSize()
	cipherText, _ := base64.StdEncoding.DecodeString(value)
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}

	return string(plainText), nil
}

func (p *hostCryptoProvider) generateEncryptionKey() ([]byte, error) {
	mid, err := p.getMachineID()
	if err != nil {
		return nil, err
	}

	sha := sha256.New()
	sha.Write([]byte(mid))

	return sha.Sum(nil), nil
}
