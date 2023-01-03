package pass

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/z"
)

var deleteCmd = &Z.Cmd{
	Name:    `delete`,
	Aliases: []string{"rm"},
	Summary: `delete an entry`,
	Usage:   `entry`,
	MinArgs: 1,
	Comp:    newComp(),
	Call: func(caller *Z.Cmd, args ...string) error {
		for _, entry := range args {
			exists, err := entryExists(entry)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("entry %q does not exist", entry)
			}

			y := confirm(fmt.Sprintf("delete %q?", entry))
			if !y {
				fmt.Printf("skipping %q\n", entry)
				continue
			}

			path := filepath.Join(Store, entry)
			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to delete %q: %w", path, err)
			}
		}

		return nil
	},
}
