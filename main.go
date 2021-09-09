package main

import (
	keystruct "basic/util/KeyStruct"
	skiplist "basic/util/skipList"
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
	for i := 0; i < 20; i++ {
		key := IntTestKey{keystruct.DefaultKey{}, i}
		sklist.InsertElement(key, "fuck you")
	}
	has, val := sklist.Search(IntTestKey{keystruct.DefaultKey{}, 60})
	fmt.Println(has, val)
	sklist.Show()
}
