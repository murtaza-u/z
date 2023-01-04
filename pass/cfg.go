package pass

import (
	"fmt"

	"github.com/murtaza-u/z/age/agelib"
	"github.com/rwxrob/bonzai/z"
	"gopkg.in/yaml.v3"
)

type pass struct {
	Store string   `yaml:"store"`
	Pubs  []string `yaml:"pubs"`
	Keys  []string `yaml:"keys"`
	Armor bool     `yaml:"armor"`
}

type cfg struct {
	Pass pass `yaml:"pass"`
}

func newCfg() (*cfg, error) {
	d, err := Z.Conf.Data()
	if err != nil {
		return nil, err
	}

	c := new(cfg)
	err = yaml.Unmarshal([]byte(d), c)
	if err != nil {
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *cfg) validate() error {
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
