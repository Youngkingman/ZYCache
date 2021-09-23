package logger

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
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
			LogItemPush(DataItem{GET, keystruct.DefaultKey{}, i, time.Minute, time.Now().Unix()})
		}
		wg.Done()
	}()
	wg.Wait()
	ShutLog()
}
