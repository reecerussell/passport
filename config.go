package passport

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

const configFilename = "config.yaml"

// Config is a struct which holds and represents the core configuration.
type Config struct {
	configDir string  `yaml:"-"`
	fs        Filesys `yaml:"-"`

	Secrets    []Secret    `yaml:"secrets"`
	Workspaces []Workspace `yaml:"workspaces"`
}

// Save writes the current config object to the config file.
func (c *Config) Save() error {
	filePath := path.Join(c.configDir, configFilename)
	bytes, _ := yaml.Marshal(c)
	err := c.fs.Write(filePath, bytes)
	if err != nil {
		return err
	}

	return nil
}

// EnsureConfigFile ensures a config file exists in the configDir.
// If a configuration file does not already exist, an empty on
// will be created.
func EnsureConfigFile(configDir string, fs Filesys) error {
	filePath := path.Join(configDir, configFilename)
	exists, err := fs.FileExists(filePath)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	cnf := Config{
		Secrets:    make([]Secret, 0),
		Workspaces: make([]Workspace, 0),
	}

	bytes, _ := yaml.Marshal(cnf)
	err = fs.Write(filePath, bytes)
	if err != nil {
		return err
	}

	return nil
}

// LoadConfig loads a configuration file from configDir. An
// error will be returned if one does not exist.
func LoadConfig(configDir string, fs Filesys) (*Config, error) {
	filePath := path.Join(configDir, configFilename)
	bytes, err := fs.Read(filePath)
	if err != nil {
		return nil, err
	}

	var c Config
	c.configDir = configDir
	c.fs = fs
	_ = yaml.Unmarshal(bytes, &c)

	return &c, nil
}

// Secret is a struct which represents a stored secret value.
type Secret struct {
	Name   string `yaml:"name"`
	Value  string `yaml:"value"`
	Secure bool   `yaml:"secure"`
}

// GetValue returns the secret's value in plain text. If the is
// encrypted, it will be decrypted before being returned.
func (s *Secret) GetValue(cp CryptoProvider) string {
	if !s.Secure {
		return s.Value
	}

	v, err := cp.DecryptString(s.Value)
	if err != nil {
		return ""
	}

	return v
}

var (
	ErrSecretNameEmpty     = errors.New("secret: name cannot be empty")
	ErrSecretValueEmpty    = errors.New("secret: value cannot be empty")
	ErrSecretNotFound      = errors.New("secret: not found")
	ErrSecretAlreadyExists = errors.New("secret: already exists")
)

// AddSecret saves a new secret in the configuration file, with the
// given name and value. If the encrypt flag is true, the value
// will be encrypted before added to the config.
func (c *Config) AddSecret(name, value string, encrypt bool, cp CryptoProvider) error {
	if name == "" {
		return ErrSecretNameEmpty
	}

	if value == "" {
		return ErrSecretValueEmpty
	}

	_, err := c.GetSecret(name)
	if err == nil {
		return ErrSecretAlreadyExists
	}

	if encrypt {
		secureString, err := cp.EncryptString(value)
		if err != nil {
			return err
		}

		value = secureString
	}

	c.Secrets = append(c.Secrets, Secret{
		Name:   name,
		Value:  value,
		Secure: encrypt,
	})

	return nil
}

// GetSecret returns a secret from the config, where the name is equal
// to name. If the secret does not exist, ErrSecretNotFound will be returned.
func (c *Config) GetSecret(name string) (*Secret, error) {
	if name == "" {
		return nil, ErrSecretNameEmpty
	}

	for _, secret := range c.Secrets {
		if secret.Name == name {
			return &secret, nil
		}
	}

	return nil, ErrSecretNotFound
}

// RemoveSecret removes the secret with the given name from the config.
// If the secret does not exist, ErrSecretNotFound will be returned.
func (c *Config) RemoveSecret(name string) error {
	if name == "" {
		return ErrSecretNameEmpty
	}

	for i := 0; i < len(c.Secrets); i++ {
		if c.Secrets[i].Name == name {
			c.Secrets = append(c.Secrets[:i], c.Secrets[i+1:]...)

			return nil
		}
	}

	return ErrSecretNotFound
}

// Common workspace errors.
var (
	ErrWorkspaceNameEmpty  = errors.New("workspace: name is empty")
	ErrWorkspaceNameExists = errors.New("workspace: name already exists")
	ErrWorkspacePathEmpty  = errors.New("workspace: path is empty")
	ErrWorkspacePathExists = errors.New("workspace: path already exists")
	ErrWorkspaceNotFound   = errors.New("workspace: not found")

	ErrWorkspaceScriptNameEmpty    = errors.New("script: name is empty")
	ErrWorkspaceScriptNameExists   = errors.New("script: name already exists")
	ErrWorkspaceScriptCommandEmpty = errors.New("script: command is empty")
	ErrWorkspaceScriptNotFound     = errors.New("script: not found")
)

// Workspace is a struct which represents a workspace. A workspace
// contains a number of scripts which can be run in a given directory.
type Workspace struct {
	Name    string            `yaml:"name"`
	Path    string            `yaml:"path"`
	Scripts []WorkspaceScript `yaml:"scripts"`
}

// WorkspaceScript represents a script which can be run within a workspace.
type WorkspaceScript struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

// AddWorkspace is a function used to add a new workspace. This creates
// a new workspace instance and adds it to the config, c.
func (c *Config) AddWorkspace(name, path string) error {
	if name == "" {
		return ErrWorkspaceNameEmpty
	}

	if path == "" {
		return ErrWorkspacePathEmpty
	}

	for _, w := range c.Workspaces {
		if w.Name == name {
			return ErrWorkspaceNameExists
		}

		if w.Path == path {
			return ErrWorkspacePathExists
		}
	}

	c.Workspaces = append(c.Workspaces, Workspace{
		Name:    name,
		Path:    path,
		Scripts: make([]WorkspaceScript, 0),
	})

	return nil
}

// GetWorkspace retrieves a workspace from config, with a matching path.
func (c *Config) GetWorkspace(path string) (*Workspace, error) {
	if path == "" {
		return nil, ErrWorkspacePathEmpty
	}

	for _, w := range c.Workspaces {
		if w.Path == path {
			return &w, nil
		}
	}

	return nil, ErrWorkspaceNotFound
}

// AddScript is used to add a new script to a workspace.
func (w *Workspace) AddScript(name, command string) error {
	if name == "" {
		return ErrWorkspaceScriptNameEmpty
	}

	if command == "" {
		return ErrWorkspaceScriptCommandEmpty
	}

	for _, s := range w.Scripts {
		if s.Name == name {
			return ErrWorkspaceScriptNameExists
		}
	}

	w.Scripts = append(w.Scripts, WorkspaceScript{
		Name:    name,
		Command: command,
	})

	return nil
}

// GetScript retrives a script from w, with the given name.
func (w *Workspace) GetScript(name string) (*WorkspaceScript, error) {
	if name == "" {
		return nil, ErrWorkspaceScriptNameEmpty
	}

	for _, s := range w.Scripts {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, ErrWorkspaceScriptNotFound
}

// Run executes the workplace script.
func (s *WorkspaceScript) Run() (int, error) {
	args := strings.Split(s.Command, " ")
	c := exec.Command(args[0], args[1:]...)

	r, _ := c.StdoutPipe()
	err := c.Start()
	if err != nil {
		return -1, err
	}

	go func() {
		buf := make([]byte, 128)

		for {
			n, _ := r.Read(buf)
			os.Stdout.Write(buf[:n])
		}
	}()

	state, _ := c.Process.Wait()
	return state.ExitCode(), nil
}
