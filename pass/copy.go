package pass

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/murtaza-u/z/age/agelib"
	"github.com/rwxrob/bonzai/z"

	"golang.design/x/clipboard"
)

const warning = "password copied to clipboard and will be cleared in 30s. Do *not* exit."

var copyCmd = &Z.Cmd{
	Name:    `copy`,
	Aliases: []string{"cp"},
	Summary: `decrypts and copies the password to clipboard`,
	Usage:   `entry`,
	NumArgs: 1,
	Comp:    newComp(),
	Call: func(caller *Z.Cmd, args ...string) error {
		arg := args[0]
		if arg == "" {
			return fmt.Errorf("missing entry")
		}

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

		id, err := agelib.ParseIdentities(key)
		if err != nil {
			return fmt.Errorf("failed to parse key %q: %w", key, err)
		}

		out := new(bytes.Buffer)

		err = agelib.Decrypt(in, out, id...)
		if err != nil {
			return err
		}

		err = clipboard.Init()
		if err != nil {
			return fmt.Errorf("failed to initialise clipboars: %q", err)
		}

		changed := clipboard.Write(clipboard.FmtText, out.Bytes())
		fmt.Println(warning)

		t := time.NewTimer(time.Second * 30)

		select {
		case <-changed:
		case <-t.C:
			clipboard.Write(clipboard.FmtText, []byte{})
		}

		return nil
	},
}
