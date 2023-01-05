package store

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/murtaza-u/z/age/agelib"
)

func (s S) Decrypt(entry string) ([]byte, error) {
	if entry == "" {
		return nil, fmt.Errorf("missing entry")
	}

	if !s.EntryExists(entry) {
		return nil, fmt.Errorf("entry %q does not exist", entry)
	}

	path := s.GetPath()
	path = filepath.Join(path, entry)

	in, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", path, err)
	}
	defer in.Close()

	id, err := agelib.ParseIdentities(s.C.Pass.Keys...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse identities: %w", err)
	}

	out := new(bytes.Buffer)

	err = agelib.Decrypt(in, out, id...)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
