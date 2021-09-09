package rbtree

import (
	keystruct "basic/util/KeyStruct"
	"sync"
)

const (
	RED   = 0
	BLACK = 1
)

type RBTreeNode struct {
	Left   *RBTreeNode
	Right  *RBTreeNode
	Parent *RBTreeNode
	Color  uint

	key keystruct.KeyStruct
	val interface{}
}

/* 5 Attributes of RBTree
1.the node is either red or black
2.the root node is black
3.every leaves node is black(and they have nil values)
4.if a node is red, its children are black
5.every route from node to its successors poccesses same counts of black nodes

*/
type RBTree struct {
	root         *RBTreeNode
	_NIL         *RBTreeNode
	elementCount int
	mtx          sync.RWMutex
}

func (rbt *RBTree) leftRotate(x *RBTreeNode) {
	if x.Right == rbt._NIL {
		return
	}
	y := x.Right
	x.Right = y.Left
	if y.Left != rbt._NIL {
		y.Left.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == rbt._NIL {
		rbt.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

func (rbt *RBTree) rightRotate(x *RBTreeNode) {
	if x.Left == rbt._NIL {
		return
	}
	y := x.Left
	x.Left = y.Right
	if y.Right != rbt._NIL {
		y.Right.Parent = x
	}
	y.Parent = x.Parent

	if x.Parent == rbt._NIL {
		rbt.root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Right = x
	x.Parent = y
}

func (rbt *RBTree) insert(z *RBTreeNode) *RBTreeNode {
	//step 1: insert the node like normal binary search tree
	//left is small and right is big
	x := rbt.root
	y := rbt._NIL

	for x != rbt._NIL {
		y = x
		if x.key.CompareBiggerThan(z.key) {
			x = x.Left
		} else if z.key.CompareBiggerThan(x.key) {
			x = x.Right
		} else {
			return x
		}
	}

	z.Parent = y
	if y == rbt._NIL {
		//there is not root
		rbt.root = z
	} else if y.key.CompareBiggerThan(z.key) {
		y.Left = z
	} else {
		y.Right = z
	}
	rbt.elementCount++
	rbt.fixAfterInsert(z)
	return z
}

func (rbt *RBTree) fixAfterInsert(z *RBTreeNode) {
	//new node z is a red node
	for z.Parent.Color == RED {
		if z.Parent == z.Parent.Parent.Left {
			//y is an uncle node
			y := z.Parent.Parent.Right
			if y.Color == RED {
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {

				}
			}
		}
	}
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
