package passport

import (
	"errors"
	"io/ioutil"
	"os"
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
		testData := []byte("secrets: []\nworkspaces: []\n")

		mockFilesys := mock.NewMockFilesys(ctrl)
		mockFilesys.EXPECT().Write("config/"+configFilename, testData).Return(nil)

		cnf := &Config{
			configDir:  "config",
			fs:         mockFilesys,
			Secrets:    make([]Secret, 0),
			Workspaces: make([]Workspace, 0),
		}

		err := cnf.Save()
		assert.Nil(t, err)
	})

	t.Run("Write Fails", func(t *testing.T) {
		testData := []byte("secrets: []\nworkspaces: []\n")
		testError := errors.New("filesys: test error")

		mockFilesys := mock.NewMockFilesys(ctrl)
		mockFilesys.EXPECT().Write("config/"+configFilename, testData).Return(testError)

		cnf := &Config{
			configDir:  "config",
			fs:         mockFilesys,
			Secrets:    make([]Secret, 0),
			Workspaces: make([]Workspace, 0),
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

func TestConfig_AddWorkspace(t *testing.T) {
	cnf := &Config{
		Workspaces: []Workspace{
			{
				Name: "MyWorkspace",
				Path: "/c/test",
			},
		},
	}

	t.Run("Given Valid Args", func(t *testing.T) {
		const testName = "MyTestWorkspace"
		const testPath = "/t/home"

		err := cnf.AddWorkspace(testName, testPath)
		assert.NoError(t, err)

		found := false

		for _, w := range cnf.Workspaces {
			if w.Name == testName {
				assert.Equal(t, testPath, w.Path)
				assert.Empty(t, w.Scripts)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("Name Is Empty", func(t *testing.T) {
		err := cnf.AddWorkspace("", "/c/home")
		assert.Equal(t, ErrWorkspaceNameEmpty, err)
	})

	t.Run("Path Is Empty", func(t *testing.T) {
		err := cnf.AddWorkspace("", "/c/home")
		assert.Equal(t, ErrWorkspaceNameEmpty, err)
	})

	t.Run("Name Already Exists", func(t *testing.T) {
		err := cnf.AddWorkspace("MyWorkspace", "/c/home")
		assert.Equal(t, ErrWorkspaceNameExists, err)
	})

	t.Run("Path Already Exists", func(t *testing.T) {
		err := cnf.AddWorkspace("MyNewWorkspace", "/c/test")
		assert.Equal(t, ErrWorkspacePathExists, err)
	})
}

func TestConfig_GetWorkspace(t *testing.T) {
	cnf := &Config{
		Workspaces: []Workspace{
			{
				Name: "MyWorkspace",
				Path: "/c/test",
			},
		},
	}

	t.Run("Given Matching Path", func(t *testing.T) {
		w, err := cnf.GetWorkspace("/c/test")
		assert.NoError(t, err)
		assert.Equal(t, &cnf.Workspaces[0], w)
	})

	t.Run("Given Empty Path", func(t *testing.T) {
		w, err := cnf.GetWorkspace("")
		assert.Nil(t, w)
		assert.Equal(t, ErrWorkspacePathEmpty, err)
	})

	t.Run("Given Invalid Path", func(t *testing.T) {
		w, err := cnf.GetWorkspace("not-a-workspace")
		assert.Nil(t, w)
		assert.Equal(t, ErrWorkspaceNotFound, err)
	})
}

func TestWorkspace_AddScript(t *testing.T) {
	w := &Workspace{
		Name: "MyWorkspace",
		Path: "/c/dev",
		Scripts: []WorkspaceScript{
			{
				Name:    "build",
				Command: "./build.sh",
			},
		},
	}

	t.Run("Given Valid Args", func(t *testing.T) {
		const testName = "MyTestScript"
		const testCommand = "main.exe"

		err := w.AddScript(testName, testCommand)
		assert.NoError(t, err)

		found := false

		for _, s := range w.Scripts {
			if s.Name == testName {
				assert.Equal(t, testCommand, s.Command)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("Name Is Empty", func(t *testing.T) {
		err := w.AddScript("", "./hello.sh")
		assert.Equal(t, ErrWorkspaceScriptNameEmpty, err)
	})

	t.Run("Command Is Empty", func(t *testing.T) {
		err := w.AddScript("my-script", "")
		assert.Equal(t, ErrWorkspaceScriptCommandEmpty, err)
	})

	t.Run("Name Already Exists", func(t *testing.T) {
		err := w.AddScript("build", "./hello.sh")
		assert.Equal(t, ErrWorkspaceScriptNameExists, err)
	})
}

func TestWorkspace_GetScript(t *testing.T) {
	w := &Workspace{
		Scripts: []WorkspaceScript{
			{
				Name: "build",
			},
		},
	}

	t.Run("Given Empty Name", func(t *testing.T) {
		s, err := w.GetScript("")
		assert.Nil(t, s)
		assert.Equal(t, ErrWorkspaceScriptNameEmpty, err)
	})

	t.Run("Given Invalid Name", func(t *testing.T) {
		s, err := w.GetScript("not-a-script")
		assert.Nil(t, s)
		assert.Equal(t, ErrWorkspaceScriptNotFound, err)
	})

	t.Run("Given Valid Name", func(t *testing.T) {
		s, err := w.GetScript("build")
		assert.NoError(t, err)
		assert.Equal(t, &w.Scripts[0], s)
	})
}

func TestWorkspaceScript_Run(t *testing.T) {
	t.Run("Given Valid Command", func(t *testing.T) {
		var command string
		if os.Getenv("GOOS") != "linux" {
			command = "cmd /C echo Hello World"
		} else {
			command = "echo Hello World"
		}

		pr, pw, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		oldStdout := os.Stdout
		os.Stdout = pw

		t.Cleanup(func() {
			pw.Close()
			os.Stdout = oldStdout
		})

		s := WorkspaceScript{
			Command: command,
		}

		code, err := s.Run()
		assert.NoError(t, err)
		assert.Equal(t, 0, code)

		pw.Close()
		os.Stdout = oldStdout

		bytes, _ := ioutil.ReadAll(pr)
		output := string(bytes)

		assert.Contains(t, output, "Hello World")
	})

	t.Run("Given Invalid Command", func(t *testing.T) {
		s := WorkspaceScript{
			Command: "no-a-valid-file.test",
		}

		code, err := s.Run()
		assert.NotNil(t, err)
		assert.Equal(t, -1, code)
	})
}
