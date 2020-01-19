package plumbing

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterleaveStripeReader(t *testing.T) {
	tables := []struct {
		stripe  int
		readers []io.Reader
		buffer  int
		wlen    int
		want    []byte
		eerr    error
	}{
		{
			1,
			[]io.Reader{
				strings.NewReader("aaaa"),
				strings.NewReader("bbbb"),
			},
			8,
			8,
			[]byte("abababab"),
			nil,
		},
		{
			2,
			[]io.Reader{
				strings.NewReader("aaaa"),
				strings.NewReader("bbbb"),
			},
			8,
			8,
			[]byte("aabbaabb"),
			nil,
		},
	}

	for _, table := range tables {
		got := make([]byte, table.buffer)
		ir := InterleaveStripeReader(table.stripe, table.readers...)
		glen, gerr := ir.Read(got[:])
		assert.Equal(t, table.eerr, gerr)
		assert.Equal(t, table.wlen, glen)
		assert.Equal(t, table.want, got)
	}
}
