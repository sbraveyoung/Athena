package ring_buffer

type ringBufferBlocking struct {
	cacheBase
	c chan interface{}
}

func newRingBufferBlocking(cb *cacheBase) *ringBufferBlocking {
	return &ringBufferBlocking{
		cacheBase: *cb,
		c:         make(chan interface{}, cb.size),
	}
}

func (rb *ringBufferBlocking) Get() (val interface{}) {
	return <-rb.c
}

func (rb *ringBufferBlocking) Append(val interface{}) {
	rb.c <- val
}
