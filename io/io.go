package io

import "io"

type Reader interface {
	io.Reader
	ReadFull(b []byte) error
	ReadN(int) ([]byte, error)
}

type Writer interface {
	io.Writer
	WriteFull([]byte) error
}
