package main

import (."./driver/heis/"
		"runtime"
		"./driver/event/"
		"./driver/queue/"
		"./driver/costfunc/"
		"./driver/network/"
		/*"./driver/phoenix/"*/)

func main(){

	//phoenix.Phoenix()

	runtime.GOMAXPROCS(runtime.NumCPU())

	//Initialize elevator
	if ELEV_init() == 0{
		println("Ikke initialisert!")
	}

	ELEV_set_motor_direction(0)

	queue.Q_init()
	queue.ReadFile()

    //go FloorSensor()
	go ButtonPush()

	go elevlog.ButtonHandle()
	go elevlog.ReceiveMessage()

	event.State = "IDLE"
	go event.StateMachine()
	go event.ReadTimer()
	//go event.ReadEvent()

	//Initialize network
	network.NETWORK_init()

	select {}
}
