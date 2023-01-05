package pass

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/murtaza-u/z/age/agelib"
	"github.com/murtaza-u/z/pass/store"

	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
)

var reencryptCmd = &Z.Cmd{
	Name:    `reencrypt`,
	Summary: `re-encrypt all secrets for different recipients`,
	Usage:   `r1 r2 r3...`,
	Comp:    compfile.New(),
	MinArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		recs, err := agelib.ParseRecipients(args...)
		if err != nil {
			return fmt.Errorf("failed to parse recipients: %w", err)
		}

		d, err := Z.Conf.Data()
		if err != nil {
			return err
		}
		c, err := store.NewConfig([]byte(d), "")
		if err != nil {
			return err
		}

		ids, err := agelib.ParseIdentities(c.Pass.Keys...)
		if err != nil {
			return fmt.Errorf("failed to parse old identities: %w", err)
		}

		var files []string

		err = filepath.WalkDir(c.Pass.Store, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			files = append(files, path)

			return nil
		})
		if err != nil {
			return err
		}

		pstore := c.Pass.Store + ".new"
		tstore := filepath.Join(pstore, "totp")
		fmt.Printf(
			"re-encrypted entries will be stored under %s\n", pstore,
		)

		err = os.MkdirAll(tstore, 0700)
		if err != nil {
			return fmt.Errorf("failed to create %q: %w", pstore, err)
		}

		for _, f := range files {
			// decrypt entry
			in, err := os.Open(f)
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
			dir := filepath.Dir(f)
			if strings.HasSuffix(dir, "totp") {
				dir = tstore
			} else {
				dir = pstore
			}

			noutF := filepath.Join(dir, filepath.Base(f))
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
