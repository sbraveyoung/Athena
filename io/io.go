package io

import (
	"bytes"
	"io"
)

type Reader interface {
	io.Reader
	ReadFull(b []byte) error
	ReadN(int) ([]byte, error)
}

type Writer interface {
	io.Writer
	WriteFull([]byte) error
}

type buffer struct {
	*bytes.Buffer
}

func NewReader(buf []byte) (reader Reader) {
	return buffer{Buffer: bytes.NewBuffer(buf)}
}

func (b buffer) ReadFull(buf []byte) error {
	n, err := io.ReadFull(b.Buffer, buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return err
	}
	return nil
}

func (b buffer) ReadN(n int) (buf []byte, err error) {
	buf = make([]byte, n)
	err = b.ReadFull(buf)
	return buf, err
}

//read dump
//write dump
