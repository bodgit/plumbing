package plumbing

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeeReaderAt(t *testing.T) {
	src := []byte("abcdef")
	dst := make([]byte, 3)
	rb := bytes.NewReader(src)
	wb := new(bytes.Buffer)

	r := TeeReaderAt(rb, wb)

	n, err := r.ReadAt(dst, 2)
	assert.Equal(t, 3, n)
	assert.Nil(t, err)
	assert.Equal(t, []byte("cde"), dst)
}
