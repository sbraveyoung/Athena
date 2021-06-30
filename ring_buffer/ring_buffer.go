package ring_buffer

const (
	TYPE_ARRAY = "array"
	TYPE_LIST  = "list"
)

type Cache interface {
	Get() interface{}
	Append(interface{})
}

type cacheBase struct {
	size  int
	block bool
}

type cacheBuilder struct {
	cacheBase
	tp string
}

func NewRingBuffer(size int) *cacheBuilder {
	if size <= 0 {
		return nil
	}
	return &cacheBuilder{
		cacheBase: cacheBase{
			size: size,
		},
		tp: TYPE_ARRAY,
	}
}

func (cb *cacheBuilder) EvictType(tp string) *cacheBuilder {
	cb.tp = tp
	return cb
}

func (cb *cacheBuilder) Array() *cacheBuilder {
	return cb.EvictType(TYPE_ARRAY)
}

func (cb *cacheBuilder) List() *cacheBuilder {
	return cb.EvictType(TYPE_LIST)
}

func (cb *cacheBuilder) Block() *cacheBuilder {
	cb.block = true
	return cb
}

func (cb *cacheBuilder) Build() Cache {
	if cb.block {
		return newRingBufferBlocking(&cb.cacheBase)
	}
	switch cb.tp {
	case TYPE_ARRAY:
		return newRingBufferWithArray(&cb.cacheBase)
	case TYPE_LIST:
		return newRingBufferWithList(&cb.cacheBase)
	default:
		panic("ring_buffer: Unknown type " + cb.tp)
	}
}
