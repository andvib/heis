package main

import (."./driver/heis/"
	/*"fmt"*/
	"runtime"
	"./driver/event/"
	"./driver/ko/"
	"./driver/costfunc/"
	/*"time"*/
	"./driver/network/")

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	if ELEV_init() == 0{
		println("Ikke initialisert!")
	}

	network.NETWORK_init()

	ELEV_set_motor_direction(0)
	
    go FloorSensor()
	go ButtonPush()
	go costfunc.ButtonHandle()
	go costfunc.ReceiveMessage()

	event.State = "IDLE"
	go event.StateMachine()

	ko.Q_init()
	//go ko.ButtonHandle()
	go event.ReadEvent()

	select {}
}
