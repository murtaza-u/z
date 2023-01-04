package pass

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/murtaza-u/z/age/agelib"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
)

var reencryptCmd = &Z.Cmd{
	Name:    `reencrypt`,
	Summary: `re-encrypt all passwords for different recipients`,
	Usage:   `new-recipient-file`,
	Comp:    compfile.New(),
	MinArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		recs, err := agelib.ParseRecipients(args...)
		if err != nil {
			return fmt.Errorf("failed to parse recipients: %w", err)
		}
		if len(recs) == 0 {
			return fmt.Errorf("no valid recipients found")
		}

		c, err := newCfg()
		if err != nil {
			return err
		}

		ids, err := agelib.ParseIdentities(c.Pass.Keys...)
		if err != nil {
			return fmt.Errorf("failed to parse keys: %w", err)
		}

		var files []string

		err = filepath.WalkDir(Store, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			files = append(files, d.Name())

			return nil
		})
		if err != nil {
			return err
		}

		_store := Store + ".new"
		fmt.Printf(
			"re-encrypted entries will be stored under %s\n", _store,
		)

		err = os.MkdirAll(_store, 0700)
		if err != nil {
			return fmt.Errorf(
				"failed to create %q directory: %w", _store, err,
			)
		}

		for _, f := range files {
			// decrypt entry
			in, err := os.Open(filepath.Join(Store, f))
			if err != nil {
				return fmt.Errorf("failed to open %q: %w", f, err)
			}
			out := new(bytes.Buffer)
			err = agelib.Decrypt(in, out, ids...)
			if err != nil {
				return fmt.Errorf("failed to decrypt %q: %w", f, err)
			}
			in.Close()

			// re-encrypt for new recipients
			noutF := filepath.Join(_store, f)
			nout, err := agelib.OpenOut(noutF)
			if err != nil {
				return err
			}

			err = agelib.Encrypt(out, nout, c.Pass.Armor, recs...)
			if err != nil {
				return fmt.Errorf("failed to re-encrypt %q: %w", f, err)
			}
		}

		return nil
	},
}
