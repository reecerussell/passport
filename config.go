package passport

import (
	"errors"
	"path"

	"gopkg.in/yaml.v3"
)

const configFilename = "config.yaml"

// Config is a struct which holds and represents the core configuration.
type Config struct {
	configDir string  `yaml:"-"`
	fs        Filesys `yaml:"-"`

	Secrets []Secret `yaml:"secrets"`
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
		Secrets: make([]Secret, 0),
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

// Secret is a struct which represents a store secret value.
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
