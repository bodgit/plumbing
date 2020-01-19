package plumbing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteCounter(t *testing.T) {
	w := &WriteCounter{}

	n, err := w.Write([]byte("abcd"))
	assert.Equal(t, 4, n)
	assert.Nil(t, err)
	assert.Equal(t, uint64(4), w.Count())

	n, err = w.Write([]byte("efgh"))
	assert.Equal(t, 4, n)
	assert.Nil(t, err)
	assert.Equal(t, uint64(8), w.Count())
}
