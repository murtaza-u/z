/*
Package vars provides a high-level abstraction over an internal
concurrency-safe map object and operations such as Get, Set, Del, and Exists
for the management of a temporary local cache. It allows the preservation of
state between command executions.
*/
package vars

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/murtaza-u/z/internal/util"

	"github.com/rogpeppe/go-internal/lockedfile"
)

// V provides an internal concurrency-safe map object as well as
// different operations to be performed on it.
type V struct {
	mu sync.Mutex
	m  map[string]string
}

// New creates a new map type internally. It is mandatory to run the
// Init() method immediately after.
func New() *V {
	return &V{
		m: make(map[string]string),
	}
}

// Init creates the necessary cache directory and file (if absent) and
// loads the keys and values into the internal map object. It is
// mandatory to call this method before performing any operations.
func (v *V) Init() error {
	path, err := v.path()
	if err != nil {
		return fmt.Errorf("failed to get cache file path: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(path), 0o751)
	if err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", path, err)
	}
	f.Close()

	err = v.read()
	if err != nil {
		return err
	}

	return nil
}

// MustInit is same as Init but panics in case of an error.
func (v *V) MustInit() {
	err := v.Init()
	if err != nil {
		panic(err)
	}
}

// Get returns the value associated with the given key. If the key does
// not exist, it returns a blank string ("").
func (v *V) Get(key string) string {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v, ok := v.m[key]; ok {
		return v
	}

	return ""
}

// Exists checks if a key is present.
func (v *V) Exists(key string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	_, ok := v.m[key]
	return ok
}

// Del deletes a key. If the key does not exist, no operation is
// performed.
func (v *V) Del(key string) error {
	v.mu.Lock()
	delete(v.m, key)
	v.mu.Unlock()
	return v.write()
}

// Set adds/updates a key-value pair.
func (v *V) Set(key, value string) error {
	v.mu.Lock()
	v.m[key] = value
	v.mu.Unlock()
	return v.write()
}

func (*V) path() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	app := filepath.Base(os.Args[0])
	return filepath.Join(cache, app), nil
}

func (v *V) read() error {
	path, err := v.path()
	if err != nil {
		return err
	}

	data, err := lockedfile.Read(path)
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", path, err)
	}

	err = v.unmarshalText(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text: %w", err)
	}

	return nil
}

func (v *V) write() error {
	path, err := v.path()
	if err != nil {
		return err
	}

	out := v.marshalText()
	err = lockedfile.Write(path, bytes.NewReader(out), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write to %q: %w", path, err)
	}
	return nil
}

func (v *V) marshalText() []byte {
	v.mu.Lock()
	defer v.mu.Unlock()

	var out string
	for k, v := range v.m {
		out += fmt.Sprintf("%s=%s\n", k, util.EscReturns(v))
	}

	return []byte(out)
}

func (v *V) unmarshalText(in []byte) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	sc := bufio.NewScanner(bytes.NewReader(in))
	for sc.Scan() {
		txt := sc.Text()
		splits := strings.SplitN(txt, "=", 2)
		if len(splits) == 2 {
			key := splits[0]
			val := util.UnEscReturns(splits[1])
			v.m[key] = val
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("failed to parse cache file: %w", err)
	}

	return nil
}
