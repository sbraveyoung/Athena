package easyio

import (
	"io"
)

type EasyReader interface {
	io.Reader
	ReadFull(b []byte) error
	ReadN(int) ([]byte, error)
	ReadAll() ([]byte, error)
}

type easyReader struct {
	io.Reader
}

func NewEasyReader(r io.Reader) (reader EasyReader) {
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

func (er easyReader) ReadAll() (buf []byte, err error) {
	return io.ReadAll(er)
}

type EasyWriter interface {
	io.Writer
	WriteFull([]byte) error
	// WriteN([]byte, int) error
}

type easyWriter struct {
	io.Writer
}

func NewEasyWriter(w io.Writer) (writer EasyWriter) {
	return easyWriter{Writer: w}
}

func (ew easyWriter) WriteFull(buf []byte) error {
	_, err := ew.Write(buf)
	return err
}

type EasyReadWriter interface {
	EasyReader
	EasyWriter
}

type easyReadWriter struct {
	EasyReader
	EasyWriter
}

func NewEasyReadWriter(rw io.ReadWriter) (readWriter EasyReadWriter) {
	return easyReadWriter{
		EasyReader: NewEasyReader(rw),
		EasyWriter: NewEasyWriter(rw),
	}
}

//read dump
//write dump

func CopyFull(dst EasyWriter, src EasyReader) (err error) {
	_, err = io.Copy(dst, src)
	return err
}
