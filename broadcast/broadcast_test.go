package broadcast

import (
	"fmt"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	bd := NewBroadcast()
	reader1 := NewBroadcastReader(bd)
	reader2 := NewBroadcastReader(bd)

	go func() {
		//consumer1
		for {
			p, alive := reader1.Read()
			if !alive {
				fmt.Println("reader1 bd disalive!")
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
				fmt.Println("reader2 bd disalive!")
				break
			}
			fmt.Println("reader2: ", p)
		}
	}()

	//producer
	for i := 0; i < 20; i++ {
		//time.Sleep(time.Second * time.Duration(1))
		bd.Write(i)
	}

	time.Sleep(time.Second * time.Duration(1))
	bd.Reset()
	for i := 'a'; i < 'h'; i++ {
		bd.Write(i)
	}

	bd.DisAlive()

	time.Sleep(time.Second * time.Duration(2))
	fmt.Println("done")
}
