package plumbing

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errClose = errors.New("error closing")

type errorWriteCloser struct {
	io.WriteCloser
}

func (errorWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (errorWriteCloser) Close() error {
	return errClose
}

type partialWriter struct {
	io.Writer
}

func (partialWriter) Write(p []byte) (n int, err error) {
	return len(p) - 1, nil
}

func TestMultiWriteCloser(t *testing.T) {
	tables := map[string]struct {
		writer   io.WriteCloser
		n        int
		writeErr error
		closeErr error
	}{
		"success": {
			NopWriteCloser(new(bytes.Buffer)),
			10,
			nil,
			nil,
		},
		"nested": {
			MultiWriteCloser(NopWriteCloser(new(bytes.Buffer))),
			10,
			nil,
			nil,
		},
		"write": {
			NopWriteCloser(errorWriter{}),
			0,
			errWrite,
			nil,
		},
		"close": {
			errorWriteCloser{},
			10,
			nil,
			errClose,
		},
		"partial": {
			NopWriteCloser(partialWriter{}),
			9,
			io.ErrShortWrite,
			nil,
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			dst := NopWriteCloser(new(bytes.Buffer))
			w := MultiWriteCloser(table.writer, dst)
			n, err := w.Write(in)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.writeErr, err)
			err = w.Close()
			assert.Equal(t, table.closeErr, err)
		})
	}
}
