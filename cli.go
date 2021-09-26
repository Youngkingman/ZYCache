package main

import (
	"basic/zhenCache/rpcdef"
	store "basic/zhenCache/storeService"
	"fmt"
)

func main() {
	for i := 0; i < 5; i++ {
		rpcdef.Set("fuck", "man", store.DefaultDuration)
	}
	val, err := rpcdef.Get("fuck")
	fmt.Println(val, err)
}
