package age

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
)

var Cmd = &Z.Cmd{
	Name:    `age`,
	Summary: `encryption tool`,
	Description: `
		A simple, modern and secure encryption tool (and Go library)
		with small explicit keys, no config options, and UNIX-style
		composability.`,
	Site:    `https://age-encryption.org`,
	Source:  `https://github.com/FiloSottile/age`,
	Issues:  `https://github.com/FiloSottile/age/issues`,
	License: `BSD-3-Clause`,
	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, keygenCmd, symmetricCmd, asymmetricCmd,
	},
}

var Store, SSHDir string

func init() {
	Z.Conf.SoftInit()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	_store, err := Z.Conf.Query(".age.store")
	if err != nil {
		log.Fatal(err)
	}
	if _store == "null" {
		_store = filepath.Join(home, ".age")
	}
	Store = _store

	err = os.MkdirAll(Store, 0700)
	if err != nil {
		log.Fatal(err)
	}

	SSHDir = filepath.Join(home, ".ssh")
}
