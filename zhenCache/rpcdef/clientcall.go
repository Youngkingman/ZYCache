package rpcdef

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func call(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
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
func Get(key string) (interface{}, error) {
	args := StoreArgs{
		Command: GET,
		Key:     key,
		Value:   "",
	}
	reply := StoreReply{}
	call("Coordinator.GetVal", &args, &reply)
	if reply.Reply == SUCCESS {
		return reply.Value, nil
	}
	return "", errors.New("key not found")
}

//set some value by cli
func Set(key string, value string, expire time.Duration) error {
	args := StoreArgs{
		Command: SET,
		Key:     key,
		Value:   value,
		Expire:  expire,
	}
	reply := StoreReply{}
	call("Coordinator.SetVal", &args, &reply)
	if reply.Reply == SUCCESS {
		return nil
	}
	return errors.New("set fail")
}

func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.Example", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}
