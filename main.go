package main

import (
	skiplistwithgo "basic/util/skipList"
	"fmt"
)

type IntTestKey struct {
	skiplistwithgo.ParameterListSkey
	key int
}

func (key IntTestKey) CompareBiggerThan(other skiplistwithgo.SListKey) bool {
	return key.key > other.KeyInt32()
}

func (key IntTestKey) KeyInt32() int {
	return key.key
}

func main() {
	sklist := skiplistwithgo.GetSkipList(5)
	for i := 0; i < 20; i++ {
		key := IntTestKey{skiplistwithgo.ParameterListSkey{}, i}
		sklist.InsertElement(key, "fuck you")
	}
	has, val := sklist.Search(IntTestKey{skiplistwithgo.ParameterListSkey{}, 60})
	fmt.Println(has, val)
	sklist.Show()
}
