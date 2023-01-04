package pass

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/murtaza-u/z/age/agelib"

	"github.com/rwxrob/bonzai/z"
)

var insertCmd = &Z.Cmd{
	Name:    `insert`,
	Summary: `insert a new password entry`,
	NumArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		entry := args[0]
		if entry == "" {
			return fmt.Errorf("missing entry name")
		}

		c, err := newCfg()
		if err != nil {
			return err
		}

		recs, err := agelib.ParseRecipients(c.Pass.Pubs...)
		if err != nil {
			return fmt.Errorf("failed to parse public keys: %w", err)
		}

		exists, err := entryExists(entry)
		if err != nil {
			return err
		}

		if exists {
			y := confirm(fmt.Sprintf(
				"entry %q exists. Overwrite?", entry,
			))
			if !y {
				return nil
			}
		}

		pswd := agelib.ReadHidden("password: ")
		_cpswd := agelib.ReadHidden("confirm password: ")
		if pswd != _cpswd {
			return fmt.Errorf("password do not match")
		}

		in := bytes.NewReader([]byte(pswd))
		out := new(bytes.Buffer)
		err = agelib.Encrypt(in, out, c.Pass.Armor, recs...)
		if err != nil {
			return err
		}

		_outF := filepath.Join(Store, entry)
		return os.WriteFile(_outF, out.Bytes(), 0600)
	},
}
