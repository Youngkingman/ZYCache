package loopqueue

import (
	"fmt"
	"sync"
	"testing"
)

var Q LoopQueue
var wg sync.WaitGroup

const MAX_TEST_COUNT = 1000

func Create() {
	index := DataItem{}
	ret := Q.Push(index)
	if ret {
		fmt.Println("PushOk", "index=", index)
	} else {
		fmt.Println("PushError", "index=", index)
	}

}
func Consum() {
	ret, data := Q.Pop()
	if ret {
		fmt.Println("PopSucc", "data=", data)
	} else {
		fmt.Println("PopError")
	}
}

func Test_ConcurrenyOperation(t *testing.T) {
	Q.InitQueue(1000, "test")
	wg.Add(2)
	go func() {
		for i := 0; i < MAX_TEST_COUNT; i++ {
			Create()
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_TEST_COUNT; i++ {
			Consum()
		}
		wg.Done()
	}()

	wg.Wait()
}
