package main

import (."./driver/heis/"
	/*"fmt"*/)

func main(){
	/*println(LAMP_CHANNEL[1][2])*/
	//println("hei")
	//Elev_set_floor_indicator(1)

	if ELEV_init() == 0{
		println("Ikke initialisert!")
	}
	
	//ELEV_set_button_lamp(0, 2, 1)

	//ELEV_set_stop_lamp(1)

	/*for ; true ; {
		println(ELEV_get_button_signal(BUTTON_CALL_UP, 2))
	}*/

	ELEV_set_motor_direction(DIRN_UP)

	for ; true ; {
		//println(ELEV_get_floor_sensor_signal())
		if ELEV_get_floor_sensor_signal() == N_FLOORS - 1{
			ELEV_set_motor_direction(DIRN_DOWN)
			//println(ELEV_get_floor_sensor_signal())
		} else if ELEV_get_floor_sensor_signal() == 0{
			ELEV_set_motor_direction(DIRN_UP)
		}

		if ELEV_get_stop_signal() == 1 {
			ELEV_set_motor_direction(DIRN_STOP)
		}
	}
}
