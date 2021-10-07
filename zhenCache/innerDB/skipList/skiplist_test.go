package skiplist

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

//Implement of SlistKey

var wg sync.WaitGroup
var skList SkipList

const MAX_OPERATION = 20000 //the max number of single CRUD operation
const LEVEL_COUNT = 32      //the level number of SkipList
const TOP_COUNT = 20030

func Test_ConcurrenyOperation(t *testing.T) {
	skList = New(LEVEL_COUNT)
	rand.Seed(time.Now().UnixNano())
	//keySlice := make([]IntTestKey, 200000)
	wg.Add(4)
	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			skList.InsertElement(key, "fuck you")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			//keySlice[i] = key
			skList.UpdateDuplicateKey(key, "fuck you")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			skList.Search(key)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := string(rune(int(rand.Uint32())))
			skList.Delete(key)
		}
		wg.Done()
	}()

	wg.Wait()
}

func Test_BasciFunctionTest(t *testing.T) {
	skList = New(LEVEL_COUNT)
	rand.Seed(time.Now().UnixNano())
	keySlice := make([]string, MAX_OPERATION)
	for i := 0; i < MAX_OPERATION; i++ {
		key := string(rune(int(rand.Uint32())))
		skList.InsertElement(key, "fuck you")
	}
	for i := 0; i < MAX_OPERATION; i++ {
		key := string(rune(int(rand.Uint32())))
		keySlice[i] = key
		skList.UpdateDuplicateKey(key, "fuck you")
	}
	for i := 0; i < TOP_COUNT; i++ {
		skList.TopN(i)
	}
	for i := 0; i < MAX_OPERATION; i++ {
		key := string(rune(int(rand.Uint32())))
		skList.Search(key)
	}
	for i := 0; i < MAX_OPERATION; i++ {
		skList.Delete(keySlice[i])
	}
	skList.Show()
}
