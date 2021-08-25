package gop

import "sync"

type GOP struct {
	alive bool
	data  []interface{}
	c     *sync.Cond
}

func NewGOP() *GOP {
	return &GOP{
		alive: true,
		c:     sync.NewCond(&sync.Mutex{}),
	}
}

func (gop *GOP) Write(p interface{}) {
	gop.c.L.Lock()
	gop.data = append(gop.data, p)
	gop.c.L.Unlock()

	gop.c.Broadcast()
}

func (gop *GOP) DisAlive() {
	gop.c.L.Lock()
	gop.alive = false
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

func (r *GOPReader) Read() (p interface{}, alive bool) {
	r.gop.c.L.Lock()
	for r.gop.alive && r.index >= len(r.gop.data) {
		r.gop.c.Wait()
	}

	alive = r.gop.alive
	if !alive && r.index < len(r.gop.data) {
		alive = true
	}
	r.gop.c.L.Unlock()

	if alive {
		p = r.gop.data[r.index]
		r.index++
	}
	return p, alive
}
