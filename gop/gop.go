package gop

import (
	"sync"
)

type GOP struct {
	data   []interface{}
	c      *sync.Cond
	alive  bool
	lrSize int //size of last round
	round  int
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

func (gop *GOP) Reset() {
	gop.c.L.Lock()
	gop.lrSize = len(gop.data)
	gop.data = gop.data[0:0:cap(gop.data)]
	gop.round++
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
	round int
}

func NewGOPReader(gop *GOP) *GOPReader {
	return &GOPReader{
		index: 0,
		gop:   gop,
		round: gop.round,
	}
}

func (r *GOPReader) reset() {
	r.index = 0
	r.round = r.gop.round
}

func (r *GOPReader) Read(who string) (p interface{}, alive bool) {
	r.gop.c.L.Lock()
	for {
		alive = r.gop.alive
		if r.round == r.gop.round-1 {
			if r.index < len(r.gop.data) {
				//read too slowly, try to reset reader
				r.reset()
			} else if r.index < r.gop.lrSize {
				//normal, do nothing
			} else {
				r.reset()
			}
		} else if r.round == r.gop.round {
			if r.index < len(r.gop.data) {
				//normal, do nothing
			} else {
				if alive {
					r.gop.c.Wait()
					continue
				} else {
					r.gop.c.L.Unlock()
					return
				}
			}
		} else {
			r.reset()
		}
		break
	}
	p = r.gop.data[r.index]
	r.gop.c.L.Unlock()
	r.index++
	return p, true
}
