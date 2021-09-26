package rpcdef

import (
	store "basic/zhenCache/storeService"
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Coordinator struct{}

func (c *Coordinator) SetVal(args *StoreArgs, reply *StoreReply) error {
	if args.Command != SET {
		reply.Reply = FAIL
		return errors.New("WRONG COMMAND")
	}
	//currently only support string
	store.SetValue(args.Key, args.Value, args.Expire)
	reply.Reply = SUCCESS
	return nil
}

func (c *Coordinator) GetVal(args *StoreArgs, reply *StoreReply) error {
	if args.Command != GET {
		reply.Reply = FAIL
		return errors.New("WRONG COMMAND")
	}
	val, err := store.GetValue(args.Key)
	if err != nil {
		reply.Reply = FAIL
		return errors.New("NO KEY")
	}
	reply.Reply = SUCCESS
	reply.Value = val
	return nil
}

func (c *Coordinator) CoodinatorServe() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")

	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//some example
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}
