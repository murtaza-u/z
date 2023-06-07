package age

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/murtaza-u/z/age/agelib"

	"filippo.io/age"
	"github.com/seehuhn/password"
	"github.com/urfave/cli/v2"
)

var keygenCmd = &cli.Command{
	Name:      "keygen",
	Usage:     "generates AGE public-private key pair",
	UsageText: "keygen NAME",
	Action: func(ctx *cli.Context) error {
		k, err := age.GenerateX25519Identity()
		if err != nil {
			return err
		}

		_outF := ctx.Args().First()
		if _outF == "" {
			return fmt.Errorf("missing output file")
		}
		_outF = filepath.Join(Store, _outF)

		defer func(k *age.X25519Identity) {
			pub := fmt.Sprintf("%s\n", k.Recipient())
			os.WriteFile(_outF+".pub", []byte(pub), 0600)
		}(k)

		pswd, err := password.Read("passphrase (optional): ")
		if err != nil {
			return fmt.Errorf("failed to read passphrase")
		}

		var _pswd []byte
		if pswd != nil && len(pswd) != 0 {
			_pswd, err = password.Read("confirm passphrase: ")
			if err != nil {
				return fmt.Errorf("failed to read passphrase")
			}
		}

		if bytes.Compare(pswd, _pswd) != 0 {
			return fmt.Errorf("passphrases do not match. Aborting!")
		}

		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "# created: %s\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(buf, "# public key: %s\n", k.Recipient())
		fmt.Fprintf(buf, "%s\n", k)

		if pswd == nil || len(pswd) == 0 {
			return os.WriteFile(_outF, buf.Bytes(), 0600)
		}

		r, err := age.NewScryptRecipient(string(pswd))
		if err != nil {
			return err
		}

		in := bytes.NewReader(buf.Bytes())
		out := new(bytes.Buffer)

		err = agelib.Encrypt(in, out, true, r)
		if err != nil {
			return fmt.Errorf(
				"failed to encrypt private key with passphrase: %w",
				err)
		}

		return os.WriteFile(_outF, out.Bytes(), 0600)
	},
}
