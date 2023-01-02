package age

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/z"
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
		help.Cmd, keygenCmd, symmetricCmd, asymmetricCmd,
	},
}

var defaultAgeDir, sshDir string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf(
			"failed to initialise default age directory: %s",
			err.Error(),
		)
	}

	defaultAgeDir = filepath.Join(home, ".age")
	sshDir = filepath.Join(home, ".ssh")
}
