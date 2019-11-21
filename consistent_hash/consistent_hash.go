//package consistentHash provide the algorithm of consistent hash.
//this package has no thread safety.
package consistentHash

import (
	"time"
)

type BaseType interface{}

type ISortable interface {
	Add(BaseType)
	Remove(BaseType)
	Sort()
	Len() int
	Index(int) BaseType
}

type Ihash interface {
	Add(BaseType)
	Remove(BaseType)
	Get() BaseType
	Next() (BaseType, bool)
}

type Hash struct {
	list       ISortable
	startIndex int
	curIndex   int
}

func NewHashObj(ls ISortable) *Hash {
	ch := &Hash{
		list: ls,
	}
	return ch.sort().hash()
}

func (ch *Hash) sort() *Hash {
	ch.list.Sort()
	return ch
}

func (ch *Hash) hash() *Hash {
	//		index/len(list)==offset/86400
	//==>   index=(offset*len(lisj))/86400
	now := time.Now()
	offset := now.Hour()*3600 + now.Minute()*60 + now.Second()
	ch.startIndex = (offset * ch.list.Len()) / 86400
	ch.curIndex = ch.startIndex
	return ch
}

func (ch *Hash) Add(ips BaseType) {
	ch.list.Add(ips)
}

func (ch *Hash) Remove(ip BaseType) {
	ch.list.Remove(ip)
}

func (ch *Hash) Get() (ip BaseType) {
	ch.sort().hash()
	return ch.list.Index(ch.startIndex)
}

//If ping the ip that obtain from Get() fail,
//please call Next() method to get a new one in loop.
//If last return true, all of elements in list had be accessed.
func (ch *Hash) Next() (ip BaseType, last bool) {
	if (ch.curIndex+1)%ch.list.Len() == ch.startIndex {
		last = true
	}
	ch.curIndex++
	ch.curIndex = ch.curIndex % ch.list.Len()
	return ch.list.Index(ch.curIndex), last
}
