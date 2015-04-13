package main

import (."./driver/heis/"
	/*"fmt"*/
		"runtime")

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())

	if ELEV_init() == 0{
		println("Ikke initialisert!")
	}
	
    //go FloorSensor()
	go ButtonPush()

	for ; true ; {
        hendelse := <- ButtonChan
        
        println("BUTTON: ", hendelse.Button)
        println("FLOOR: ", hendelse.Floor)
		
		
		/*if ELEV_get_floor_sensor_signal() == N_FLOORS - 1{
		    ELEV_set_motor_direction(DIRN_DOWN)
			//println(ELEV_get_floor_sensor_signal())
		} else if ELEV_get_floor_sensor_signal() == 0{
			ELEV_set_motor_direction(DIRN_UP)
		}

		if ELEV_get_stop_signal() == 1 {
			ELEV_set_motor_direction(DIRN_STOP)
		}
	*/		
	}
}
