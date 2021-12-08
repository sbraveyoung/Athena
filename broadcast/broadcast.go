package broadcast

import (
	"sync"
)

type Broadcast struct {
	data   []interface{} //TODO: rewrite with type parameter
	c      *sync.Cond
	start  bool
	alive  bool
	lrSize int //size of last round
	round  int
	wIndex int
}

func NewBroadcast() *Broadcast {
	return &Broadcast{
		start: false,
		alive: true,
		//TODO: implement a RWMutex to replace Mutex here!
		c: sync.NewCond(&sync.Mutex{}),
	}
}

func (bd *Broadcast) Start() {
	bd.start = true
}

func (bd *Broadcast) Write(p interface{}) {
	bd.c.L.Lock()
	if bd.wIndex >= len(bd.data) {
		bd.data = append(bd.data, p)
	} else {
		bd.data[bd.wIndex] = p
	}
	bd.wIndex++
	bd.c.L.Unlock()

	bd.c.Broadcast()
}

func (bd *Broadcast) Reset() {
	bd.c.L.Lock()
	bd.lrSize = bd.wIndex
	bd.round++
	bd.wIndex = 0
	bd.c.L.Unlock()
	bd.c.Broadcast()
}

func (bd *Broadcast) DisAlive() {
	bd.c.L.Lock()
	bd.alive = false
	bd.c.L.Unlock()

	bd.c.Broadcast()
}

type BroadcastReader struct {
	rIndex int
	bd     *Broadcast
	round  int
}

func NewBroadcastReader(bd *Broadcast) *BroadcastReader {
	return &BroadcastReader{
		rIndex: 0,
		bd:     bd,
		round:  bd.round,
	}
}

func (r *BroadcastReader) reset() {
	r.rIndex = 0
	r.round = r.bd.round
}

func (r *BroadcastReader) Read() (p interface{}, alive bool) {
	r.bd.c.L.Lock()
	defer r.bd.c.L.Unlock()

	for {
		if !r.bd.start {
			r.bd.c.Wait()
			continue
		}

		alive = r.bd.alive
		if r.round == r.bd.round-1 {
			if r.rIndex < r.bd.wIndex {
				//read too slowly, try to reset reader
				r.reset()
			} else if r.rIndex < r.bd.lrSize {
				//normal, do nothing
			} else {
				r.reset()
			}
		} else if r.round == r.bd.round {
			if r.rIndex < r.bd.wIndex {
				//normal, do nothing
			} else {
				if alive {
					r.bd.c.Wait()
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
	p = r.bd.data[r.rIndex] //BUG: maybe panic with index 0
	r.rIndex++
	return p, true
}
