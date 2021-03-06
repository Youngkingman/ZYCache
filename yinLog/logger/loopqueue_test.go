package logger

import (
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

const MAX_TEST_COUNT = Q_LENGTH

func Test_ConcurrenyOperation(t *testing.T) {
	wg.Add(1)
	go func() {
		for i := 0; i < MAX_TEST_COUNT; i++ {
			LogItemPush(DataItem{GET, "", i, 0, time.Now().Unix()})
		}
		wg.Done()
	}()
	wg.Wait()
	time.Sleep(10 * time.Second)
	ShutLog()
}
