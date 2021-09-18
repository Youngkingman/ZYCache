package loopqueue

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

type DataItem struct {
	Commandtype int
	Inner       interface{}
	TimeStamp   int64
}

//for log system buffer
var RingQueueService LoopQueue

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
func (lq *LoopQueue) Push(data DataItem) bool {
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
func (lq *LoopQueue) Pop() (bool, interface{}) {
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
