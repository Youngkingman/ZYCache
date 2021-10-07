package rpcdef

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func call(rpcname string, args interface{}, reply interface{}, serverAddr string) bool {
	c, err := rpc.DialHTTP("tcp", serverAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

//get some value by cli
func Get(key string, serverAddr string) (interface{}, error) {
	args := StoreArgs{
		Command: GET,
		Key:     key,
		Value:   "",
	}
	reply := StoreReply{}
	call("Coordinator.GetVal", &args, &reply, serverAddr)
	if reply.Reply == SUCCESS {
		return reply.Value, nil
	}
	return "", errors.New("key not found")
}

//set some value by cli
//need value to be marshal
func Set(key string, value string, expire time.Duration, serverAddr string) error {
	args := StoreArgs{
		Command: SET,
		Key:     key,
		Value:   value,
		Expire:  expire,
	}
	reply := StoreReply{}
	call("Coordinator.SetVal", &args, &reply, serverAddr)
	if reply.Reply == SUCCESS {
		return nil
	}
	return errors.New("set fail")
}
