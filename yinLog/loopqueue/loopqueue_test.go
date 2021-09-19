package loopqueue

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
	"fmt"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

const MAX_TEST_COUNT = 200000

func Test_ConcurrenyOperation(t *testing.T) {
	wg.Add(2)
	go func() {
		for i := 0; i < MAX_TEST_COUNT; i++ {
			LogItemPush(DataItem{GET, keystruct.DefaultKey{}, i, time.Minute, time.Now().Unix()})
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_TEST_COUNT; i++ {
			has, item := LogItemPop()
			fmt.Println(has, item)
		}
		wg.Done()
	}()

	wg.Wait()
}
