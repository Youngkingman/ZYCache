package rpcdef

import (
	keystruct "basic/zhenCache/innerDB/keystruct"
	"time"
)

const (
	GET = iota
	SET
	SUCCESS
	FAIL
)

type StoreArgs struct {
	Command int
	Key     keystruct.KeyStruct
	Value   interface{}
	Expire  time.Duration
}

type StoreReply struct {
	Reply int
	Key   keystruct.KeyStruct
	Value interface{}
}
