package age

import (
	"fmt"
	"strings"

	"github.com/murtaza-u/z/age/agelib"

	"filippo.io/age"
	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
	"github.com/rwxrob/help"
)

var symmetricCmd = &Z.Cmd{
	Name:    `symmetric`,
	Summary: `symmetric encryption/decryption`,
	Usage:   `(encrypt|decrypt) file [--out file]`,
	Commands: []*Z.Cmd{
		help.Cmd, symmetricEncryptCmd, symmetricDecryptCmd,
	},
}

var symmetricEncryptCmd = &Z.Cmd{
	Name:     `encrypt`,
	Summary:  `symmetric encryption`,
	Comp:     newComp(),
	Commands: []*Z.Cmd{help.Cmd},
	Usage:    `file [--out file]`,
	NumArgs:  1,
	Keys: Z.Keys{
		{
			Name:  `out`,
			Usage: `write result to a file`,
			Comp:  compfile.New(),
		},
	},
	MinArgs: 1,
	Call: func(caller *Z.Cmd, args ...string) error {
		in, err := agelib.ReadIn(args[0])
		if err != nil {
			return err
		}

		_outF := caller.GetVal("out")
		out, err := agelib.OpenOut(_outF)
		if err != nil {
			return err
		}
		defer out.Close()

		pswd := agelib.ReadHidden("password: ")
		_pswd := agelib.ReadHidden("confirm password: ")
		if _pswd != pswd {
			return fmt.Errorf("passwords do not match")
		}

		r, err := age.NewScryptRecipient(pswd)
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

var symmetricDecryptCmd = &Z.Cmd{
	Name:     `decrypt`,
	Summary:  `symmetric decryption`,
	Comp:     newComp(),
	Commands: []*Z.Cmd{help.Cmd},
	Usage:    `file [--out file]`,
	NumArgs:  1,
	Keys: Z.Keys{
		{
			Name:  `out`,
			Usage: `write result to a file`,
			Comp:  compfile.New(),
		},
	},
	Call: func(caller *Z.Cmd, args ...string) error {
		in, err := agelib.ReadIn(args[0])
		if err != nil {
			return err
		}

		_outF := caller.GetVal("out")
		out, err := agelib.OpenOut(_outF)
		if err != nil {
			return err
		}
		defer out.Close()

		pswd := agelib.ReadHidden("password: ")

		i, err := age.NewScryptIdentity(pswd)
		if err != nil {
			return err
		}

		return agelib.Decrypt(in, out, i)
	},
}
