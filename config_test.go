package passport

import (
	"errors"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/passport/mock"
)

func TestEnsureConfigFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Where File Does Not Exist", func(t *testing.T) {
		testDir := ".config"
		testFilePath := path.Join(testDir, configFilename)

		fs := mock.NewMockFilesys(ctrl)
		fs.EXPECT().FileExists(testFilePath).Return(false, nil)
		fs.EXPECT().Write(testFilePath, gomock.Any()).Return(nil)

		err := EnsureConfigFile(testDir, fs)
		assert.NoError(t, err)
	})

	t.Run("Where FileExists Returns Error", func(t *testing.T) {
		testDir := ".config"
		testErr := errors.New("fs: error")

		fs := mock.NewMockFilesys(ctrl)
		fs.EXPECT().FileExists(path.Join(testDir, configFilename)).Return(false, testErr)

		err := EnsureConfigFile(testDir, fs)
		assert.Equal(t, testErr, err)
	})

	t.Run("Where File Already Exists", func(t *testing.T) {
		testDir := ".config"

		fs := mock.NewMockFilesys(ctrl)
		fs.EXPECT().FileExists(path.Join(testDir, configFilename)).Return(true, nil)

		err := EnsureConfigFile(testDir, fs)
		assert.NoError(t, err)
	})

	t.Run("Where Write Returns Error", func(t *testing.T) {
		testDir := ".config"
		testFilePath := path.Join(testDir, configFilename)
		testErr := errors.New("fs: error")

		fs := mock.NewMockFilesys(ctrl)
		fs.EXPECT().FileExists(testFilePath).Return(false, nil)
		fs.EXPECT().Write(testFilePath, gomock.Any()).Return(testErr)

		err := EnsureConfigFile(testDir, fs)
		assert.Equal(t, testErr, err)
	})
}

func TestLoadConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Where Config File Exists", func(t *testing.T) {
		testDir := ".config"
		testFilePath := path.Join(testDir, configFilename)
		testData := `secrets:
- name: MySecret
  value: Hello World
  secure: false`

		fs := mock.NewMockFilesys(ctrl)
		fs.EXPECT().Read(testFilePath).Return([]byte(testData), nil)

		c, err := LoadConfig(testDir, fs)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(c.Secrets))
		assert.Equal(t, "MySecret", c.Secrets[0].Name)
		assert.Equal(t, "Hello World", c.Secrets[0].Value)
		assert.Equal(t, false, c.Secrets[0].Secure)
	})

	t.Run("Where Read Fails", func(t *testing.T) {
		testDir := ".config"
		testFilePath := path.Join(testDir, configFilename)
		testErr := errors.New("fs: error")

		fs := mock.NewMockFilesys(ctrl)
		fs.EXPECT().Read(testFilePath).Return(nil, testErr)

		c, err := LoadConfig(testDir, fs)
		assert.Nil(t, c)
		assert.Equal(t, testErr, err)
	})
}

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

func TestConfig_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Saves Config", func(t *testing.T) {
		mockFilesys := mock.NewMockFilesys(ctrl)
		mockFilesys.EXPECT().Write("config/"+configFilename, []byte("secrets: []\n")).Return(nil)

		cnf := &Config{
			configDir: "config",
			fs:        mockFilesys,
			Secrets:   make([]Secret, 0),
		}

		err := cnf.Save()
		assert.Nil(t, err)
	})

	t.Run("Write Fails", func(t *testing.T) {
		testError := errors.New("filesys: test error")

		mockFilesys := mock.NewMockFilesys(ctrl)
		mockFilesys.EXPECT().Write("config/"+configFilename, []byte("secrets: []\n")).Return(testError)

		cnf := &Config{
			configDir: "config",
			fs:        mockFilesys,
			Secrets:   make([]Secret, 0),
		}

		err := cnf.Save()
		assert.Equal(t, testError, err)
	})
}

func TestSecret_GetValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Where Secret Is Secure", func(t *testing.T) {
		testSecureValue := "3287ykshd"
		testValue := "hello world"

		s := &Secret{
			Value:  testSecureValue,
			Secure: true,
		}

		cp := mock.NewMockCryptoProvider(ctrl)
		cp.EXPECT().DecryptString(testSecureValue).Return(testValue, nil)

		v := s.GetValue(cp)
		assert.Equal(t, testValue, v)
	})

	t.Run("Where Secret Is Plain-Text", func(t *testing.T) {
		testValue := "hello world"

		s := &Secret{
			Value:  testValue,
			Secure: false,
		}

		cp := mock.NewMockCryptoProvider(ctrl)

		v := s.GetValue(cp)
		assert.Equal(t, testValue, v)
	})

	t.Run("Where CryptoProvider Fails", func(t *testing.T) {
		testSecureValue := "3287ykshd"
		testError := errors.New("crypto: test error")

		s := &Secret{
			Value:  testSecureValue,
			Secure: true,
		}

		cp := mock.NewMockCryptoProvider(ctrl)
		cp.EXPECT().DecryptString(testSecureValue).Return("", testError)

		v := s.GetValue(cp)
		assert.Equal(t, "", v)
	})
}
