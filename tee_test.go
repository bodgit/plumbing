package plumbing

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errWrite = errors.New("error writing")

type errorWriter struct {
	io.Writer
}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 0, errWrite
}

func TestTeeReaderAt(t *testing.T) {
	in := []byte("abcdefghij")

	tables := map[string]struct {
		reader io.ReaderAt
		writer io.Writer
		length int
		offset int64
		n      int
		err    error
	}{
		"success": {
			bytes.NewReader(in),
			ioutil.Discard,
			3,
			2,
			3,
			nil,
		},
		"fail": {
			bytes.NewReader(in),
			errorWriter{},
			3,
			2,
			0,
			errWrite,
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			r := TeeReaderAt(table.reader, table.writer)
			dst := make([]byte, table.length)
			n, err := r.ReadAt(dst, table.offset)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.err, err)
		})
	}
}

func TestTeeReadCloser(t *testing.T) {
	in := []byte("abcdefghij")

	tables := map[string]struct {
		reader io.ReadCloser
		writer io.Writer
		n      int64
		err    error
	}{
		"success": {
			ioutil.NopCloser(bytes.NewReader(in)),
			ioutil.Discard,
			10,
			nil,
		},
		"fail": {
			ioutil.NopCloser(bytes.NewReader(in)),
			errorWriter{},
			0,
			errWrite,
		},
	}

	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			r := TeeReadCloser(table.reader, table.writer)
			defer r.Close()
			n, err := io.Copy(ioutil.Discard, r)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.err, err)
		})
	}
}
