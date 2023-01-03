package pass

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/murtaza-u/z/age/agelib"
	Z "github.com/rwxrob/bonzai/z"
)

var showCmd = &Z.Cmd{
	Name:    `show`,
	Aliases: []string{"ls", "list"},
	Summary: `list entries / decrypt an entry`,
	Usage:   `[entry]`,
	MaxArgs: 1,
	Comp:    newComp(),
	Call: func(caller *Z.Cmd, args ...string) error {
		if len(args) == 0 || args[0] == "" {
			for _, e := range listEntries() {
				fmt.Println(e)
			}

			return nil
		}

		arg := args[0]
		path := filepath.Join(Store, arg)
		in, err := os.OpenFile(path, os.O_RDONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to open %q: %w", path, err)
		}
		defer in.Close()

		key, err := Z.Conf.Query(".pass.key")
		if err != nil {
			return err
		}

		if key == "null" {
			return fmt.Errorf(
				".pass.key (private key) not set in config",
			)
		}

		id, err := agelib.ParseIdentityFile(key)
		if err != nil {
			return fmt.Errorf("failed to parse key %q: %w", key, err)
		}

		err = agelib.Decrypt(in, os.Stdout, id...)
		if err != nil {
			return err
		}
		fmt.Println()

		return nil
	},
}
