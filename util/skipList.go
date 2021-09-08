package skiplistwithgo

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

func (node *SkipListNode) setVal(val interface{}) {
	node.val = val
}

func (node *SkipListNode) getVal() interface{} {
	return node.val
}

func (node *SkipListNode) getKey() SListKey {
	return node.key
}

//as an ordered data struct, the key should be able to compare
type SListKey interface {
	CompareBiggerThan(other SListKey) bool
	KeyString() string
	KeyInt() int
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
func (skList *SkipList) getRamdomLevel() (ret int) {
	ret = 1
	//the random algorithm isn't good enough
	rand.Seed(time.Now().UnixNano())
	for rand.Uint32()%2 == 0 {
		ret++
		rand.Seed(time.Now().UnixNano() + int64(ret*(ret<<8))*time.Now().UnixNano())
	}
	if ret >= skList.levelMax {
		ret = skList.levelMax
	}
	return
}

//node initialize function
func (skList *SkipList) createNode(key SListKey, val interface{}, level int) *SkipListNode {
	ret := SkipListNode{
		level:       level,
		key:         key,
		val:         val,
		ForwardList: make([]*SkipListNode, level+1),
	}
	return &ret
}

//if the key already exists, return -1, else insert the node
func (skList *SkipList) InsertElement(key SListKey, val interface{}) int {
	skList.mtx.Lock()
	current := skList.head
	update := make([]*SkipListNode, skList.levelMax+1)

	for i := skList.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.CompareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.CompareBiggerThan(current.ForwardList[i].getKey()) {
			current = current.ForwardList[i]
		}
		update[i] = current
	}
	current = current.ForwardList[0] //reach the 0 level to insert an element

	if current != nil && current.getKey() == key {
		//log.Println("the key is already exists")
		skList.mtx.Unlock()
		return -1
	}

	if current == nil || current.getKey() != key {
		randomLevel := skList.getRamdomLevel()

		if randomLevel > skList.currentLevel {
			for i := skList.currentLevel + 1; i < skList.levelMax; i++ {
				update[i] = skList.head
			}
			skList.currentLevel = randomLevel
		}

		nodeToInsert := skList.createNode(key, val, randomLevel)

		//insert the new node here, update every previous level node
		for i := 0; i < randomLevel; i++ {
			nodeToInsert.ForwardList[i] = update[i].ForwardList[i]
			update[i].ForwardList[i] = nodeToInsert
		}
		skList.elementCount++
		//log.Println("successfully insert a new node")
	}

	skList.mtx.Unlock()
	return 0
}

//if the key already exists, change the value, else insert the node
func (skList *SkipList) UpdateDuplicateKey(key SListKey, val interface{}) {
	skList.mtx.Lock()
	current := skList.head
	update := make([]*SkipListNode, skList.levelMax+1)

	for i := skList.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.compareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.CompareBiggerThan(current.ForwardList[i].getKey()) {
			current = current.ForwardList[i]
		}
		update[i] = current
	}

	current = current.ForwardList[0] //reach the 0 level to insert an element

	if current != nil && current.getKey() == key {
		current.setVal(val)
		skList.mtx.Unlock()
		//log.Println("already have val, sucessfully")
		return
	}

	if current == nil || current.getKey() != key {
		randomLevel := skList.getRamdomLevel()

		if randomLevel > skList.currentLevel {
			for i := skList.currentLevel + 1; i < skList.levelMax; i++ {
				update[i] = skList.head
			}
			skList.currentLevel = randomLevel
		}

		nodeToInsert := skList.createNode(key, val, randomLevel)

		//insert the new node here, update every previous level node
		for i := 0; i < randomLevel; i++ {
			nodeToInsert.ForwardList[i] = update[i].ForwardList[i]
			update[i].ForwardList[i] = nodeToInsert
		}
		skList.elementCount++
		//log.Println("successfully insert a new node")
	}

	skList.mtx.Unlock()
}

//with given key, return if the key exists
func (skList *SkipList) Search(key SListKey) (bool, interface{}) {
	//skList.mtx.RLock()
	current := skList.head
	for i := skList.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.compareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.CompareBiggerThan(current.ForwardList[i].getKey()) {
			current = current.ForwardList[i]
		}
	}
	current = current.ForwardList[0]

	if current != nil && current.key == key {
		//skList.mtx.RUnlock()
		//log.Println("find key")
		return true, current.getVal()
	}
	//skList.mtx.RUnlock()
	//log.Println("key not found")
	return false, nil
}

func (skList *SkipList) Delete(key SListKey) int {
	skList.mtx.Lock()
	current := skList.head
	update := make([]*SkipListNode, skList.levelMax+1)

	for i := skList.currentLevel; i >= 0; i-- {
		//simple change the latter into
		//"current.ForwardList[i].GetKey.compareBiggerThan(key)"
		//can change the order
		for current.ForwardList[i] != nil && key.CompareBiggerThan(current.ForwardList[i].getKey()) {
			current = current.ForwardList[i]
		}
		update[i] = current
	}

	current = current.ForwardList[0] //reach the 0 level to insert an element
	if current != nil && current.getKey() == key {
		//cut the relation ship between current and previous node from lower level to higher one
		for i := 0; i < current.level; i++ {
			// already find the highest level of previous node
			if update[i] == nil || update[i].ForwardList[i] != current {
				break
			}
			update[i].ForwardList[i] = current.ForwardList[i]
		}

		//adjust new level of the skipList
		for skList.currentLevel > 0 && skList.head.ForwardList[current.level] == nil {
			skList.currentLevel--
		}

		skList.elementCount--
		skList.mtx.Unlock()
		//log.Println("element sucessfully delete")
		return 0
	}
	//log.Println("not such element for delete")
	skList.mtx.Unlock()
	return -1
}

func GetSkipList(maxLevel int) (ret SkipList) {
	ret = SkipList{
		head:         ret.createNode(nil, "I am the head", maxLevel),
		tail:         nil,
		elementCount: 0,
		currentLevel: 0,
		levelMax:     maxLevel,
		mtx:          sync.RWMutex{},
	}
	return
}