package rpcdef

import (
	store "basic/zhenCache/storeService"
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type Coordinator struct{}

func (c *Coordinator) Set(args *StoreArgs, reply *StoreReply, expire time.Duration) error {
	if args.Command != SET {
		reply.Reply = FAIL
		return errors.New("WRONG COMMAND")
	}
	store.SetValue(args.Key, args.Value, expire)
	reply.Reply = SUCCESS
	return nil
}

func (c *Coordinator) Get(args *StoreArgs, reply *StoreReply) error {
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

func (c *Coordinator) Serve() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")

	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
