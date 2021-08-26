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
	wIndex int
}

func NewGOP() *GOP {
	return &GOP{
		alive: true,
		c:     sync.NewCond(&sync.Mutex{}),
	}
}

func (gop *GOP) Write(p interface{}) {
	gop.c.L.Lock()
	if gop.round == 0 || gop.wIndex >= len(gop.data) {
		gop.data = append(gop.data, p)
	} else {
		gop.data[gop.wIndex] = p
	}
	gop.wIndex++
	gop.c.L.Unlock()

	gop.c.Broadcast()
}

func (gop *GOP) Reset() {
	gop.c.L.Lock()
	gop.lrSize = len(gop.data)
	gop.round++
	gop.wIndex = 0
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
	rIndex int
	gop    *GOP
	round  int
}

func NewGOPReader(gop *GOP) *GOPReader {
	return &GOPReader{
		rIndex: 0,
		gop:    gop,
		round:  gop.round,
	}
}

func (r *GOPReader) reset() {
	r.rIndex = 0
	r.round = r.gop.round
}

func (r *GOPReader) Read() (p interface{}, alive bool) {
	r.gop.c.L.Lock()
	defer r.gop.c.L.Unlock()

	for {
		alive = r.gop.alive
		if r.round == r.gop.round-1 {
			if r.rIndex < r.gop.wIndex {
				//read too slowly, try to reset reader
				r.reset()
			} else if r.rIndex < r.gop.lrSize {
				//normal, do nothing
			} else {
				r.reset()
			}
		} else if r.round == r.gop.round {
			if r.rIndex < r.gop.wIndex {
				//normal, do nothing
			} else {
				if alive {
					r.gop.c.Wait()
					continue
				} else {
					return
				}
			}
		} else {
			r.reset()
		}
		break
	}
	p = r.gop.data[r.rIndex]
	r.rIndex++
	return p, true
}
