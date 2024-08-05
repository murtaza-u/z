package vars_test

import (
	"testing"

	"github.com/murtaza-u/z/internal/vars"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	v := vars.New()
	err := v.Init()
	assert.Nil(t, err, "failed to initialize vars")
}

func TestSet(t *testing.T) {
	v := vars.New()
	err := v.Init()
	assert.Nil(t, err, "failed to initialize vars")

	err = v.Set("foo", "bar")
	assert.Nil(t, err, "failed to set foo=bar")
}

func TestExists(t *testing.T) {
	v := vars.New()
	err := v.Init()
	assert.Nil(t, err, "failed to initialize vars")

	ok := v.Exists("foo")
	assert.True(t, ok, "key `foo` exists, but false returned")

	ok = v.Exists("bar")
	assert.False(t, ok, "key `bar` does not exists, but true returned")
}

func TestGet(t *testing.T) {
	vars := vars.New()
	err := vars.Init()
	assert.Nil(t, err, "failed to initialize vars")

	v := vars.Get("foo")
	assert.Equal(t, "bar", v)

	v = vars.Get("bar")
	assert.Empty(t, v)
}

func TestDel(t *testing.T) {
	vars := vars.New()
	err := vars.Init()
	assert.Nil(t, err, "failed to initialize vars")

	err = vars.Del("foo")
	assert.Nil(t, err, "failed to delete key `foo`")

	v := vars.Get("foo")
	assert.Empty(t, v, "key `foo` not deleted")
}
