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

	"github.com/murtaza-u/conf"
	"github.com/urfave/cli/v2"
)

var reEncryptCmd = &cli.Command{
	Name:      "re-encrypt",
	Usage:     "re-encrypt all secrets for different recipients",
	UsageText: "r1 r2 r3...",
	Action: func(ctx *cli.Context) error {
		conf := conf.New()
		conf.MustInit()

		d, err := conf.Data()
		if err != nil {
			return err
		}
		c, err := store.NewConfig([]byte(d))
		if err != nil {
			return err
		}

		recs, err := agelib.ParseRecipients(ctx.Args().Slice()...)
		if err != nil {
			return fmt.Errorf("failed to parse recipients: %w", err)
		}

		ids, err := agelib.ParseIdentities(c.Pass.Keys...)
		if err != nil {
			return fmt.Errorf("failed to parse old identities: %w", err)
		}

		var files []string
		err = filepath.WalkDir(c.Pass.Store,
			func(path string, d fs.DirEntry, err error) error {
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

		newStore := c.Pass.Store + ".new"
		fmt.Printf(
			"re-encrypted entries will be stored under %s\n", newStore,
		)

		err = os.MkdirAll(newStore, 0700)
		if err != nil {
			return fmt.Errorf(
				"failed to create directory %q: %w", newStore, err,
			)
		}

		for _, f := range files {
			// decrypt file
			in, err := agelib.ReadIn(f)
			if err != nil {
				return err
			}
			out := new(bytes.Buffer)
			err = agelib.Decrypt(in, out, ids...)
			if err != nil {
				return fmt.Errorf("failed to decrypt %q: %w", f, err)
			}

			path := strings.TrimPrefix(f, c.Pass.Store)
			path = strings.TrimPrefix(path, string(filepath.Separator))
			path = filepath.Join(newStore, path)
			dir := filepath.Dir(path)
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				return err
			}

			newOut, err := agelib.OpenOut(path)
			if err != nil {
				return err
			}
			// re-encrypt for new recipient(s)
			err = agelib.Encrypt(out, newOut, false, recs...)
			if err != nil {
				return fmt.Errorf("failed to re-encrypt %q: %w", f, err)
			}
			newOut.Close()
		}

		return nil
	},
}
