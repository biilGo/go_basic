package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	rpc_objects "mywork/09.用rpc实现远程过程调用/rpc_objects"
)

func main() {
	calc := new(rpc_objects.Args)
	rpc.Register(calc)
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", "localhost:1234")
	if e != nil {
		log.Fatal("Starting RPC-serve -listen error:", e)
	}
	go http.Serve(listener, nil)
	time.Sleep(1000e9)
}
