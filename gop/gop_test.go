package gop

import (
	"fmt"
	"testing"
	"time"
)

func TestGOP(t *testing.T) {
	gop := NewGOP()
	reader1 := NewGOPReader(gop)
	reader2 := NewGOPReader(gop)

	go func() {
		//consumer1
		for {
			p := reader1.Read().(int)
			fmt.Println("reader1: ", p)
		}
	}()
	go func() {
		//consumer2
		for {
			p := reader2.Read().(int)
			fmt.Println("reader2: ", p)
		}
	}()
	//producer
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second * time.Duration(2))
		gop.Write(i)
	}
}
