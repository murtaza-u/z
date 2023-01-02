package agelib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"filippo.io/age/armor"
	"golang.org/x/crypto/ssh"
)

const (
	PrivateKeySizeLimit    = 1 << 24  // 16 MiB
	RecipientFileSizeLimit = 16 << 20 // 16 MiB
	LineLengthLimit        = 8 << 10  // 8 KiB, same as sshd(8)
)

func ParseIdentityFile(name string) ([]age.Identity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	b := bufio.NewReader(f)

	// length of "age-encryption" and "-----BEGIN AGE"
	p, _ := b.Peek(14)
	peeked := string(p)

	// an age encrypted file, plain or armored.
	if peeked == "age-encryption" || peeked == "-----BEGIN AGE" {
		var r io.Reader = b
		if peeked == "-----BEGIN AGE" {
			r = armor.NewReader(r)
		}

		contents, err := io.ReadAll(io.LimitReader(r, PrivateKeySizeLimit))
		if err != nil {
			return nil, fmt.Errorf("failed to read %q: %v", name, err)
		}

		if len(contents) == PrivateKeySizeLimit {
			return nil, fmt.Errorf("failed to read %q: file too long", name)
		}

		return []age.Identity{&EncryptedIdentity{
			Contents: contents,
			Passphrase: func() (string, error) {
				pswd := ReadHidden("passphrase for identity file %q: ", name)
				return pswd, nil
			},
			NoMatchWarning: func() {
				fmt.Fprintf(
					os.Stderr,
					"encrypted identity file %q didn't match file's recipients\n",
					name,
				)
			},
		}}, nil
	}

	// another PEM file, possibly an SSH private key.
	if strings.HasPrefix(peeked, "-----BEGIN") {
		contents, err := io.ReadAll(io.LimitReader(b, PrivateKeySizeLimit))
		if err != nil {
			return nil, fmt.Errorf("failed to read %q: %v", name, err)
		}

		if len(contents) == PrivateKeySizeLimit {
			return nil, fmt.Errorf("failed to read %q: file too long", name)
		}

		return parseSSHIdentity(name, contents)
	}

	ids, err := parseIdentities(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %v", name, err)
	}

	return ids, nil
}

type EncryptedIdentity struct {
	Contents       []byte
	Passphrase     func() (string, error)
	NoMatchWarning func()

	identities []age.Identity
}

func (i *EncryptedIdentity) Unwrap(stanzas []*age.Stanza) (fileKey []byte, err error) {
	if i.identities == nil {
		if err := i.decrypt(); err != nil {
			return nil, err
		}
	}

	for _, id := range i.identities {
		fileKey, err = id.Unwrap(stanzas)
		if errors.Is(err, age.ErrIncorrectIdentity) {
			continue
		}

		if err != nil {
			return nil, err
		}

		return fileKey, nil
	}
	i.NoMatchWarning()

	return nil, age.ErrIncorrectIdentity
}

func (i *EncryptedIdentity) decrypt() error {
	pswd, err := i.Passphrase()
	if err != nil {
		return err
	}

	id, err := age.NewScryptIdentity(pswd)
	if err != nil {
		return err
	}

	r, err := age.Decrypt(dearmor(bytes.NewReader(i.Contents)), id)
	if err != nil {
		return err
	}

	i.identities, err = parseIdentities(r)

	return err
}

func parseIdentity(s string) (age.Identity, error) {
	if strings.HasPrefix(s, "AGE-SECRET-KEY-1") {
		return age.ParseX25519Identity(s)
	}

	return nil, fmt.Errorf("unknown identity type")
}

func parseIdentities(f io.Reader) ([]age.Identity, error) {
	s := bufio.NewScanner(io.LimitReader(f, PrivateKeySizeLimit))

	var ids []age.Identity
	var n int

	for s.Scan() {
		n++
		line := s.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		i, err := parseIdentity(line)
		if err != nil {
			return nil, fmt.Errorf("error at line %d: %v", n, err)
		}
		ids = append(ids, i)

	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("failed to read secret keys file: %v", err)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no secret keys found")
	}

	return ids, nil
}

func parseSSHIdentity(name string, pemBytes []byte) ([]age.Identity, error) {
	id, err := agessh.ParseIdentity(pemBytes)
	if sshErr, ok := err.(*ssh.PassphraseMissingError); ok {
		pubK := sshErr.PublicKey
		if pubK == nil {
			pubK, err = readPubFile(name)
			if err != nil {
				return nil, err
			}
		}

		prompt := func() ([]byte, error) {
			pswd := ReadHidden("passphrase for %q: ", name)
			return []byte(pswd), nil
		}

		i, err := agessh.NewEncryptedSSHIdentity(pubK, pemBytes, prompt)
		if err != nil {
			return nil, err
		}

		return []age.Identity{i}, nil
	}

	if err != nil {
		return nil, fmt.Errorf(
			"malformed SSH identity in %q: %v", name, err,
		)
	}

	return []age.Identity{id}, nil
}
