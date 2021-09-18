package easyio

import (
	"container/ring"
	"io"
	"sync"

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
	return &easyReader{Reader: r}
}

func (er *easyReader) ReadFull(b []byte) (err error) {
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

func (er *easyReader) ReadN(n uint32) (b []byte, err error) {
	b = make([]byte, n)
	err = er.ReadFull(b)
	return b, err
}

func (er *easyReader) ReadAll() (b []byte, err error) {
	return io.ReadAll(er)
}

type EasyWriter interface {
	io.Writer
	WriteFull([]byte) error
}

type easyWriter struct {
	io.Writer
}

func NewEasyWriter(w io.Writer) (writer EasyWriter) {
	return &easyWriter{Writer: w}
}

func (ew *easyWriter) WriteFull(b []byte) error {
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
	Close()
}

type rwBase struct {
	buf      []byte
	writeOff int
	lrOff    int //writeOff of last round
	readOff  int
}

func newRWBase() *rwBase {
	v := &rwBase{
		// buf:      make([]byte, 0, 512), //TODO: 512 should be redefined by benchmark
		buf:      make([]byte, 0), //TODO: 512 should be redefined by benchmark
		writeOff: 0,
		readOff:  0,
		lrOff:    0,
	}
	return v
}

func (v *rwBase) reset() {
	v.writeOff = 0
	v.readOff = 0
	v.lrOff = 0
}

//If v.buf has enough space to write b, write it and return true,
//else return false and write nothing.
func (v *rwBase) write(b []byte) bool {
	if v.readOff <= v.writeOff {
		writed := 0
		copyAble := len(v.buf) - v.writeOff
		appendAble := cap(v.buf) - len(v.buf)
		//BUG: 如果给一个空节点插入超过 max 的值，会永远返回失败，永远无法插入
		if copyAble+appendAble+v.readOff < len(b) && len(v.buf) > 64 { //TODO: 4096 should be redefined by benchmark
			return false
		}

		n := copy(v.buf[v.writeOff:], b)
		writed += n
		v.writeOff += n
		// fmt.Printf("%x", b[:n])
		if writed == len(b) {
			return true
		}

		v.buf = append(v.buf, b[writed:]...)
		writed += (len(v.buf) - v.writeOff)
		v.writeOff = len(v.buf)
		// fmt.Printf("%x", b[:n])
		if writed == len(b) {
			return true
		}

		if v.readOff > len(b)-writed {
			n = copy(v.buf[:v.readOff], b[writed:])
			writed += n
			v.writeOff = n
			//assert(writed==len(b))
			return true
		} else {
			v.buf = append(v.buf, b[writed:]...)
			writed += n
			v.writeOff += len(v.buf)
			//assert(writed==len(b))
			return true
		}
	} else {
		writeAble := v.readOff - v.writeOff
		if writeAble < len(b) {
			return false
		} else {
			v.writeOff += copy(v.buf[v.writeOff:], b)
			return true
		}
	}
}

func (v *rwBase) read(b []byte) (n int, err error) {
	if v.readOff <= v.writeOff {
		readAble := v.writeOff - v.readOff
		if readAble <= len(b) {
			v.readOff += copy(b, v.buf[v.readOff:v.writeOff])
			return readAble, nil
		} else {
			v.readOff += copy(b, v.buf[v.readOff:])
			return len(b), nil
		}
	} else {
		suffixReadAble := v.lrOff - v.readOff
		if suffixReadAble <= len(b) {
			v.readOff += copy(b, v.buf[v.readOff:v.lrOff])

			prefixReadAble := v.writeOff
			remain := len(b) - suffixReadAble
			if prefixReadAble <= remain {
				if n := copy(b[suffixReadAble:], v.buf[:prefixReadAble]); n != 0 {
					v.readOff = n
				}
			} else {
				if n := copy(b[suffixReadAble:], v.buf[:remain]); n != 0 {
					v.readOff = n
				}
			}
			return suffixReadAble + prefixReadAble, nil
		} else {
			v.readOff += copy(b, v.buf[v.readOff:])
			return len(b), nil
		}
	}
}

type easyReadWriter struct {
	writeBufRing *ring.Ring //the elements of ring is rwBase
	readBufRing  *ring.Ring
	c            *sync.Cond
	alive        bool
}

func NewEasyReadWriter() (rw EasyReadWriter) {
	r := ring.New(1)
	return &easyReadWriter{
		writeBufRing: r,
		readBufRing:  r,
		c:            sync.NewCond(&sync.Mutex{}),
		alive:        true,
	}
}

//func (rw *easyReadWriter) debug(from string) {
//	fmt.Printf("-------------------------------------------------------\nrw.writeBufRing: %p, rw.readBufRing: %p, rw.alive:%t :", rw.writeBufRing, rw.readBufRing, rw.alive)
//	var start *ring.Ring = nil
//	var end *ring.Ring = rw.readBufRing
//	if from == "write" {
//		end = rw.writeBufRing
//	}
//	for ; start != end; start = start.Next() {
//		if start == nil {
//			start = rw.readBufRing
//			if from == "write" {
//				start = rw.writeBufRing
//			}
//		}
//
//		if start.Value == nil {
//			break
//		}
//		v := start.Value.(*rwBase)
//		fmt.Printf("{addr:%p, node buf:%x, readOff:%d, writeOff:%d, lrOff:%d, len(buf):%d, cap(buf):%d, node prev:%p, node next:%p} ---> ", start, v.buf, v.readOff, v.writeOff, v.lrOff, len(v.buf), cap(v.buf), start.Prev(), start.Next())
//	}
//	fmt.Printf("nil\n--------------------------------------------------------\n")
//}

func (rw *easyReadWriter) Close() {
	rw.c.L.Lock()
	rw.alive = false
	rw.c.L.Unlock()

	rw.c.Broadcast()
}

func (rw *easyReadWriter) Write(b []byte) (n int, err error) {
	defer rw.c.Broadcast()
	rw.c.L.Lock()
	defer rw.c.L.Unlock()

	defer func() {
		// fmt.Printf("write b, len:%d, data:%x\n", len(b), b)
		// rw.debug("write")
	}()

insert:
	if rw.writeBufRing.Value == nil {
		rw.writeBufRing.Value = newRWBase()
	}

	v := rw.writeBufRing.Value.(*rwBase)
	if v.write(b) {
		return len(b), nil
	}

	if rw.writeBufRing.Next() == rw.readBufRing {
		rw.writeBufRing.Link(ring.New(1))
	}
	rw.writeBufRing = rw.writeBufRing.Next()
	goto insert
}

func (rw *easyReadWriter) WriteFull(b []byte) (err error) {
	n, err := rw.Write(b)
	if err != nil {
		return errors.Wrap(err, "conn.Write")
	}
	if n != len(b) {
		return errors.New("do not write enough data to conn")
	}
	return nil
}

func (rw *easyReadWriter) Read(b []byte) (n int, err error) {
	rw.c.L.Lock()
	defer rw.c.L.Unlock()

	defer func() {
		if n == 0 && !rw.alive {
			err = io.EOF
		}
		// fmt.Printf("read b, len:%d, n:%d, data:%x\n", len(b), n, b)
		// rw.debug("read")
	}()

start:
	node := rw.readBufRing
	totalRead := 0
	for {
		if node.Value == nil {
			rw.c.Wait()
			goto start
		}
		v := node.Value.(*rwBase)
		n, _ := v.read(b[totalRead:])
		totalRead += n
		if totalRead == len(b) {
			rw.readBufRing = node
			return totalRead, nil
		}

		if node == rw.writeBufRing {
			return totalRead, nil
		}
		v.reset()
		node = node.Next()
	}
	return 0, nil
}

func (rw *easyReadWriter) ReadFull(b []byte) (err error) {
	var n int
	n, err = io.ReadFull(rw, b)
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
func (rw *easyReadWriter) ReadN(n uint32) (b []byte, err error) {
	b = make([]byte, n)
	err = rw.ReadFull(b)
	return b, err
}
func (rw *easyReadWriter) ReadAll() (b []byte, err error) {
	return io.ReadAll(rw)
}

//read dump
//write dump

func CopyFull(dst EasyWriter, src EasyReader) (err error) {
	_, err = io.Copy(dst, src)
	return err
}
