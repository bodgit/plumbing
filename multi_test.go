package plumbing_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/bodgit/plumbing"
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

//nolint:funlen
func TestMultiWriteCloser(t *testing.T) {
	t.Parallel()

	in := []byte("abcdefghij")

	tables := []struct {
		name     string
		writer   io.WriteCloser
		n        int
		writeErr error
		closeErr error
	}{
		{
			"success",
			plumbing.NopWriteCloser(new(bytes.Buffer)),
			10,
			nil,
			nil,
		},
		{
			"nested",
			plumbing.MultiWriteCloser(plumbing.NopWriteCloser(new(bytes.Buffer))),
			10,
			nil,
			nil,
		},
		{
			"write",
			plumbing.NopWriteCloser(errorWriter{}),
			0,
			errWrite,
			nil,
		},
		{
			"close",
			errorWriteCloser{},
			10,
			nil,
			errClose,
		},
		{
			"partial",
			plumbing.NopWriteCloser(partialWriter{}),
			9,
			io.ErrShortWrite,
			nil,
		},
	}

	for _, table := range tables {
		table := table
		t.Run(table.name, func(t *testing.T) {
			t.Parallel()
			dst := plumbing.NopWriteCloser(new(bytes.Buffer))
			w := plumbing.MultiWriteCloser(table.writer, dst)
			n, err := w.Write(in)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.writeErr, err)
			err = w.Close()
			assert.Equal(t, table.closeErr, err)
		})
	}
}

type errorReadCloser struct {
	io.Reader
}

func (errorReadCloser) Close() error {
	return errClose
}

type earlyEOFBytesReader struct {
	r *bytes.Reader
}

func (r *earlyEOFBytesReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	// Return an early EOF when bytes have been read and there's none left
	if n > 0 && err == nil && r.r.Len() == 0 {
		err = io.EOF
	}

	return
}

func TestMultiReadCloser(t *testing.T) {
	t.Parallel()

	tables := []struct {
		name     string
		readers  []io.ReadCloser
		n        int64
		expected []byte
		readErr  error
		closeErr error
	}{
		{
			"success",
			[]io.ReadCloser{
				io.NopCloser(&earlyEOFBytesReader{bytes.NewReader([]byte("abcde"))}),
				io.NopCloser(bytes.NewReader([]byte("fghij"))),
			},
			10,
			[]byte("abcdefghij"),
			nil,
			nil,
		},
		{
			"nested",
			[]io.ReadCloser{
				plumbing.MultiReadCloser(io.NopCloser(bytes.NewReader([]byte("abcdefghij")))),
			},
			10,
			[]byte("abcdefghij"),
			nil,
			nil,
		},
		{
			"close",
			[]io.ReadCloser{
				io.NopCloser(bytes.NewReader([]byte("abcde"))),
				&errorReadCloser{bytes.NewReader([]byte("fghij"))},
			},
			10,
			[]byte("abcdefghij"),
			nil,
			errClose,
		},
	}

	for _, table := range tables {
		table := table
		t.Run(table.name, func(t *testing.T) {
			t.Parallel()
			r := plumbing.MultiReadCloser(table.readers...)
			b := new(bytes.Buffer)
			n, err := io.Copy(b, r)
			assert.Equal(t, table.n, n)
			assert.Equal(t, table.readErr, err)
			assert.Equal(t, table.expected, b.Bytes())
			err = r.Close()
			assert.Equal(t, table.closeErr, err)
		})
	}
}
