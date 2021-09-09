package skiplistwithgo

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

//Implement of SlistKey
type IntTestKey struct {
	ParameterListSkey
	key int
}

func (key IntTestKey) CompareBiggerThan(other SListKey) bool {
	return key.key > other.KeyInt32()
}

func (key IntTestKey) KeyInt() int {
	return key.key
}

var wg sync.WaitGroup
var skList SkipList

const MAX_OPERATION = 300000

func Test_ConcurrenyOperation(t *testing.T) {
	skList = GetSkipList(5)
	rand.Seed(time.Now().UnixNano())
	//keySlice := make([]IntTestKey, 200000)
	wg.Add(4)
	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := IntTestKey{ParameterListSkey{}, int(rand.Uint32())}
			skList.InsertElement(key, "fuck you")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := IntTestKey{ParameterListSkey{}, int(rand.Uint32())}
			//keySlice[i] = key
			skList.UpdateDuplicateKey(key, "fuck you")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := IntTestKey{ParameterListSkey{}, int(rand.Uint32())}
			skList.Search(key)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < MAX_OPERATION; i++ {
			key := IntTestKey{ParameterListSkey{}, int(rand.Uint32())}
			skList.Delete(key)
		}
		wg.Done()
	}()

	wg.Wait()
}

func Test_BasciFunctionTest(t *testing.T) {
	skList = GetSkipList(5)
	rand.Seed(time.Now().UnixNano())
	keySlice := make([]IntTestKey, 20000)
	for i := 0; i < 20000; i++ {
		key := IntTestKey{ParameterListSkey{}, int(rand.Uint32())}
		skList.InsertElement(key, "fuck you")
	}
	for i := 0; i < 20000; i++ {
		key := IntTestKey{ParameterListSkey{}, int(rand.Uint32())}
		keySlice[i] = key
		skList.UpdateDuplicateKey(key, "fuck you")
	}
	for i := 0; i < 20000; i++ {
		key := IntTestKey{ParameterListSkey{}, int(rand.Uint32())}
		skList.Search(key)
	}
	for i := 0; i < 20000; i++ {
		skList.Delete(keySlice[i])
	}
}
