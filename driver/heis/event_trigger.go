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


/*func FloorSensor(){
	//Reads the floor-sensors and sends new-floor events to event manager
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
}*/ 


func ButtonPush() {
	//Reads when a button is pushed and sends the event to the elevlogic-module
	//Uses a timer to avoid spamming the buttons
    var buttonEvent ButtonEvent

	currentFloor := -1
	var floor int
    var floorEvent FloorEvent


	
	timer := time.Now()

	for ; true ; {
		for i := 0 ; i < 4 ; i++{
			if (i != 3) && (ELEV_get_button_signal(0,i) == 1) && !((buttonEvent.Button == "U") && (buttonEvent.Floor == i)){

                buttonEvent.Button = "U"
                buttonEvent.Floor = i
                ButtonChan <- buttonEvent
				timer = time.Now()

			}else if (i != 0) && (ELEV_get_button_signal(1,i) == 1) && !((buttonEvent.Button == "D") && (buttonEvent.Floor == i)){

				buttonEvent.Button = "D"
                buttonEvent.Floor = i
                ButtonChan <- buttonEvent
				timer = time.Now()
    
			}else if (ELEV_get_button_signal(2,i) == 1) && !((buttonEvent.Button == "C") && (buttonEvent.Floor == i)){

				buttonEvent.Button = "C"
                buttonEvent.Floor = i
                ButtonChan <- buttonEvent
				timer = time.Now()
			}
        }

		if (time.Since(timer) > 1000*time.Millisecond){
			buttonEvent.Button = ""
			buttonEvent.Floor = -1
		}

		floor = ELEV_get_floor_sensor_signal()

		if (floor != -1) && (currentFloor != floor){
		    currentFloor = floor
			ELEV_set_floor_indicator(floor)    
		    floorEvent.Floor = floor
		    floorEvent.Event = "NEW_FLOOR"
		    ElevChan <- floorEvent
	    }
		time.Sleep(500*time.Microsecond)
    }
}
