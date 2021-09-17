package rbtree

import (
	keystruct "basic/util/KeyStruct"
	"errors"
	"sync"
)

func (rbt *RBTree) InsertElement(key keystruct.KeyStruct, val interface{}) {
	node := RBTreeNode{
		rbt._NIL,
		rbt._NIL,
		rbt._NIL,
		RED,
		key,
		val,
	}
	rbt.insert(&node, false)
}

func (rbt *RBTree) UpdateDuplicateKey(key keystruct.KeyStruct, val interface{}) {
	node := RBTreeNode{
		rbt._NIL,
		rbt._NIL,
		rbt._NIL,
		RED,
		key,
		val,
	}
	rbt.insert(&node, true)
}

func (rbt *RBTree) Search(key keystruct.KeyStruct) (interface{}, bool) {
	ret := rbt.search(key)
	if ret == rbt._NIL {
		return nil, false
	}
	return ret.item, true
}

func (rbt *RBTree) Delete(key keystruct.KeyStruct) error {
	ret := rbt.delete(key)
	if ret == rbt._NIL {
		return errors.New("no such element")
	}
	return nil
}

func (rbt *RBTree) Range(condition func(interface{}) bool) (keys []keystruct.KeyStruct) {
	rbt.mtx.RLock()
	rbt.preOreder(rbt.root, condition, keys)
	rbt.mtx.RUnlock()
	return
}

func New() (rbt RBTree) {
	initNode := &RBTreeNode{nil, nil, nil, BLACK, keystruct.DefaultKey{}, nil}
	rbt = RBTree{
		root:         initNode,
		_NIL:         initNode,
		elementCount: 0,
		mtx:          sync.RWMutex{},
	}
	return
}
