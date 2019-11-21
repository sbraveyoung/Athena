package bitmap

import (
	"errors"
)

var (
	ERROR = errors.New("index is out of range.")
)

const (
	BIT   = 1
	BYTE  = 8 * BIT
	WORD  = 4 * BYTE
	DWORD = 2 * WORD
)

type Bitmap struct {
	bits   uint32
	buffer []byte
}

func New(bits uint32) (*Bitmap, error) {
	if 0 == bits {
		return nil, ERROR
	}
	bm := &Bitmap{
		bits:   bits,
		buffer: make([]byte, (bits-1)/BYTE+1),
	}
	return bm, nil
}

func NewWithString(str string) (*Bitmap, error) {
	// bm, err := New(uint32(len(str) * BYTE))
	// if err != nil {
	// return bm, err
	// }
	// for index, c := range str {
	// for i := 0; i < BYTE; i++ {
	// if (0x1<<i)&c != 0 {
	// err := bm.Set(uint32(index*BYTE + i))
	// if err != nil {
	// return bm, err
	// }
	// }
	// }
	// }
	// return bm, nil
	return &Bitmap{
		bits:   uint32(len(str) * BYTE),
		buffer: []byte(str),
	}, nil
}

func (bm *Bitmap) Range(f func(pos uint32)) {
	for i := uint32(1); i <= bm.bits; i++ {
		f(i)
	}
}

func (bm *Bitmap) Set(pos uint32) error {
	pos -= 1
	index := pos / BYTE
	subIndex := pos % BYTE
	if index >= uint32(len(bm.buffer)) || 0 == pos+1 {
		return ERROR
	}
	bm.buffer[index] |= 0x1 << subIndex
	return nil
}

func (bm *Bitmap) Reset(pos uint32) error {
	pos -= 1
	index := pos / BYTE
	subIndex := pos % BYTE
	if index >= uint32(len(bm.buffer)) || 0 == pos+1 {
		return ERROR
	}
	bm.buffer[index] &= (^(0x1 << subIndex))
	return nil
}

func (bm *Bitmap) Get(pos uint32) bool {
	pos -= 1
	index := pos / BYTE
	subIndex := pos % BYTE
	if index >= uint32(len(bm.buffer)) || 0 == pos+1 {
		return false
	}
	if (0x1<<subIndex)&bm.buffer[index] == 0 {
		return false
	}
	return true
}

func (bm *Bitmap) String() string {
	return string(bm.buffer[:])
}
