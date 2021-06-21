package ring_buffer

//empty:head==tail
//full:(head+1)%len(buffer)==tail
type ringBufferWithArray struct {
	cacheBase
	buffer []interface{}
	//head points a empty index that store value next.
	//tail points the oldest index in buffer that has not be read.
	head, tail int
}

func newRingBufferWithArray(cb *cacheBase) *ringBufferWithArray {
	return &ringBufferWithArray{
		cacheBase: *cb,
		buffer:    make([]interface{}, cb.size+1),
		head:      0,
		tail:      0,
	}
}

func (rb *ringBufferWithArray) Get() (val interface{}) {
	if rb.tail == rb.head {
		return nil
	}
	val = rb.buffer[rb.tail]
	rb.tail = (rb.tail + 1) % len(rb.buffer)
	return val
}

func (rb *ringBufferWithArray) Append(val interface{}) {
	rb.buffer[rb.head] = val
	rb.head = (rb.head + 1) % len(rb.buffer)
	if rb.head == rb.tail {
		//ring is full, conver oldest index
		rb.tail = (rb.tail + 1) % len(rb.buffer)
	}
}
