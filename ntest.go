package main

import (."./driver/network"
	"runtime")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	NETWORK_init()
	
	//select{}

	AppendConn(nil,"78.91.8.288")
	AppendConn(nil,"78.91.8.253")
	AppendConn(nil,"78.91.8.253")

	WhosMaster()
}
