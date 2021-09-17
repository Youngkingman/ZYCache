package main

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
	skiplist "basic/zhenCache/innerDB/skipList"
	"fmt"
)

type IntTestKey struct {
	keystruct.DefaultKey
	key int
}

func (key IntTestKey) CompareBiggerThan(other keystruct.KeyStruct) bool {
	return key.key > other.KeyInt32()
}

func (key IntTestKey) KeyInt32() int {
	return key.key
}

func main() {
	sklist := skiplist.New(5)
	tmp := IntTestKey{}
	for i := 0; i < 20; i++ {
		key := IntTestKey{keystruct.DefaultKey{}, i}
		tmp = key
		sklist.InsertElement(key, "fuck you")
	}
	has, val := sklist.Search(tmp)
	fmt.Println(has, val)
	sklist.Show()
}
