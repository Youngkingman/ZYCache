package skiplist

import (
	keystruct "basic/util/KeyStruct"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/* ************************************************************************
/********************IMPLEMENT OF SKIPLIST********************************
 ************************************************************************/

//basic element for skiplist
type skipListNode struct {
	level       int
	key         keystruct.KeyStruct
	item        interface{}
	ForwardList []*skipListNode
}

func (node *skipListNode) setItem(item interface{}) {
	node.item = item
}

func (node *skipListNode) getItem() interface{} {
	return node.item
}

func (node *skipListNode) getKey() keystruct.KeyStruct {
	return node.key
}

//main struct of SkipList,
type SkipList struct {
	head         *skipListNode //point to the head of SkipList
	tail         *skipListNode //point to the tail of SkipList, not used now
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
func (skList *SkipList) createNode(key keystruct.KeyStruct, item interface{}, level int) *skipListNode {
	ret := skipListNode{
		level:       level,
		key:         key,
		item:        item,
		ForwardList: make([]*skipListNode, level+1),
	}
	return &ret
}

//if the key already exists, return -1, else insert the node
func (skList *SkipList) InsertElement(key keystruct.KeyStruct, item interface{}) int {
	skList.mtx.Lock()
	current := skList.head
	update := make([]*skipListNode, skList.levelMax+1)

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

		nodeToInsert := skList.createNode(key, item, randomLevel)

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

//if the key already exists, change the itemue, else insert the node
func (skList *SkipList) UpdateDuplicateKey(key keystruct.KeyStruct, item interface{}) {
	skList.mtx.Lock()
	current := skList.head
	update := make([]*skipListNode, skList.levelMax+1)

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
		current.setItem(item)
		skList.mtx.Unlock()
		//log.Println("already have item, sucessfully")
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

		nodeToInsert := skList.createNode(key, item, randomLevel)

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
func (skList *SkipList) Search(key keystruct.KeyStruct) (interface{}, bool) {
	skList.mtx.RLock()
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
		skList.mtx.RUnlock()
		//log.Println("find key")
		return current.getItem(), true
	}
	skList.mtx.RUnlock()
	//log.Println("key not found")
	return nil, false
}

func (skList *SkipList) Delete(key keystruct.KeyStruct) error {
	skList.mtx.Lock()
	current := skList.head
	update := make([]*skipListNode, skList.levelMax+1)

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
		return nil
	}
	//log.Println("not such element for delete")
	skList.mtx.Unlock()
	return errors.New("no such element")
}

//this is to show the member of the list
func (skList *SkipList) Show() {
	fmt.Println("/*****THE SKIPLIST*****/")
	for i := 0; i <= skList.currentLevel; i++ {
		node := skList.head.ForwardList[i]
		for node != nil {
			fmt.Print("|", node.getKey(), ":", node.getItem(), "|")
			node = node.ForwardList[i]
		}
	}
}

//contructor of the SkipList
func New(maxLevel int) (ret SkipList) {
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

//return top n elements
func (skList *SkipList) TopN(count int) (keys []keystruct.KeyStruct, items []interface{}) {
	if count < skList.elementCount {
		for i := 0; i < skList.currentLevel; i++ {
			current := skList.head.ForwardList[i]
			for current.ForwardList[i] != nil {
				current = current.ForwardList[i]
				keys = append(keys, current.key)
				items = append(items, current.item)
			}
		}
	} else {
		tmp := 0
		for i := 0; i < skList.currentLevel; i++ {
			current := skList.head.ForwardList[i]
			for current.ForwardList[i] != nil && tmp < count {
				current = current.ForwardList[i]
				keys = append(keys, current.key)
				items = append(items, current.item)
				tmp++
			}
		}
	}

	return
}

func (skList *SkipList) Range(condition func(interface{}) bool) (keys []keystruct.KeyStruct) {
	skList.mtx.RLock()
	for i := 0; i <= skList.currentLevel; i++ {
		node := skList.head.ForwardList[i]
		for node != nil {
			if condition(node.item) {
				keys = append(keys, node.key)
			}
			node = node.ForwardList[i]
		}
	}
	skList.mtx.RUnlock()
	return
}
