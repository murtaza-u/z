package util_test

import (
	"testing"

	"github.com/murtaza-u/z/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestEscReturns(t *testing.T) {
	in := "hello\nworld\ragain"
	exp := "hello\\nworld\\ragain"
	out := util.EscReturns(in)
	assert.Equal(t, exp, out)
}

func TestUnEscReturns(t *testing.T) {
	in := "hello\\nworld\\ragain"
	exp := "hello\nworld\ragain"
	out := util.UnEscReturns(in)
	assert.Equal(t, exp, out)
}
