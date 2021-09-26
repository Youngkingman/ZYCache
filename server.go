package main

import "basic/zhenCache/rpcdef"

func main() {
	//server code
	cood := new(rpcdef.Coordinator)
	cood.CoodinatorServe()
	for {
	}
}
