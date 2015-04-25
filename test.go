package main

import (."./driver/heis/"
	"runtime"
	"./driver/event/"
	"./driver/ko/"
	"./driver/costfunc/"
	"./driver/network/")

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Initialize elevator
	if ELEV_init() == 0{
		println("Ikke initialisert!")
	}

	ELEV_set_motor_direction(0)

	ko.Q_init()
	ko.ReadFile()

    go FloorSensor()
	go ButtonPush()

	go costfunc.ButtonHandle()
	go costfunc.ReceiveMessage()

	event.State = "IDLE"
	go event.StateMachine()
	go event.ReadTimer()
	go event.ReadEvent()

	//Initialize network
	network.NETWORK_init()


	select {}
}
