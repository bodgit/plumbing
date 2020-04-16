// Package plumbing is a collection of assorted I/O helpers.
package plumbing

import "io"

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}
