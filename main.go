package main

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
	"basic/zhenCache/rpcdef"
	store "basic/zhenCache/storeService"
	"fmt"
)

func main() {
	//server code
	cood := rpcdef.Coordinator{}
	go cood.Serve()
	//client code
	for i := 0; i < 10000; i++ {
		key := keystruct.DefaultKey{}
		rpcdef.Set(key, i, store.DefaultDuration)
	}
	val, err := rpcdef.Get(keystruct.DefaultKey{})
	fmt.Print(val, err)
}
