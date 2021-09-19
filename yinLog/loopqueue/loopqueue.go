package loopqueue

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
	"sync"
	"time"
)

const (
	GET = iota
	SET
	RANGE
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
	Commandtype int
	Key         keystruct.KeyStruct
	Value       interface{}
	Expire      time.Duration
	TimeStamp   int64
}

//for log system buffer,private queue
var ringQueueService LoopQueue
var lqonce sync.Once

//buffer area of queue
const Q_LENGTH = 4096

func LogItemPush(data DataItem) bool {
	return getService().push(data)
}

func LogItemPop() (bool, interface{}) {
	return getService().pop()
}

func getService() *LoopQueue {
	lqonce.Do(func() {
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
