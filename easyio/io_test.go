package easyio

import (
	"io"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/SmartBrave/Athena/easyerrors"
)

const (
	TEST_BUF_LEN = 1024 * 256
)

func TestEasyReadWriter1(t *testing.T) {
	src := make([]byte, TEST_BUF_LEN)
	dst := make([]byte, 0)
	rand.Seed(time.Now().Unix())
	_, err := rand.Read(src)
	if err != nil {
		t.Errorf("rand.Read error:%v", err)
	}
	// fmt.Printf("src  : %x\n", src)

	rw := NewEasyReadWriter()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		total := 0
		for {
			r := rand.Int()%20 + 1
			if total+r > len(src) {
				r = len(src) - total
			}
			_, err := rw.Write(src[total : total+r])
			if err != nil {
				// fmt.Println("rw.write error:", err)
				break
			}
			total += r

			if total >= len(src) {
				rw.Close()
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			r := rand.Int()%20 + 1
			if len(dst)+r > len(src) {
				r = len(src) - len(dst)
			}
			if r == 0 {
				break
			}
			b := make([]byte, r)
			// fmt.Printf("before read, rw:%+v\n", rw)
			n, err := rw.Read(b)
			// fmt.Printf("read from rw, n:%d, len(b):%d, b:%x\n", n, len(b), b)
			// fmt.Printf("after read, rw:%+v\n", rw)
			if err == io.EOF {
				dst = append(dst, b...)
				break
			}
			if err != nil {
				// fmt.Println("rw.write error:", err)
				break
			}
			dst = append(dst, b[:n]...)
		}
	}()
	wg.Wait()

	// fmt.Printf("dst  : %x\n", dst)
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("read data is:%+x, expect:%+x", dst, src)
	}
}

func TestEasyReadWriter2(t *testing.T) {
	src := make([]byte, TEST_BUF_LEN)
	rand.Seed(time.Now().Unix())
	_, err := rand.Read(src)
	if err != nil {
		t.Errorf("rand.Read error:%v", err)
	}
	// fmt.Printf("src  : %x\n", src)

	rw := NewEasyReadWriter()
	err1 := rw.WriteFull(src)
	rw.Close()
	dst, err2 := rw.ReadAll()
	if err = easyerrors.HandleMultiError(easyerrors.Simple(), err1, err2); err != nil {
		t.Errorf("write of read error, err1:%v, err2:%v", err1, err2)
	}
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("read data is:%+x, expect:%+x", dst, src)
	}
}

func TestEasyReadWriter3(t *testing.T) {
	src := make([]byte, TEST_BUF_LEN)
	rand.Seed(time.Now().Unix())
	_, err := rand.Read(src)
	if err != nil {
		t.Errorf("rand.Read error:%v", err)
	}
	dst := make([]byte, TEST_BUF_LEN)

	rw := NewEasyReadWriter()
	err1 := rw.WriteFull(src)
	rw.Close()
	err2 := rw.ReadFull(dst)
	if err = easyerrors.HandleMultiError(easyerrors.Simple(), err1, err2); err != nil {
		t.Errorf("write of read error, err1:%v, err2:%v", err1, err2)
	}
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("read data is:%+x, expect:%+x", dst, src)
	}
}
