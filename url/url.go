package ourl

import (
	"net/url"
	"sort"
	"strings"
	"sync"

	pool "github.com/SmartBrave/utils/easypool"
)

const (
	NORMAL_MODE    = iota //do not use BLACKLIST AND WHITELIST
	WHITELIST_MODE        //only the keys in WHITELIST can be set to index
	BLACKLIST_MODE        //only the keys in BLACKLIST can NOT be set to index
)

var (
	//If we have lots of keys and values,
	//the index could be very large and spent much time on gc.
	mode = NORMAL_MODE

	oneInit sync.Once
	//key: string
	//value: bool
	indexList sync.Map
	length    int
)

func SetMode(m int) {
	oneInit.Do(func() {
		mode = m
	})
}

func AddKeys(keys []string) {
	for _, key := range keys {
		indexList.Store(key, true)
		length++
	}
}

func ParseQuery( /*m Values,*/ rawQuery string) Values {
	var err error
	m := NewValues()
	for rawQuery != "" {
		key := rawQuery
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, rawQuery = key[:i], key[i+1:]
		} else {
			rawQuery = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m.Add(key, value)
	}
	return m
}

//We assume that any key in the map has only one value.
//The type of Values is unsafe for concurrent goroutines.
//Please note: do not use url.URL.Query() to get Values, use ParseQuery() instead.
type Values struct {
	iv map[int]int
	sv map[string]string
}

func NewValues() (m Values) {
	m = Values{
		iv: make(map[int]int, 30),
	}
	switch mode {
	case NORMAL_MODE:
		m.sv = make(map[string]string)
		//do nothing
	case WHITELIST_MODE:
		fallthrough
	case BLACKLIST_MODE:
		m.sv = make(map[string]string, length)
	default:
		//panic
	}

	return m
}

func (v *Values) Get(skey string) string {
	if v == nil {
		return ""
	}
	switch mode {
	case NORMAL_MODE:
		ikey := getIndex(skey)
		return getString(v.iv[ikey])
	case WHITELIST_MODE:
		if _, ok := indexList.Load(skey); ok {
			ikey := getIndex(skey)
			return getString(v.iv[ikey])
		}
		return v.sv[skey]
	case BLACKLIST_MODE:
		if _, ok := indexList.Load(skey); ok {
			return v.sv[skey]
		}
		ikey := getIndex(skey)
		return getString(v.iv[ikey])
	default:
		//panic
		return ""
	}
}

func (v *Values) Set(skey, svalue string) {
	if skey == "" || svalue == "" {
		return
	}
	switch mode {
	case NORMAL_MODE:
		v.iv[getIndex(skey)] = getIndex(svalue)
	case WHITELIST_MODE:
		if _, ok := indexList.Load(skey); ok {
			v.iv[getIndex(skey)] = getIndex(svalue)
			return
		}
		v.sv[skey] = svalue
	case BLACKLIST_MODE:
		if _, ok := indexList.Load(skey); ok {
			v.sv[skey] = svalue
			return
		}
		v.iv[getIndex(skey)] = getIndex(svalue)
	default:
		//panic
		return
	}
}
func (v *Values) Setnx(key, val string) bool {
	if v.Get(key) != "" {
		return false
	}
	v.Set(key, val)
	return true
}

func (v *Values) Add(skey, svalue string) {
	v.Set(skey, svalue)
}

func (v *Values) Del(skey string) {
	switch mode {
	case NORMAL_MODE:
		ikey := getIndex(skey)
		delete(v.iv, ikey)
	case WHITELIST_MODE:
		if _, ok := indexList.Load(skey); ok {
			ikey := getIndex(skey)
			delete(v.iv, ikey)
			return
		}
		delete(v.sv, skey)
	case BLACKLIST_MODE:
		if _, ok := indexList.Load(skey); ok {
			delete(v.sv, skey)
			return
		}
		ikey := getIndex(skey)
		delete(v.iv, ikey)
	default:
		//panic
		return
	}
}

func (v *Values) Encode() string {
	if v == nil {
		return ""
	}
	kvs := make(KVS, 0, len(v.iv)+len(v.sv))
	for ikey, ivalue := range v.iv {
		kvs = append(kvs, kv{
			key:   getString(ikey),
			value: getString(ivalue),
		})
	}
	for skey, svalue := range v.sv {
		kvs = append(kvs, kv{
			key:   skey,
			value: svalue,
		})
	}
	//XXX:why need to sort it?
	sort.Sort(kvs)

	buf, putBuf := pool.GetBuffer()
	defer putBuf(buf)
	for i, _ := range kvs {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(kvs[i].key))
		buf.WriteByte('=')
		//XXX:could do better?
		buf.WriteString(url.QueryEscape(kvs[i].value))
	}
	return buf.String()
}

func (v *Values) Range(f func(key, value string) bool) {
one:
	for ikey, ivalue := range v.iv {
		if !f(getString(ikey), getString(ivalue)) {
			break one
		}
	}
two:
	for skey, svalue := range v.sv {
		if !f(skey, svalue) {
			break two
		}
	}
}

type kv struct {
	key   string
	value string
}

type KVS []kv

func (s KVS) Len() int           { return len(s) }
func (s KVS) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s KVS) Less(i, j int) bool { return s[i].key < s[j].key }
