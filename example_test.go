package plumbing

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

func ExampleWriteCounter() {
	in := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	writer := WriteCounter{}
	reader := io.TeeReader(bytes.NewReader(in), &writer)
	if _, err := io.CopyN(ioutil.Discard, reader, 4); err != nil {
		panic(err)
	}
	if _, err := io.Copy(ioutil.Discard, reader); err != nil {
		panic(err)
	}

	fmt.Println(writer.Count())
	// Output: 10
}

func ExampleTeeReaderAt() {
	// Smallest valid zip archive
	in := []byte{80, 75, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	writer := WriteCounter{}
	if _, err := zip.NewReader(TeeReaderAt(bytes.NewReader(in), &writer), int64(len(in))); err != nil {
		panic(err)
	}

	fmt.Println(writer.Count())
	// Output: 44
}

func ExampleTeeReadCloser() {
	in := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	writer := WriteCounter{}
	reader := TeeReadCloser(ioutil.NopCloser(bytes.NewReader(in)), &writer)
	defer reader.Close()
	if _, err := io.Copy(ioutil.Discard, reader); err != nil {
		panic(err)
	}

	fmt.Println(writer.Count())
	// Output: 10
}

func ExamplePaddedReader() {
	in := []byte{1, 2, 3, 4}
	reader := PaddedReader(bytes.NewReader(in), 8, 0)
	writer := new(bytes.Buffer)
	if _, err := io.Copy(writer, reader); err != nil {
		panic(err)
	}

	fmt.Println(writer.Bytes())
	// Output: [1 2 3 4 0 0 0 0]
}

func ExampleNopWriteCloser() {
	writer := NopWriteCloser(new(bytes.Buffer))

	fmt.Println(writer.Close())
	// Output: <nil>
}

func ExampleMultiWriteCloser() {
	in := []byte{0, 1, 2, 3}
	b1, b2 := new(bytes.Buffer), new(bytes.Buffer)
	writer := MultiWriteCloser(NopWriteCloser(b1), NopWriteCloser(b2))
	if _, err := writer.Write(in); err != nil {
		panic(err)
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}

	fmt.Println(b1.Bytes(), b2.Bytes())
	// Output: [0 1 2 3] [0 1 2 3]
}
