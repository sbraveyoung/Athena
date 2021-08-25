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
			fmt.Println("reader1: ", p)
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
			fmt.Println("reader2: ", p)
		}
	}()

	//producer
	for i := 0; i < 20; i++ {
		//time.Sleep(time.Second * time.Duration(1))
		gop.Write(i)
	}

	time.Sleep(time.Second * time.Duration(1))
	gop.Reset()
	for i := 'a'; i < 'h'; i++ {
		gop.Write(i)
	}

	gop.DisAlive()

	time.Sleep(time.Second * time.Duration(2))
	fmt.Println("done")
}
