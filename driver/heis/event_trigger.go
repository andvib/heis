package driver

import ("time")

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
	
	timer := time.Now()

	for ; true ; {
		for i := 0 ; i < 4 ; i++{
			if (i != 3) && (ELEV_get_button_signal(0,i) == 1) && !((event.Button == "U") && (event.Floor == i)){

                event.Button = "U"
                event.Floor = i
                ButtonChan <- event
				timer = time.Now()

			}else if (i != 0) && (ELEV_get_button_signal(1,i) == 1) && !((event.Button == "D") && (event.Floor == i)){

				event.Button = "D"
                event.Floor = i
                ButtonChan <- event
				timer = time.Now()
    
			}else if (ELEV_get_button_signal(2,i) == 1) && !((event.Button == "C") && (event.Floor == i)){

				event.Button = "C"
                event.Floor = i
                ButtonChan <- event
				timer = time.Now()
			}
        }

		if (time.Since(timer) > 1000*time.Millisecond){
			event.Button = ""
			event.Floor = -1
		}
    }
}
