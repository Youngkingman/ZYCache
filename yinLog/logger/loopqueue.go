package logger

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
	"sync"
)

const (
	GET = iota
	SET
	RANGE
	INITMESSAGE
)

//support for log request
type LoopQueue struct {
	start  int
	end    int
	length int
	name   string
	data   []DataItem
}

//command log item
type DataItem struct {
	Commandtype int                 `json:"commandtype"`
	Key         keystruct.KeyStruct `json:"key"`
	Value       interface{}         `json:"value"`
	Expire      int64               `json:"duration"`
	TimeStamp   int64               `json:"log_time"`
}

//for log system buffer,private queue
var ringQueueService LoopQueue

//init operation
var lqonce sync.Once

//buffer area of queue
const Q_LENGTH = 32768

//enable the service
func init() {
	getQueue()
}

//user interface for log push
func LogItemPush(data DataItem) bool {
	return getQueue().push(data)
}

//user interface for log pop
func LogItemPop() (bool, interface{}) {
	return getQueue().pop()
}

//if you have ever used LogItemPush or LogItemPop
//please make sure that you will call ShutLog() before
//you close the servise
var shutdown chan struct{}

func ShutLog() {
	shutdown = make(chan struct{})
	lt.stop <- struct{}{}
	<-shutdown
}

func getQueue() *LoopQueue {
	lqonce.Do(func() {
		go startLogAppendServe()
		q := LoopQueue{
			start:  0,
			end:    0,
			length: Q_LENGTH,
			name:   "log_buffer",
			data:   make([]DataItem, Q_LENGTH),
		}
		ringQueueService = q
	})
	return &ringQueueService
}

func (lq *LoopQueue) InitQueue(length int, name string) bool {
	if nil == lq || length <= 0 {
		return false
	}
	lq.data = make([]DataItem, length)
	lq.length = length
	lq.name = name
	lq.start = 0
	lq.end = 0
	return true
}

func (lq *LoopQueue) push(data DataItem) bool {
	if nil == lq {
		panic("LoopQueue is nil")
	}
	if lq.isFull() {
		return false
	}
	var end int = lq.getEnd()
	lq.data[end] = data
	lq.end = (end + 1) % lq.length
	return true
}

func (lq *LoopQueue) pop() (bool, interface{}) {
	if nil == lq {
		panic("LoopQueue is nil")
	}
	if lq.isEmpty() {
		return false, nil
	}
	var start = lq.getStart()
	var startValue interface{} = lq.data[start]
	lq.start = (start + 1) % lq.length
	return true, startValue
}

func (lq *LoopQueue) isEmpty() bool {
	if nil == lq {
		panic("LoopQueue is nil")
	}
	if lq.getStart() == lq.getEnd() {
		return true
	}
	return false
}

func (lq *LoopQueue) isFull() bool {
	if nil == lq {
		panic("LoopQueue is nil")
	}
	if lq.getEnd()+1 == lq.getStart() {
		return true
	}
	return false
}

func (lq *LoopQueue) getStart() int {
	return lq.start % lq.length
}

func (lq *LoopQueue) getEnd() int {
	return lq.end % lq.length
}
