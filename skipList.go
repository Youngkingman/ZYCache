package skipList

import (
	"math/rand"
	"sync"
	"time"
)

type SkipListNode struct {
	Level       int
	Key         SListKey
	Val         interface{}
	ForwardList []*SkipListNode
}

func (this *SkipListNode) SetVal(val interface{}) {
	this.Val = val
}

func (this *SkipListNode) GetVal() interface{} {
	return this.Val
}

//as an ordered data struct, the key should be able to compare
type SListKey interface {
	compareBiggerThan(other SListKey) bool
}

//only for stress test
type SimpleKey struct {
	myKey int
}

func (this *SimpleKey) compareBiggerThan(other SimpleKey) bool {
	return this.myKey > other.myKey
}

type SkipList struct {
	head         *SkipListNode
	tail         *SkipListNode
	elementCount int
	currentLevel int
	levelMax     int
	mtx          sync.RWMutex
}

func (this *SkipList) getRamdomLevel() int {
	k := 1
	rand.Seed(time.Now().Unix())
	for rand.Uint32()%2 == 0 {
		k++
	}
	if k >= this.levelMax {
		k = this.levelMax
	}
	return k
}

func (this *SkipList) createNode(key SListKey, val interface{}, level int) *SkipListNode {
	ret := SkipListNode{
		Level:       level,
		Key:         key,
		Val:         val,
		ForwardList: make([]*SkipListNode, level+1),
	}
	return &ret
}

func (this *SkipList) InsertElement(key SListKey, val interface{}) int {
	this.mtx.Lock()
	//current := this.head

	this.mtx.Unlock()
	return 0
}

func (this *SkipList) Search(key SListKey) bool {
	this.mtx.RLock()

	this.mtx.RUnlock()
	return false
}

func (this *SkipList) Delete(key SListKey) int { //mut.Lock()
	this.mtx.Lock()

	this.mtx.Unlock()
	return 0
}
