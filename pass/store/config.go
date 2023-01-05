package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/murtaza-u/z/age/agelib"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Pass    pass `yaml:"pass"`
	subPath string
}

type pass struct {
	Store string   `yaml:"store"`
	Pubs  []string `yaml:"pubs"`
	Keys  []string `yaml:"keys"`
	Armor bool     `yaml:"armor"`
}

func NewConfig(data []byte, subpath string) (*Config, error) {
	c := new(Config)
	err := yaml.Unmarshal([]byte(data), c)
	if err != nil {
		return nil, err
	}

	c.subPath = subpath

	err = c.validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) validate() error {
	if c.Pass.Store == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		c.Pass.Store = filepath.Join(home, ".agepass")
	}

	_path := c.Pass.Store
	if c.subPath != "" {
		_path = filepath.Join(c.Pass.Store, c.subPath)
	}

	os.MkdirAll(_path, 0700)

	if len(c.Pass.Pubs) == 0 {
		return fmt.Errorf(
			".pass.pubs (public keys) not set in config",
		)
	}

	if len(c.Pass.Keys) == 0 {
		return fmt.Errorf(
			".pass.keys (private keys) not set in config",
		)
	}

	_, err := agelib.ParseRecipients(c.Pass.Pubs...)
	if err != nil {
		return fmt.Errorf("failed to parse public keys: %w", err)
	}

	return nil
}
