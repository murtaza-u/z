package agelib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
)

const PrivateKeySizeLimit = 1 << 24 // 16 MiB

func ParseIdentities(fnames ...string) ([]age.Identity, error) {
	buf := new(bytes.Buffer)

	var ids []age.Identity

	for _, name := range fnames {
		f, err := os.Open(name)
		if err != nil {
			return nil, fmt.Errorf("failed to open %q: %w", name, err)
		}
		r := bufio.NewReader(f)

		isenc, err := isEncrypted(r)
		if err != nil {
			return nil, err
		}

		if isenc {
			id, err := parseEncID(r, name)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)

			continue
		}

		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		data = append(data, []byte("\n")...)
		_, err = buf.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %w", err)
		}
	}

	if buf.Len() != 0 {
		_ids, err := age.ParseIdentities(buf)
		if err != nil {
			return nil, err
		}
		ids = append(ids, _ids...)
	}

	return ids, nil
}

func header(b *bufio.Reader) string {
	p, err := b.Peek(14)
	if err != nil {
		return ""
	}
	return string(p)
}

func isEncrypted(b *bufio.Reader) (bool, error) {
	h := header(b)
	if h == "age-encryption" || h == "-----BEGIN AGE" {
		return true, nil
	}

	return false, nil
}

func parseEncID(b *bufio.Reader, fname string) (age.Identity, error) {
	var r io.Reader = b

	h := header(b)
	if h == "-----BEGIN AGE" {
		r = armor.NewReader(r)
	}

	content, err := io.ReadAll(io.LimitReader(r, PrivateKeySizeLimit))
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", fname, err)
	}

	if len(content) == PrivateKeySizeLimit {
		return nil, fmt.Errorf("file %q is too long", fname)
	}

	return newEncID(fname, content), nil
}

type encID struct {
	content    []byte
	passphrase func() (string, error)
	warning    func()
	identities []age.Identity
}

func newEncID(fname string, content []byte) *encID {
	return &encID{
		content: content,
		passphrase: func() (string, error) {
			pswd := ReadHidden(
				"passphrase for identity file %q: ", fname,
			)
			return pswd, nil
		},
		warning: func() {
			fmt.Fprintf(
				os.Stderr,
				"encrypted identity file %q didn't match any recipients\n",
				fname,
			)
		},
	}
}

func (i *encID) Unwrap(stanzas []*age.Stanza) ([]byte, error) {
	if i.identities == nil {
		if err := i.decrypt(); err != nil {
			return nil, err
		}
	}

	for _, id := range i.identities {
		key, err := id.Unwrap(stanzas)
		if errors.Is(err, age.ErrIncorrectIdentity) {
			continue
		}

		if err != nil {
			return nil, err
		}

		return key, nil
	}
	i.warning()

	return nil, age.ErrIncorrectIdentity
}

var errIncorrectPass = errors.New("no identity matched any of the recipients")

func (i *encID) decrypt() error {
	pswd, err := i.passphrase()
	if err != nil {
		return err
	}

	id, err := age.NewScryptIdentity(pswd)
	if err != nil {
		return err
	}

	r, err := age.Decrypt(dearmor(bytes.NewReader(i.content)), id)
	if err != nil {
		if err.Error() == errIncorrectPass.Error() {
			return fmt.Errorf("incorrect passphrase")
		}

		return err
	}

	i.identities, err = age.ParseIdentities(r)

	return err
}
