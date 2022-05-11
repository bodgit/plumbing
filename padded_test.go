package plumbing_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/bodgit/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestPaddedReader(t *testing.T) {
	t.Parallel()

	src := []byte("abcdef")
	rb := bytes.NewReader(src)
	wb := new(bytes.Buffer)

	r := plumbing.PaddedReader(rb, 8, 0)

	n, err := io.Copy(wb, r)
	assert.Equal(t, int64(8), n)
	assert.Nil(t, err)
	assert.Equal(t, []byte{'a', 'b', 'c', 'd', 'e', 'f', 0, 0}, wb.Bytes())
}
