package io

type Reader interface {
	Read(b []byte) error
	ReadN(int) ([]byte, error)
}

type Writer interface {
	Write([]byte) error
}
