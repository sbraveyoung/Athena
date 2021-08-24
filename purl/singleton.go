package purl

import (
	"sync"
)

var (
	num   = 1 //0 is nil
	mutex sync.Mutex

	once  sync.Once
	index *indexTable
)

func getNext() int {
	mutex.Lock()
	ret := num
	num++
	mutex.Unlock()
	return ret
}

//https://github.com/golang/go/issues/9477
type indexTable struct {
	//key: string, the original key or value
	//val: int
	forwardIndex sync.Map

	//key: int
	//val: string, the original key or value
	backwardIndex sync.Map

	//key: string, the original key or value
	//val: *sync.WaitGroup
	waitMap sync.Map

	//key: string, the original key or value
	//val: *sync.RWMutex
	forwardMutex sync.Map

	//key: int
	//val: *sync.RWMutex
	backwardMutex sync.Map
}

func init() {
	once.Do(func() {
		index = new(indexTable)
	})
}

func getIndex(skey string) (ikey int) {
	if locker, ok := index.forwardMutex.Load(skey); ok {
		locker.(*sync.RWMutex).RLock()
		key, ok := index.forwardIndex.Load(skey)
		locker.(*sync.RWMutex).RUnlock()
		if ok {
			ikey = key.(int)
			return
		}
	}
	obj := sync.WaitGroup{}
	obj.Add(1)
	wg, wgLoaded := index.waitMap.LoadOrStore(skey, &obj)
	if !wgLoaded {
		//only one goroutine execute here at a same time.
		num := getNext()
		locker := sync.RWMutex{}
		index.forwardMutex.Store(skey, &locker)
		index.backwardMutex.Store(num, &locker)
		locker.Lock()
		index.forwardIndex.Store(skey, num)
		index.backwardIndex.Store(num, skey)
		locker.Unlock()

		wg.(*sync.WaitGroup).Done()
		ikey = num
	} else {
		wg.(*sync.WaitGroup).Wait()
		key, _ := index.forwardIndex.Load(skey)
		ikey = key.(int)
	}
	return ikey
}

func getString(ikey int) (skey string) {
	if locker, ok := index.backwardMutex.Load(ikey); ok {
		locker.(*sync.RWMutex).RLock()
		if key, ok := index.backwardIndex.Load(ikey); ok {
			skey = key.(string)
		}
		locker.(*sync.RWMutex).RUnlock()
	}
	return skey
}
