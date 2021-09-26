package main

import "basic/zhenCache/rpcdef"

func main() {
	//server code
	cood := new(rpcdef.Coordinator)
	cood.CoodinatorServe()
	// // //client code
	// // // for i := 0; i < 5; i++ {
	// // // 	key := TestKey{keystruct.DefaultKey{}, "fuck you"}
	// // // 	rpcdef.Set(key, i, store.DefaultDuration)
	// // // }
	// // // val, err := rpcdef.Get(keystruct.DefaultKey{})
	// // // fmt.Print(val, err)
	for {
	}
}
