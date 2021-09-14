package rbtree

import (
	keystruct "basic/util/KeyStruct"
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

func (rbt *RBTree) Search(key keystruct.KeyStruct) (bool, interface{}) {
	ret := rbt.search(key)
	if ret == rbt._NIL {
		return false, nil
	}
	return true, ret.val
}

func (rbt *RBTree) Delete(key keystruct.KeyStruct) int {
	ret := rbt.delete(key)
	if ret == rbt._NIL {
		return -1
	}
	return 0
}

func (rbt *RBTree) Range() {

}
