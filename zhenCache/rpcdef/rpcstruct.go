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
	Value   string
	Expire  time.Duration
}

type StoreReply struct {
	Reply int
	Value string
}

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}
