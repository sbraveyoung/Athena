package pool

import "sync"

var (
	bufferPool = sync.Pool{
		New: func() interface{} {
			return make(Buffer, 0, 150)
		},
	}
	slicePool = sync.Pool{
		New: func() interface{} {
			return make([]Buffer, 50)
		},
	}
)

func GetBuffer() (buf Buffer, putFunc func(giveUpStr Buffer)) {
	buf = bufferPool.Get().(Buffer)
	//clear slice without gc
	buf = buf[:0]
	return buf, PutBuffer
}

func PutBuffer(giveUpBuf Buffer) {
	bufferPool.Put(giveUpBuf)
}

func GetSlice() (slice []Buffer, putFunc func(giveUpSlice []Buffer)) {
	slice = slicePool.Get().([]Buffer)
	slice = slice[:0]
	return slice, func(giveUpSlice []Buffer) {
		slicePool.Put(giveUpSlice)
	}
}
