package lrucache

import "sync"

type BiList struct {
	Key  string
	Val  interface{}
	Prev *BiList
	Next *BiList
}

type LRUCache struct {
	Head     *BiList
	Tail     *BiList
	Cache    map[string]*BiList
	Size     int
	Capacity int
	mu       sync.Mutex
}

func New(capacity int) (ret LRUCache) {
	head := &BiList{}
	tail := &BiList{}
	head.Next = tail
	tail.Prev = head
	ret.Head, ret.Tail, ret.Size, ret.Capacity = head, tail, 0, capacity
	ret.Cache = make(map[string]*BiList)
	return
}

func (lru *LRUCache) Search(key string) (interface{}, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	val := lru.get(key)
	if val != nil {
		return val, true
	}
	return nil, false
}

func (lru *LRUCache) InsertElement(key string, value interface{}) int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	_, has := lru.Cache[key]
	if has {
		return -1
	}
	lru.put(key, value)
	return 0
}

func (lru *LRUCache) UpdateDuplicateKey(key string, value interface{}) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if _, has := lru.Cache[key]; !has {
		lru.put(key, value)
	} else {
		tmp := lru.Cache[key]
		tmp.Val = value
		lru.moveToHead(tmp)
	}
}

func (lru *LRUCache) get(key string) interface{} {
	if _, has := lru.Cache[key]; !has {
		return nil
	}
	tmp := lru.Cache[key]
	lru.moveToHead(tmp)
	return tmp.Val
}

func (lru *LRUCache) put(key string, value interface{}) {
	node := &BiList{
		key,
		value,
		nil,
		nil,
	}
	lru.Cache[key] = node
	lru.addToHead(node)
	lru.Size++
	if lru.Size > lru.Capacity {
		tmp := lru.Tail.Prev
		lru.trimFromList(tmp)
		delete(lru.Cache, tmp.Key)
		lru.Size--
	}
}

func (lru *LRUCache) addToHead(node *BiList) {
	node.Next = lru.Head.Next
	node.Prev = lru.Head
	lru.Head.Next.Prev = node
	lru.Head.Next = node
}
func (lru *LRUCache) trimFromList(node *BiList) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	node.Next, node.Prev = nil, nil
}
func (lru *LRUCache) moveToHead(node *BiList) {
	lru.trimFromList(node)
	lru.addToHead(node)
}
