package age

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/murtaza-u/z/age/agelib"

	"filippo.io/age"
	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
)

var keygenCmd = &Z.Cmd{
	Name:     `keygen`,
	Summary:  `generates AGE public-private key pair`,
	Commands: []*Z.Cmd{help.Cmd},
	NumArgs:  1,
	Usage:    `file`,
	Call: func(caller *Z.Cmd, args ...string) error {
		k, err := age.GenerateX25519Identity()
		if err != nil {
			return err
		}

		_outF := args[0]
		if _outF == "" {
			return fmt.Errorf("missing output file")
		}
		_outF = filepath.Join(Store, _outF)

		defer func(k *age.X25519Identity) {
			pub := fmt.Sprintf("%s\n", k.Recipient())
			os.WriteFile(_outF+".pub", []byte(pub), 0600)
		}(k)

		pswd := term.PromptHidden("passphrase (optional): ")
		fmt.Println()

		var _pswd string
		if pswd != "" {
			_pswd = term.PromptHidden("confirm passphrase: ")
			fmt.Println()
		}

		if pswd != _pswd {
			return fmt.Errorf("passphrases do not match. Aborting!")
		}

		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "# created: %s\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(buf, "# public key: %s\n", k.Recipient())
		fmt.Fprintf(buf, "%s\n", k)

		if pswd == "" {
			return os.WriteFile(_outF, buf.Bytes(), 0600)
		}

		r, err := age.NewScryptRecipient(pswd)
		if err != nil {
			return err
		}

		in := bytes.NewReader(buf.Bytes())
		out := new(bytes.Buffer)

		err = agelib.Encrypt(in, out, true, r)
		if err != nil {
			return fmt.Errorf(
				"failed to encrypt private key with passphrase: %w",
				err,
			)
		}

		return os.WriteFile(_outF, out.Bytes(), 0600)
	},
}
