package main

import (."./driver/heis/"
	/*"fmt"*/
	"runtime"
	"./driver/event/"
	"./driver/ko/"
	"./driver/costfunc/"
	"time")

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	if ELEV_init() == 0{
		println("Ikke initialisert!")
	}

	ELEV_set_motor_direction(0)
	
    go FloorSensor()
	go ButtonPush()

	event.State = "IDLE"
	go event.StateMachine()

	ko.Q_init()
	go ko.ButtonHandle()
	go event.ReadEvent()

	//go costfunc.ButtonHandle()
	time.Sleep(10000*time.Millisecond)
	costfunc.Cost()
	select {}
	/*for ; true ; {
		costfunc.Cost()
		time.Sleep(1000*time.Millisecond)
	}*/
}
