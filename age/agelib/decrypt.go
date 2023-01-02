package agelib

import (
	"bufio"
	"fmt"
	"io"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func Decrypt(in io.Reader, out io.Writer, ids ...age.Identity) error {
	if len(ids) == 0 {
		return fmt.Errorf("missing identities")
	}

	r, err := age.Decrypt(dearmor(in), ids...)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return nil
}

func dearmor(in io.Reader) io.Reader {
	r := bufio.NewReader(in)
	start, _ := r.Peek(len(armor.Header))
	if string(start) == armor.Header {
		return armor.NewReader(r)
	}
	return r
}
