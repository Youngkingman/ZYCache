package rpcdef

import (
	"time"
)

const (
	GET = iota
	SET
	SUCCESS
	FAIL
)

//to utilize personal key, we need some protobuf to realize it
type StoreArgs struct {
	Command int
	Key     string
	Value   interface{}
	Expire  time.Duration
}

//the reply value is a string, the client need to assert it
type StoreReply struct {
	Reply int
	Value interface{}
}

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}
