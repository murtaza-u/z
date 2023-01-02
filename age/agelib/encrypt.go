package agelib

import (
	"fmt"
	"io"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func Encrypt(in io.Reader, out io.Writer, arm bool, recs ...age.Recipient) error {
	if len(recs) == 0 {
		return fmt.Errorf("missing recipients")
	}

	var _out io.Writer
	_out = out

	if arm {
		a := armor.NewWriter(out)
		defer a.Close()
		_out = a
	}

	w, err := age.Encrypt(_out, recs...)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, in)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
