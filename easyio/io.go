package easyio

import (
	"io"
)

type Reader interface {
	io.Reader
	ReadFull(b []byte) error
	ReadN(int) ([]byte, error)
}

type easyReader struct {
	io.Reader
}

func NewEasyReader(r io.Reader) (reader Reader) {
	return easyReader{Reader: r}
}

func (er easyReader) ReadFull(buf []byte) error {
	_, err := io.ReadFull(er, buf)
	if err != nil {
		return err
	}
	return nil
}

func (er easyReader) ReadN(n int) (buf []byte, err error) {
	buf = make([]byte, n)
	err = er.ReadFull(buf)
	return buf, err
}

type Writer interface {
	io.Writer
	// WriteFull([]byte) error
	// WriteN([]byte, int) error
}

type easyWriter struct {
	io.Writer
}

func NewEasyWriter(w io.Writer) (writer Writer) {
	return easyWriter{Writer: w}
}

// func (ew easyWriter) WriteFull(buf []byte) (err error) {
// }

// func (ew easyWriter) WriteN(buf []byte, n int) (err error) {
// }

//read dump
//write dump
