package store

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/murtaza-u/z/age/agelib"
	"github.com/rwxrob/fs/file"
)

type Storer interface {
	GetPath() string
	List() []string
	ListFull() []string
	Delete(...string) error
	EntryExists(string) bool
	WriteEntry(string, []byte) error
	Insert(string) ([]byte, error)
	Decrypt(string) ([]byte, error)
}

type S struct {
	C           *Config
	InputSecret func() (string, error)
}

func New(c *Config) *S {
	return &S{C: c, InputSecret: func() (string, error) {
		pswd := agelib.ReadHidden("password: ")
		_cpswd := agelib.ReadHidden("confirm password: ")
		if pswd != _cpswd {
			return "", fmt.Errorf("passwords do not match")
		}

		return pswd, nil
	}}
}

func (s S) GetPath() string {
	if s.C.subPath == "" {
		return s.C.Pass.Store
	}

	return filepath.Join(s.C.Pass.Store, s.C.subPath)
}

func (s S) List() []string {
	path := s.GetPath()

	var list []string

	err := filepath.WalkDir(path, func(sub string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if sub == path {
			return nil
		}

		if filepath.Dir(sub) != path {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		list = append(list, d.Name())

		return nil
	})

	if err != nil {
		return []string{}
	}

	return list
}

func (s S) ListFull() []string {
	path := s.GetPath()

	var list []string

	err := filepath.WalkDir(path, func(sub string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if sub == path {
			return nil
		}

		if filepath.Dir(sub) != path {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		list = append(list, sub)

		return nil
	})

	if err != nil {
		return nil
	}

	return list
}

func (s S) Delete(entries ...string) error {
	for _, e := range entries {
		if !s.EntryExists(e) {
			return fmt.Errorf("entry %q does not exist", e)
		}

		y := confirm(fmt.Sprintf("delete %q?", e))
		if !y {
			fmt.Printf("skipping %q\n", e)
			continue
		}

		path := s.GetPath()
		path = filepath.Join(path, e)

		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("failed to delete %q: %w", path, err)
		}
	}

	return nil
}

func (s S) EntryExists(entry string) bool {
	path := s.GetPath()
	path = filepath.Join(path, entry)
	return file.Exists(path)
}

func (s S) WriteEntry(entry string, out []byte) error {
	path := s.GetPath()
	path = filepath.Join(path, entry)
	return os.WriteFile(path, out, 0600)
}
