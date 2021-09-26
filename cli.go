package main

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
	"basic/zhenCache/rpcdef"
	store "basic/zhenCache/storeService"
	"fmt"
)

func main() {
	key := keystruct.StringKey{"fuck you"}
	for i := 0; i < 5; i++ {
		rpcdef.Set(key, "man", store.DefaultDuration)
	}
	val, err := rpcdef.Get(key)
	fmt.Println(val, err)
}
