package gop

import "sync"

type GOP struct {
	data []interface{}
	c    *sync.Cond
}

func NewGOP() *GOP {
	return &GOP{
		c: sync.NewCond(&sync.Mutex{}),
	}
}

func (gop *GOP) Write(p interface{}) {
	gop.c.L.Lock()
	gop.data = append(gop.data, p)
	gop.c.L.Unlock()

	gop.c.Broadcast()
}

type GOPReader struct {
	index int
	gop   *GOP
}

func NewGOPReader(gop *GOP) *GOPReader {
	return &GOPReader{
		index: 0,
		gop:   gop,
	}
}

func (r *GOPReader) Read() (p interface{}) {
	r.gop.c.L.Lock()
	for r.index >= len(r.gop.data) {
		r.gop.c.Wait()
	}
	r.gop.c.L.Unlock()
	p = r.gop.data[r.index]
	r.index++
	return p
}
