package ring_buffer

import "container/ring"

//empty:head==tail
//full:head.Next==tail
type ringBufferWithList struct {
	cacheBase
	//head points a empty node that store value next.
	//tail points the oldest node in ring that has not be read.
	head, tail *ring.Ring
}

func newRingBufferWithList(cb *cacheBase) *ringBufferWithList {
	r := ring.New(cb.size + 1)
	return &ringBufferWithList{
		cacheBase: *cb,
		head:      r,
		tail:      r,
	}
}

func (rb *ringBufferWithList) Get() (val interface{}) {
	if rb.tail == rb.head {
		return nil
	}
	val = rb.tail.Value
	rb.tail = rb.tail.Next()
	return val
}

func (rb *ringBufferWithList) Append(val interface{}) {
	rb.head.Value = val
	rb.head = rb.head.Next()
	if rb.head == rb.tail {
		//ring is full, conver oldest node
		rb.tail = rb.tail.Next()
	}
}
