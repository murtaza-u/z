package age

import (
	"log"
	"os"
	"path/filepath"

	"github.com/murtaza-u/conf"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "age",
	Usage: "encryption tool",
	Description: `Age is a simple, modern and secure encryption tool
with small explicit keys.`,
	Subcommands: []*cli.Command{keygenCmd, symmetricCmd, asymmetricCmd},
}

var Store string

func init() {
	conf := conf.New()
	conf.Init()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	_store, err := conf.Query(".age.store")
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
}
