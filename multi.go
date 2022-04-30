package plumbing

import "io"

type multiWriteCloser struct {
	writeClosers []io.WriteCloser
}

func (t *multiWriteCloser) Write(p []byte) (n int, err error) {
	for _, wc := range t.writeClosers {
		n, err = wc.Write(p)
		if err != nil {
			return
		}

		if n != len(p) {
			err = io.ErrShortWrite

			return
		}
	}

	return len(p), nil
}

func (t *multiWriteCloser) Close() (err error) {
	for _, wc := range t.writeClosers {
		err = wc.Close()
		if err != nil {
			return
		}
	}

	return
}

// MultiWriteCloser creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
//
// Each write is written to each listed writer, one at a time.
// If a listed writer returns an error, that overall write operation
// stops and returns the error; it does not continue down the list.
func MultiWriteCloser(writeClosers ...io.WriteCloser) io.WriteCloser {
	allWriteClosers := make([]io.WriteCloser, 0, len(writeClosers))

	for _, wc := range writeClosers {
		if mwc, ok := wc.(*multiWriteCloser); ok {
			allWriteClosers = append(allWriteClosers, mwc.writeClosers...)
		} else {
			allWriteClosers = append(allWriteClosers, wc)
		}
	}

	return &multiWriteCloser{allWriteClosers}
}
