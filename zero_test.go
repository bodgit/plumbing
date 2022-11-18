package plumbing_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/bodgit/plumbing"
	"github.com/stretchr/testify/assert"
)

const limit = 10

func TestDevZero(t *testing.T) {
	t.Parallel()

	rw := plumbing.DevZero()
	b := new(bytes.Buffer)

	n, err := io.Copy(b, io.LimitReader(rw, limit))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, limit, int(n))
	assert.Equal(t, limit, b.Len())
	assert.Equal(t, bytes.Repeat([]byte{0x00}, limit), b.Bytes())

	n, err = io.Copy(rw, b)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, limit, int(n))
	assert.Equal(t, 0, b.Len())
}
