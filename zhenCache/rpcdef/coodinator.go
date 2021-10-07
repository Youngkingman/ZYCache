package rpcdef

import (
	"basic/zhenCache/consistenthash"
	store "basic/zhenCache/storeService"
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

const (
	defaultReplicas = 10
)

type Coordinator struct {
	self  string                          //hostname of current serving machine
	port  string                          //port of current serving machine
	peers *consistenthash.ConsisteHashMap //hostname of other serving machines
	mu    sync.Mutex
}

func (c *Coordinator) SetVal(args *StoreArgs, reply *StoreReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if args.Command != SET {
		reply.Reply = FAIL
		return errors.New("WRONG COMMAND")
	}
	//confirm whether the key belongs to me
	peer := c.peers.Get(args.Key)
	if peer == c.self {
		store.SetValue(args.Key, args.Value, args.Expire)
		reply.Reply = SUCCESS
		return nil
	}
	//let other peer handle it
	err := Set(args.Key, args.Value.(string), args.Expire, peer)
	if err == nil {
		reply.Reply = SUCCESS
	}
	return err
}

func (c *Coordinator) GetVal(args *StoreArgs, reply *StoreReply) error {
	if args.Command != GET {
		reply.Reply = FAIL
		return errors.New("WRONG COMMAND")
	}
	//confirm whether the key belongs to me
	peer := c.peers.Get(args.Key)
	if peer == c.self {
		val, err := store.GetValue(args.Key)
		if err != nil {
			reply.Reply = FAIL
			return errors.New("NO KEY")
		}
		reply.Reply = SUCCESS
		reply.Value = val
		return nil
	}
	//let other peer handle it
	val, err := Get(args.Key, peer)
	if err != nil {
		reply.Reply = FAIL
	} else {
		reply.Reply = SUCCESS
		reply.Value = val
	}
	return err
}

func (c *Coordinator) SetPeers(peers ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.peers = consistenthash.New(defaultReplicas, nil)
	c.peers.Add(peers...)
}

func (c *Coordinator) CoodinatorServe() {
	rpc.Register(c)
	rpc.HandleHTTP()
	serverAddr := c.self
	port := c.port
	l, e := net.Listen("tcp", serverAddr+":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func New(serverAddr string, port string, peers []string) *Coordinator {
	c := &Coordinator{
		self: serverAddr,
		port: port,
		mu:   sync.Mutex{},
	}
	c.SetPeers(peers...)
	return c
}
