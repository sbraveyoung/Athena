package broadcast

import (
	"sync"
)

type Broadcast struct {
	data   []interface{} //TODO: rewrite with type parameter
	meta   []interface{}
	c      *sync.Cond
	alive  bool
	lrSize int //size of last round
	round  int
	wIndex int
}

func NewBroadcast(metaNum int) *Broadcast {
	return &Broadcast{
		alive: true,
		//TODO: implement a RWMutex to replace Mutex here!
		c:    sync.NewCond(&sync.Mutex{}),
		meta: make([]interface{}, 0, metaNum),
	}
}

func (bd *Broadcast) WriteMeta(meta interface{}) {
	bd.c.L.Lock()
	defer bd.c.L.Unlock()

	if len(bd.meta) < cap(bd.meta) {
		bd.meta = append(bd.meta, meta)
	}
}

func (bd *Broadcast) Write(p interface{}) {
	bd.c.L.Lock()
	defer bd.c.L.Unlock()

	if len(bd.meta) < cap(bd.meta) {
		return
	}

	if bd.wIndex >= len(bd.data) {
		bd.data = append(bd.data, p)
	} else {
		bd.data[bd.wIndex] = p
	}
	bd.wIndex++

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
	mIndex int
	bd     *Broadcast
	round  int
}

func NewBroadcastReader(bd *Broadcast) *BroadcastReader {
	return &BroadcastReader{
		rIndex: 0,
		mIndex: 0,
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
		if len(r.bd.meta) < cap(r.bd.meta) {
			r.bd.c.Wait()
			continue
		}

		alive = r.bd.alive

		if r.mIndex < len(r.bd.meta) {
			p = r.bd.meta[r.mIndex]
			r.mIndex++
			return
		}

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
