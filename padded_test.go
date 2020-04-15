package plumbing

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaddedReader(t *testing.T) {
	src := []byte("abcdef")
	rb := bytes.NewReader(src)
	wb := new(bytes.Buffer)

	r := PaddedReader(rb, 8, 0)

	n, err := io.Copy(wb, r)
	assert.Equal(t, int64(8), n)
	assert.Nil(t, err)
	assert.Equal(t, []byte{'a', 'b', 'c', 'd', 'e', 'f', 0, 0}, wb.Bytes())
}
