package store

import (
	"bytes"
	"fmt"

	"github.com/murtaza-u/z/age/agelib"
)

func (s *S) Insert(entry string) ([]byte, error) {
	if entry == "" {
		return nil, fmt.Errorf("missing entry name")
	}

	recs, err := agelib.ParseRecipients(s.C.Pass.Pubs...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public keys: %w", err)
	}

	if s.EntryExists(entry) {
		y := confirm(fmt.Sprintf(
			"entry %q exists. Overwrite?", entry,
		))
		if !y {
			return nil, nil
		}
	}

	pswd, err := s.InputSecret()
	if err != nil {
		return nil, err
	}

	in := bytes.NewReader([]byte(pswd))
	out := new(bytes.Buffer)

	err = agelib.Encrypt(in, out, s.C.Pass.Armor, recs...)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
