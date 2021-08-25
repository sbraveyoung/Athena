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
			p, alive := reader1.Read()
			if !alive {
				fmt.Println("reader1 gop disalive!")
				break
			}
			fmt.Println("reader1: ", p.(int))
		}
	}()
	go func() {
		//consumer2
		for {
			p, alive := reader2.Read()
			if !alive {
				fmt.Println("reader2 gop disalive!")
				break
			}
			fmt.Println("reader2: ", p.(int))
		}
	}()
	//producer
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * time.Duration(1))
		gop.Write(i)
	}
	gop.DisAlive()

	time.Sleep(time.Second * time.Duration(2))
	fmt.Println("done")
}
