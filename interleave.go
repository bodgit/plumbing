package plumbing

import (
	"errors"
	"io"
)

var (
	errInvalidStripe    = errors.New("invalid stripe size")
	errIncompleteStripe = errors.New("buffer size should be a multiple of the product of the number of readers and stripe size")
	errUnexpectedEOF    = errors.New("got EOF when the first reader didn't return that")
	errMissingEOF       = errors.New("first reader returned EOF, subsequent reader didn't")
)

type interleaveReader struct {
	stripe  int
	readers []io.Reader
}

func (ir *interleaveReader) Read(p []byte) (int, error) {
	if ir.stripe < 1 {
		return 0, errInvalidStripe
	}

	if len(p)%len(ir.readers)*ir.stripe != 0 {
		return 0, errIncompleteStripe
	}

	if len(ir.readers) > 0 {
		total, seenEOF := 0, false
		for i := 0; i < len(p); i += ir.stripe {
			n, err := ir.readers[i/ir.stripe%len(ir.readers)].Read(p[i : i+ir.stripe])
			total += n
			if err == io.EOF {
				if !seenEOF && i > 0 {
					return total, errUnexpectedEOF
				}
				seenEOF = true
			} else if seenEOF {
				return total, errMissingEOF
			}
			if err != nil && err != io.EOF {
				return total, err
			}
		}
		if seenEOF {
			ir.readers = nil
		}
		return total, nil
	}
	return 0, io.EOF
}

// InterleaveReader returns an io.Reader that reads alternate bytes from the
// provided input readers. It is required that each input is the same length.
func InterleaveReader(readers ...io.Reader) io.Reader {
	return InterleaveStripeReader(1, readers...)
}

// InterleaveStripeReader returns an io.Reader that reads alternate stripe
// number of bytes from the provided input readers. It is required that each
// input is the same length.
func InterleaveStripeReader(stripe int, readers ...io.Reader) io.Reader {
	r := make([]io.Reader, len(readers))
	copy(r, readers)
	return &interleaveReader{stripe, r}
}
