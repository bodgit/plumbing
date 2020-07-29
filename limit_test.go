package plumbing

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitReadCloser(t *testing.T) {
	tables := map[string]struct {
		reader io.ReadCloser
		limit  int64
		n      int64
		err    error
	}{
		"success": {
			ioutil.NopCloser(strings.NewReader("abcdefghij")),
			5,
			5,
			nil,
		},
		"partial": {
			ioutil.NopCloser(strings.NewReader("abcde")),
			10,
			5,
			nil,
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			r := LimitReadCloser(table.reader, table.limit)
			defer r.Close()
			n, err := io.Copy(ioutil.Discard, r)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.err, err)
		})
	}
}
