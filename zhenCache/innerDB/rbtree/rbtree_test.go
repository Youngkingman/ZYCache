package rbtree

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup
var rbt RBTree

const MAX_OPERATION = 600 //the max number of single CRUD operation

func Test_ConcurrenyOperation(t *testing.T) {
	rbt = New()
	rand.Seed(time.Now().UnixNano())
	//keySlice := make([]IntTestKey, 200000)
	wg.Add(4)
	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			key1 := string(rune(i))
			rbt.InsertElement(key, "fuck you")
			rbt.InsertElement(key1, "FUCK YOU")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			//keySlice[i] = key
			rbt.UpdateDuplicateKey(key, "fuck you")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			rbt.Search(key)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			rbt.Delete(key)
		}
		wg.Done()
	}()

	wg.Wait()
}

func Test_BasciFunctionTest(t *testing.T) {
	rbt = New()
	rand.Seed(time.Now().UnixNano())
	keySlice := make([]string, MAX_OPERATION)
	for i := 0; i < MAX_OPERATION; i++ {
		key := string(rune(int(rand.Uint32())))
		rbt.InsertElement(key, "fuck you")
	}
	for i := 0; i < MAX_OPERATION; i++ {
		key := string(rune(int(rand.Uint32())))
		keySlice[i] = key
		rbt.UpdateDuplicateKey(key, "fuck you")
	}
	for i := 0; i < MAX_OPERATION; i++ {
		key := string(rune(int(rand.Uint32())))
		rbt.Search(key)
	}
	for i := MAX_OPERATION - 1; i >= MAX_OPERATION/2; i-- {
		rbt.Delete(keySlice[i])
	}
	for i := 0; i < MAX_OPERATION/2; i++ {
		rbt.Delete(keySlice[i])
	}
}
