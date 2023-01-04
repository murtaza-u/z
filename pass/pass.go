package pass

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
)

var Cmd = &Z.Cmd{
	Name:    `pass`,
	Summary: `password manager based on AGE`,
	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, showCmd, checkCmd, insertCmd, deleteCmd,
		copyCmd, reencryptCmd,
	},
}

var Store string

func init() {
	Z.Conf.SoftInit()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	_store, err := Z.Conf.Query(".pass.store")
	if err != nil {
		log.Fatal(err)
	}
	if _store == "null" {
		_store = filepath.Join(home, ".agestore")
	}
	Store = _store

	err = os.MkdirAll(Store, 0700)
	if err != nil {
		log.Fatal(err)
	}
}
