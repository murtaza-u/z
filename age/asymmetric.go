package age

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/murtaza-u/z/age/agelib"

	"github.com/urfave/cli/v2"
)

var asymmetricCmd = &cli.Command{
	Name:      "asymmetric",
	Usage:     "asymmetric encryption/decryption",
	UsageText: "(encrypt|decrypt) FILE [--out FILE]",
	Subcommands: []*cli.Command{
		asymmetricEncryptCmd, asymmetricDecryptCmd,
	},
}

var asymmetricEncryptCmd = &cli.Command{
	Name:      "encrypt",
	Usage:     "asymmetric encryption",
	UsageText: "FILE --recipient r1,r2,... [--out FILE]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      `out`,
			Aliases:   []string{"o"},
			Usage:     "write result to a `FILE`",
			TakesFile: true,
		},
		&cli.StringFlag{
			Name:     "recipient",
			Aliases:  []string{"r"},
			Usage:    "comma-seperated list of recipient",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		arg := ctx.Args().First()
		if arg == "" {
			arg = os.Stdin.Name()
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

		_recs := ctx.String("recipient")
		if _recs == "" {
			return fmt.Errorf("missing recipients")
		}

		recsList := strings.Split(_recs, ",")
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf(
				"failed to determine user's home directory: %w", err)
		}
		for i := 0; i < len(recsList); i++ {
			if strings.HasPrefix(recsList[i], "~") {
				recsList[i] = strings.Replace(recsList[i], "~", home, 1)
			}
		}

		recs, err := agelib.ParseRecipients(recsList...)
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

var asymmetricDecryptCmd = &cli.Command{
	Name:      "decrypt",
	Usage:     "asymmetric decryption",
	UsageText: "FILE [--identity i1,i2,...] [--out FILE]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "out",
			Aliases:   []string{"o"},
			Usage:     "write result to a `FILE`",
			TakesFile: true,
		},
		&cli.StringFlag{
			Name:      "identity",
			Aliases:   []string{"i"},
			Usage:     "comma-seperated list of identity `FILE`",
			TakesFile: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		arg := ctx.Args().First()
		if arg == "" {
			arg = os.Stdin.Name()
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

		var files []string

		_id := ctx.String("identity")
		if _id != "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf(
					"failed to determine user's home directory: %w", err)
			}
			for _, f := range strings.Split(_id, ",") {
				if strings.HasPrefix(f, "~") {
					f = strings.Replace(f, "~", home, 1)
				}
				files = append(files, f)
			}
		}

		if len(files) == 0 {
			// lookup every file under age store
			filepath.WalkDir(Store, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() {
					return nil
				}

				if strings.HasSuffix(d.Name(), ".pub") {
					return nil
				}

				files = append(files, path)

				return nil
			})
		}

		ids, err := agelib.ParseIdentities(files...)
		if err != nil {
			return fmt.Errorf("failed to parse identity files: %w", err)
		}

		if len(ids) == 0 {
			return fmt.Errorf("no valid identity file found")
		}

		return agelib.Decrypt(in, out, ids...)
	},
}
