package age

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/murtaza-u/z/age/agelib"

	"filippo.io/age"
	"github.com/seehuhn/password"
	"github.com/urfave/cli/v2"
)

var symmetricCmd = &cli.Command{
	Name:      "symmetric",
	Usage:     "symmetric encryption/decryption",
	UsageText: "symmetric (encrypt|decrypt) FILE [--out file]",
	Subcommands: []*cli.Command{
		symmetricEncryptCmd, symmetricDecryptCmd,
	},
}

var symmetricEncryptCmd = &cli.Command{
	Name:      "encrypt",
	Usage:     "symmetric encryption",
	UsageText: "FILE [--out FILE]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "out",
			Aliases:   []string{"o"},
			Usage:     "write result to a `FILE`",
			TakesFile: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		arg := ctx.Args().First()
		if arg == "" {
			return fmt.Errorf("missing file to encrypt")
		}

		_outF := ctx.String("out")
		out, err := agelib.OpenOut(_outF)
		if err != nil {
			return err
		}
		defer out.Close()

		in, err := agelib.ReadIn(arg)
		if err != nil {
			return err
		}

		pswd, err := password.Read("password: ")
		if err != nil {
			return fmt.Errorf("unable to get password: %w", err)
		}

		_pswd, err := password.Read("confirm password: ")
		if err != nil {
			return fmt.Errorf("unable to get password: %w", err)
		}

		if bytes.Compare(pswd, _pswd) != 0 {
			return fmt.Errorf("passwords do not match")
		}

		r, err := age.NewScryptRecipient(string(pswd))
		if err != nil {
			return err
		}

		var armor bool
		if strings.HasSuffix(_outF, ".asc") {
			armor = true
		}

		return agelib.Encrypt(in, out, armor, r)
	},
}

var symmetricDecryptCmd = &cli.Command{
	Name:      "decrypt",
	Usage:     "symmetric decryption",
	UsageText: "FILE [--out FILE]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "out",
			Aliases:   []string{"o"},
			Usage:     "write result to a `FILE`",
			TakesFile: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		arg := ctx.Args().First()
		if arg == "" {
			return fmt.Errorf("missing file to decrypt")
		}

		in, err := agelib.ReadIn(arg)
		if err != nil {
			return err
		}

		_outF := ctx.String("out")
		out, err := agelib.OpenOut(_outF)
		if err != nil {
			return err
		}
		defer out.Close()

		pswd, err := password.Read("password: ")
		if err != nil {
			return fmt.Errorf("unable to get password: %w", err)
		}

		i, err := age.NewScryptIdentity(string(pswd))
		if err != nil {
			return err
		}

		return agelib.Decrypt(in, out, i)
	},
}
