package plumbing_test

import (
	"io"
	"strings"
	"testing"

	"github.com/bodgit/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestLimitReadCloser(t *testing.T) {
	t.Parallel()

	tables := []struct {
		name   string
		reader io.ReadCloser
		limit  int64
		n      int64
		err    error
	}{
		{
			"success",
			io.NopCloser(strings.NewReader("abcdefghij")),
			5,
			5,
			nil,
		},
		{
			"partial",
			io.NopCloser(strings.NewReader("abcde")),
			10,
			5,
			nil,
		},
	}

	for _, table := range tables {
		table := table
		t.Run(table.name, func(t *testing.T) {
			t.Parallel()
			r := plumbing.LimitReadCloser(table.reader, table.limit)
			defer r.Close()
			n, err := io.Copy(io.Discard, r)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.err, err)
		})
	}
}
