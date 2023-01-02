package agelib

import (
	"fmt"
	"io"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func Encrypt(in io.Reader, out io.Writer, recs ...age.Recipient) error {
	if len(recs) == 0 {
		return fmt.Errorf("missing recipients")
	}

	a := armor.NewWriter(out)

	w, err := age.Encrypt(a, recs...)
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

	err = a.Close()
	if err != nil {
		return err
	}

	return nil
}
