package skipList

import (
	"math/rand"
	"sync"
	"time"
)

//basic element of a SkipList
type SkipListNode struct {
	level       int
	key         SListKey
	val         interface{}
	ForwardList []*SkipListNode
}

func (this *SkipListNode) SetVal(val interface{}) {
	this.val = val
}

func (this *SkipListNode) GetVal() interface{} {
	return this.val
}

func (this *SkipListNode) GetKey() SListKey {
	return this.key
}

//as an ordered data struct, the key should be able to compare
type SListKey interface {
	compareBiggerThan(other SListKey) bool
}

//only for stress test
type SimpleKey struct {
	myKey int
}

//implement of SlistKey interface
func (this *SimpleKey) compareBiggerThan(other SimpleKey) bool {
	return this.myKey > other.myKey
}

//main struct of SkipList,
type SkipList struct {
	head         *SkipListNode //point to the head of SkipList
	tail         *SkipListNode //point to the tail of SkipList, not used now
	elementCount int           //update during every insert operation
	currentLevel int           //current level of the SkipList
	levelMax     int           //initialized to limit the max level of SkipList
	mtx          sync.RWMutex  //gurantee the safety of multi gorotiues operation
}

//give a random level for a node, bigger level is less possible to occur
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

//node initialize function
func (this *SkipList) createNode(key SListKey, val interface{}, level int) *SkipListNode {
	ret := SkipListNode{
		level:       level,
		key:         key,
		val:         val,
		ForwardList: make([]*SkipListNode, level+1),
	}
	return &ret
}

//if the key already exists, return -1, else insert the node
func (this *SkipList) InsertElement(key SListKey, val interface{}) int {
	this.mtx.Lock()
	current := this.head
	update := make([]*SkipListNode, this.levelMax+1)

	for i := this.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.compareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.compareBiggerThan(current.ForwardList[i].GetKey()) {
			current = current.ForwardList[i]
		}
		update[i] = current
	}

	current = current.ForwardList[0] //reach the 0 level to insert an element

	if current != nil && current.GetKey() == key {
		//log.Println("the key is already exists")
		this.mtx.Unlock()
		return -1
	}

	if current == nil || current.GetKey() != key {
		randomLevel := this.getRamdomLevel()

		if randomLevel > this.currentLevel {
			for i := this.currentLevel + 1; i < this.levelMax; i++ {
				update[i] = this.head
			}
			this.currentLevel = randomLevel
		}

		nodeToInsert := this.createNode(key, val, randomLevel)

		//insert the new node here, update every previous level node
		for i := 0; i < randomLevel; i++ {
			nodeToInsert.ForwardList[i] = update[i].ForwardList[i]
			update[i].ForwardList[i] = nodeToInsert
		}
		this.elementCount++
		//log.Println("successfully insert a new node")
	}

	this.mtx.Unlock()
	return 0
}

//if the key already exists, change the value, else insert the node
func (this *SkipList) UpdateDuplicateKey(key SListKey, val interface{}) {
	this.mtx.Lock()
	current := this.head
	update := make([]*SkipListNode, this.levelMax+1)

	for i := this.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.compareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.compareBiggerThan(current.ForwardList[i].GetKey()) {
			current = current.ForwardList[i]
		}
		update[i] = current
	}

	current = current.ForwardList[0] //reach the 0 level to insert an element

	if current != nil && current.GetKey() == key {
		current.SetVal(val)
		this.mtx.Unlock()
		//log.Println("already have val, sucessfully")
		return
	}

	if current == nil || current.GetKey() != key {
		randomLevel := this.getRamdomLevel()

		if randomLevel > this.currentLevel {
			for i := this.currentLevel + 1; i < this.levelMax; i++ {
				update[i] = this.head
			}
			this.currentLevel = randomLevel
		}

		nodeToInsert := this.createNode(key, val, randomLevel)

		//insert the new node here, update every previous level node
		for i := 0; i < randomLevel; i++ {
			nodeToInsert.ForwardList[i] = update[i].ForwardList[i]
			update[i].ForwardList[i] = nodeToInsert
		}
		this.elementCount++
		//log.Println("successfully insert a new node")
	}

	this.mtx.Unlock()
	return
}

//with given key, return if the key exists
func (this *SkipList) Search(key SListKey) bool {
	//this.mtx.RLock()
	current := this.head
	for i := this.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.compareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.compareBiggerThan(current.ForwardList[i].GetKey()) {
			current = current.ForwardList[i]
		}
	}
	current = current.ForwardList[0]

	if current != nil && current.key == key {
		//this.mtx.RUnlock()
		//log.Println("find key")
		return true
	}
	//this.mtx.RUnlock()
	//log.Println("key not found")
	return false
}

func (this *SkipList) Delete(key SListKey) int { //mut.Lock()
	this.mtx.Lock()
	current := this.head
	update := make([]*SkipListNode, this.levelMax+1)

	for i := this.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.compareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.compareBiggerThan(current.ForwardList[i].GetKey()) {
			current = current.ForwardList[i]
		}
		update[i] = current
	}

	current = current.ForwardList[0] //reach the 0 level to insert an element
	if current != nil && current.GetKey() == key {
		//cut the relation ship between current and previous node from lower level to higher one
		for i := 0; i < current.level; i++ {
			// already find the highest level of previous node
			if update[i].ForwardList[i] != current {
				break
			}
			update[i].ForwardList[i] = current.ForwardList[i]
		}

		//adjust new level of the skipList
		for this.currentLevel > 0 && this.head.ForwardList[current.level] == nil {
			this.currentLevel--
		}

		this.elementCount--
		this.mtx.Unlock()
		//log.Println("element sucessfully delete")
		return 0
	}
	//log.Println("not such element for delete")
	this.mtx.Unlock()
	return -1
}

func GetSkipList() (ret SkipList) {
	return
}
