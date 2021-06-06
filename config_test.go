package passport

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/passport/mock"
)

func TestConfig_AddSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cnf := &Config{
		Secrets: []Secret{
			{
				Name:   "mySecret3",
				Value:  "Hello World",
				Secure: false,
			},
		},
	}

	t.Run("Given Empty Name", func(t *testing.T) {
		cp := mock.NewMockCryptoProvider(ctrl)

		err := cnf.AddSecret("", "myValue", false, cp)
		assert.Equal(t, ErrSecretNameEmpty, err)
	})

	t.Run("Given Empty Value", func(t *testing.T) {
		cp := mock.NewMockCryptoProvider(ctrl)

		err := cnf.AddSecret("myName", "", false, cp)
		assert.Equal(t, ErrSecretValueEmpty, err)
	})

	t.Run("Where CryptoProvider Returns Error", func(t *testing.T) {
		testName := "mySecret0"
		testValue := "233432"
		testError := errors.New("crypto: test error")

		cp := mock.NewMockCryptoProvider(ctrl)
		cp.EXPECT().EncryptString(testValue).Return("", testError)

		err := cnf.AddSecret(testName, testValue, true, cp)
		assert.Equal(t, testError, err)
	})

	t.Run("Where Encrypt Flag Is Passed", func(t *testing.T) {
		testName := "mySecret1"
		testValue := "233432"
		testEncryptedValue := "97324723"

		cp := mock.NewMockCryptoProvider(ctrl)
		cp.EXPECT().EncryptString(testValue).Return(testEncryptedValue, nil)

		err := cnf.AddSecret(testName, testValue, true, cp)
		assert.NoError(t, err)

		found := false
		for _, s := range cnf.Secrets {
			if s.Name == testName {
				found = true
				assert.Equal(t, testEncryptedValue, s.Value)
			}
		}

		assert.True(t, found)
	})

	t.Run("Where Encrypt Flag Is Not Passed", func(t *testing.T) {
		testName := "mySecret2"
		testValue := "233432"

		cp := mock.NewMockCryptoProvider(ctrl)

		err := cnf.AddSecret(testName, testValue, false, cp)
		assert.NoError(t, err)

		found := false
		for _, s := range cnf.Secrets {
			if s.Name == testName {
				found = true
				assert.Equal(t, testValue, s.Value)
			}
		}

		assert.True(t, found)
	})

	t.Run("Where Secret Already Exists", func(t *testing.T) {
		testName := "mySecret3"
		testValue := "233432"

		cp := mock.NewMockCryptoProvider(ctrl)

		err := cnf.AddSecret(testName, testValue, false, cp)
		assert.Equal(t, ErrSecretAlreadyExists, err)
	})
}

func TestConfig_GetSecret(t *testing.T) {
	cnf := &Config{
		Secrets: []Secret{
			{
				Name:   "mySecret",
				Value:  "Hello World",
				Secure: false,
			},
		},
	}

	t.Run("Given Empty Name", func(t *testing.T) {
		s, err := cnf.GetSecret("")
		assert.Nil(t, s)
		assert.Equal(t, ErrSecretNameEmpty, err)
	})

	t.Run("Given Valid Name", func(t *testing.T) {
		s, err := cnf.GetSecret("mySecret")
		assert.NoError(t, err)
		assert.Equal(t, "mySecret", s.Name)
		assert.Equal(t, "Hello World", s.Value)
		assert.False(t, s.Secure)
	})

	t.Run("Where Secret Does Not Exist", func(t *testing.T) {
		s, err := cnf.GetSecret("myFaveSecret")
		assert.Nil(t, s)
		assert.Equal(t, ErrSecretNotFound, err)
	})
}

func TestConfig_RemoveSecret(t *testing.T) {
	cnf := &Config{
		Secrets: []Secret{
			{
				Name:   "mySecret",
				Value:  "Hello World",
				Secure: false,
			},
		},
	}

	t.Run("Given Empty Name", func(t *testing.T) {
		err := cnf.RemoveSecret("")
		assert.Equal(t, ErrSecretNameEmpty, err)
	})

	t.Run("Given Valid Name", func(t *testing.T) {
		err := cnf.RemoveSecret("mySecret")
		assert.NoError(t, err)
	})

	t.Run("Where Secret Does Not Exist", func(t *testing.T) {
		err := cnf.RemoveSecret("myFaveSecret")
		assert.Equal(t, ErrSecretNotFound, err)
	})
}
