package main

import (."./driver/network"
		"runtime")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	NETWORK_init()
	
	for ; true ; {
	}
}
