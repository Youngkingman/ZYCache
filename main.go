package main

import (
	skiplistwithgo "basic/util"
	"fmt"
)

type IntTestKey struct {
	key int
}

func (key IntTestKey) CompareBiggerThan(other skiplistwithgo.SListKey) bool {
	return key.key > other.KeyInt()
}

func (key IntTestKey) KeyInt() int {
	return key.key
}

func (key IntTestKey) KeyString() (ret string) {
	return
}

func main() {
	sklist := skiplistwithgo.GetSkipList(5)
	for i := 0; i < 2000; i++ {
		key := IntTestKey{i}
		sklist.InsertElement(key, "fuck you")
	}

	has, val := sklist.Search(IntTestKey{60})
	fmt.Println(has, val)
}
