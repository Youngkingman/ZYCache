package rpcdef

import (
	store "basic/zhenCache/storeService"
	"encoding/json"
	"fmt"
)

type FUCK struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z int64  `json:"z"`
}

func Fuck_Client_example() {
	test := FUCK{"sdf", "sdfsde", 65465654}
	//currently need to marshal
	//later will introduce protobuf and msg protocal
	a, _ := json.Marshal(test)
	for i := 0; i < 5; i++ {
		//rpcdef.Set("fuck", string(a), store.DefaultDuration)
		Set("fuck", string(a), store.DefaultDuration, "127.0.0.1:1234")
	}
	//val, _ := rpcdef.Get("fuck")
	val, _ := Get("fuck", "127.0.0.1:1234")
	rettest := FUCK{}
	json.Unmarshal([]byte(val.(string)), &rettest)
	fmt.Println(rettest)
}

func Server_example() {
	//cood := new(rpcdef.Coordinator)
	cood := new(Coordinator)
	cood.CoodinatorServe()
	select {}
}
