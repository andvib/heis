package driver

type FloorEvent struct{
	Floor int
	Event string	
}

type ButtonEvent struct{
    Floor int    
    Button string
}

var ElevChan = make(chan FloorEvent,100)
var ButtonChan = make(chan ButtonEvent)

func FloorSensor(){
	currentFloor := -1
	var floor int
    var event FloorEvent

	for ; true ; {
		floor = ELEV_get_floor_sensor_signal()

		if (floor != -1) && (currentFloor != floor){
		    currentFloor = floor
			ELEV_set_floor_indicator(floor)    
		    event.Floor = floor
		    event.Event = "NEW_FLOOR"
		    ElevChan <- event
	    }
	}
}  


func ButtonPush() {
    var event ButtonEvent

	for ; true ; {
		for i := 0 ; i < 4 ; i++{
			if (i != 3) && (ELEV_get_button_signal(0,i) == 1){

   				for ; (ELEV_get_button_signal(0,i)) == 1 ; {
				}
                event.Button = "UP"
                event.Floor = i
                ButtonChan <- event

			}else if (i != 0) && (ELEV_get_button_signal(1,i) == 1){

				for ; (ELEV_get_button_signal(1,i) == 1) ; {
				}
				event.Button = "DOWN"
                event.Floor = i
                ButtonChan <- event
    
			}else if (ELEV_get_button_signal(2,i) == 1){

				for ; (ELEV_get_button_signal(2,i) == 1) ; {
				}
				event.Button = "CMD"
                event.Floor = i
                ButtonChan <- event
			}
        }
    }
}
