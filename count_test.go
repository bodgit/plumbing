package plumbing_test

import (
	"testing"

	"github.com/bodgit/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestWriteCounter(t *testing.T) {
	t.Parallel()

	w := &plumbing.WriteCounter{}

	n, err := w.Write([]byte("abcd"))
	assert.Equal(t, 4, n)
	assert.Nil(t, err)
	assert.Equal(t, uint64(4), w.Count())

	n, err = w.Write([]byte("efgh"))
	assert.Equal(t, 4, n)
	assert.Nil(t, err)
	assert.Equal(t, uint64(8), w.Count())
}
