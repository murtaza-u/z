package age

import (
	"fmt"
	"strings"

	"github.com/murtaza-u/z/age/agelib"
	"github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
	"github.com/rwxrob/help"
)

var asymmetricCmd = &Z.Cmd{
	Name:    `asymmetric`,
	Summary: `asymmetric encryption/decryption`,
	Usage:   `(encrypt|decrypt) file [--out file]`,
	Commands: []*Z.Cmd{
		help.Cmd, asymmetricEncryptCmd, asymmetricDecryptCmd,
	},
}

var asymmetricEncryptCmd = &Z.Cmd{
	Name:     `encrypt`,
	Summary:  `asymmetric encryption`,
	Usage:    `file --recipients r1,r2,... [--out file]`,
	Comp:     newComp(),
	Commands: []*Z.Cmd{help.Cmd},
	Keys: Z.Keys{
		{
			Name:  `out`,
			Usage: `write result to a file`,
			Comp:  compfile.New(),
		},
		{
			Name:  `recipients`,
			Usage: `comma seperated recipients / recipient's file`,
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

		_recs := caller.GetVal("recipients")
		if _recs == "" {
			return fmt.Errorf("missing recipients")
		}

		recs, err := agelib.GetRecipients(_recs)
		if err != nil {
			return err
		}

		var armor bool
		if strings.HasSuffix(_outF, ".asc") {
			armor = true
		}

		return agelib.Encrypt(in, out, armor, recs...)
	},
}

var asymmetricDecryptCmd = &Z.Cmd{
	Name:     `decrypt`,
	Summary:  `asymmetric decryption`,
	Usage:    `file --identity privkeyfile [--out file]`,
	Comp:     newComp(),
	Commands: []*Z.Cmd{help.Cmd},
	Keys: Z.Keys{
		{
			Name:  `out`,
			Usage: `write result to a file`,
			Comp:  compfile.New(),
		},
		{
			Name:  `identity`,
			Usage: `path to identity file`,
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

		_id := caller.GetVal("identity")
		if _id == "" {
			return fmt.Errorf("missing identity file")
		}

		id, err := agelib.ParseIdentityFile(_id)

		return agelib.Decrypt(in, out, id...)
	},
}
