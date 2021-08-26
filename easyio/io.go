package easyio

import (
	"io"

	"github.com/pkg/errors"
)

type EasyReader interface {
	io.Reader
	ReadFull(b []byte) error
	ReadN(uint32) ([]byte, error)
	ReadAll() ([]byte, error)
}

type easyReader struct {
	io.Reader
}

func NewEasyReader(r io.Reader) (reader EasyReader) {
	return easyReader{Reader: r}
}

func (er easyReader) ReadFull(b []byte) (err error) {
	var n int
	n, err = io.ReadFull(er, b)
	if err == io.EOF {
		return err
	}
	if err != nil {
		return errors.Wrap(err, "rtmp.conn.Read")
	}
	if n != len(b) {
		return errors.New("do not read enough data from conn")
	}
	return nil
}

func (er easyReader) ReadN(n uint32) (b []byte, err error) {
	b = make([]byte, n)
	err = er.ReadFull(b)
	return b, err
}

func (er easyReader) ReadAll() (b []byte, err error) {
	return io.ReadAll(er)
}

type EasyWriter interface {
	io.Writer
	WriteFull([]byte) error
	// Flush() error
	// WriteN([]byte, int) error
}

type easyWriter struct {
	io.Writer
	// *bufio.Writer
}

func NewEasyWriter(w io.Writer) (writer EasyWriter) {
	//return easyWriter{Writer: bufio.NewWriter(w)}
	return easyWriter{Writer: w}
}

func (ew easyWriter) WriteFull(b []byte) error {
	n, err := ew.Write(b)
	if err != nil {
		return errors.Wrap(err, "conn.Write")
	}
	if n != len(b) {
		return errors.New("do not write enough data to conn")
	}
	return nil
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
