package plumbing_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/bodgit/plumbing"
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
	t.Parallel()

	in := []byte("abcdefghij")

	tables := []struct {
		name   string
		reader io.ReaderAt
		writer io.Writer
		length int
		offset int64
		n      int
		err    error
	}{
		{
			"success",
			bytes.NewReader(in),
			io.Discard,
			3,
			2,
			3,
			nil,
		},
		{
			"fail",
			bytes.NewReader(in),
			errorWriter{},
			3,
			2,
			0,
			errWrite,
		},
	}

	for _, table := range tables {
		table := table
		t.Run(table.name, func(t *testing.T) {
			t.Parallel()
			r := plumbing.TeeReaderAt(table.reader, table.writer)
			dst := make([]byte, table.length)
			n, err := r.ReadAt(dst, table.offset)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.err, err)
		})
	}
}

func TestTeeReadCloser(t *testing.T) {
	t.Parallel()

	in := []byte("abcdefghij")

	tables := []struct {
		name   string
		reader io.ReadCloser
		writer io.Writer
		n      int64
		err    error
	}{
		{
			"success",
			io.NopCloser(bytes.NewReader(in)),
			io.Discard,
			10,
			nil,
		},
		{
			"fail",
			io.NopCloser(bytes.NewReader(in)),
			errorWriter{},
			0,
			errWrite,
		},
	}

	for _, table := range tables {
		table := table
		t.Run(table.name, func(t *testing.T) {
			t.Parallel()
			r := plumbing.TeeReadCloser(table.reader, table.writer)
			defer r.Close()
			n, err := io.Copy(io.Discard, r)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.err, err)
		})
	}
}
