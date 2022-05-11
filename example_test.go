package plumbing_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/bodgit/plumbing"
)

func ExampleWriteCounter() {
	in := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	writer := plumbing.WriteCounter{}
	reader := io.TeeReader(bytes.NewReader(in), &writer)

	if _, err := io.CopyN(io.Discard, reader, 4); err != nil {
		panic(err)
	}

	if _, err := io.Copy(io.Discard, reader); err != nil {
		panic(err)
	}

	fmt.Println(writer.Count())
	// Output: 10
}

func ExampleTeeReaderAt() {
	// Smallest valid zip archive
	in := []byte{80, 75, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	writer := plumbing.WriteCounter{}
	if _, err := zip.NewReader(plumbing.TeeReaderAt(bytes.NewReader(in), &writer), int64(len(in))); err != nil {
		panic(err)
	}

	fmt.Println(writer.Count())
	// Output: 44
}

func ExampleTeeReadCloser() {
	in := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	writer := plumbing.WriteCounter{}
	reader := plumbing.TeeReadCloser(io.NopCloser(bytes.NewReader(in)), &writer)

	defer reader.Close()

	if _, err := io.Copy(io.Discard, reader); err != nil {
		panic(err)
	}

	fmt.Println(writer.Count())
	// Output: 10
}

func ExamplePaddedReader() {
	in := []byte{1, 2, 3, 4}

	reader := plumbing.PaddedReader(bytes.NewReader(in), 8, 0)
	writer := new(bytes.Buffer)

	if _, err := io.Copy(writer, reader); err != nil {
		panic(err)
	}

	fmt.Println(writer.Bytes())
	// Output: [1 2 3 4 0 0 0 0]
}

func ExampleNopWriteCloser() {
	writer := plumbing.NopWriteCloser(new(bytes.Buffer))

	fmt.Println(writer.Close())
	// Output: <nil>
}

func ExampleMultiWriteCloser() {
	in := []byte{0, 1, 2, 3}
	b1, b2 := new(bytes.Buffer), new(bytes.Buffer)
	writer := plumbing.MultiWriteCloser(plumbing.NopWriteCloser(b1), plumbing.NopWriteCloser(b2))

	if _, err := writer.Write(in); err != nil {
		panic(err)
	}

	if err := writer.Close(); err != nil {
		panic(err)
	}

	fmt.Println(b1.Bytes(), b2.Bytes())
	// Output: [0 1 2 3] [0 1 2 3]
}

func ExampleMultiReadCloser() {
	b1, b2 := bytes.NewReader([]byte{0, 1, 2, 3}), bytes.NewReader([]byte{4, 5, 6, 7})
	r := plumbing.MultiReadCloser(io.NopCloser(b1), io.NopCloser(b2))
	w := new(bytes.Buffer)

	if _, err := io.Copy(w, r); err != nil {
		panic(err)
	}

	if err := r.Close(); err != nil {
		panic(err)
	}

	fmt.Println(w.Bytes())
	// Output: [0 1 2 3 4 5 6 7]
}

func ExampleLimitReadCloser() {
	in := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	reader := plumbing.LimitReadCloser(io.NopCloser(bytes.NewReader(in)), 5)
	writer := new(bytes.Buffer)

	if _, err := io.Copy(writer, reader); err != nil {
		panic(err)
	}

	if err := reader.Close(); err != nil {
		panic(err)
	}

	fmt.Println(writer.Bytes())
	// Output: [0 1 2 3 4]
}
