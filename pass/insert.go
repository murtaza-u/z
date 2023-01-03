package pass

import (
	"bytes"
	"fmt"
	"path/filepath"

	"filippo.io/age"
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

		pub, err := Z.Conf.Query(`.pass.pub`)
		if err != nil {
			return err
		}

		var arm bool

		_arm, err := Z.Conf.Query(`.pass.armor`)
		if err != nil {
			return err
		}

		if _arm == "true" {
			arm = true
		}

		if pub == "null" {
			return fmt.Errorf(".pass.pub not set in config")
		}

		r, err := age.ParseX25519Recipient(pub)
		if err != nil {
			return fmt.Errorf("failed to parse public key: %w", err)
		}

		exists, err := entryExists(entry)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("entry %q exists", entry)
		}

		pswd := agelib.ReadHidden("password: ")
		_cpswd := agelib.ReadHidden("confirm password: ")
		if pswd != _cpswd {
			return fmt.Errorf("password do not match")
		}

		in := bytes.NewReader([]byte(pswd))

		_outF := filepath.Join(Store, entry)
		out, err := agelib.OpenOut(_outF)
		if err != nil {
			return err
		}
		defer out.Close()

		return agelib.Encrypt(in, out, arm, []age.Recipient{r}...)
	},
}
