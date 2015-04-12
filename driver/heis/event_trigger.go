package driver

type Event struct{
	floor int
	event string	
}

var c = make(chan Event, 1)

/*type floorReachedEvent struct{
	floor int
}

type buttonPressedEvent struct{
	floor int
	button elev_button_type
}*/


func buttonPoll(button elev_button_type, event chan) {
	alreadyPushed := [N_FLOORS]int{0,0,0,0}
	for j:=0; j < N_FLOORS-1; j++{
		buttonPushed := ELEV_get_button_signal(button, j)
		if buttonPushed == 1 && alreadyPushed[j] == 0{
			event <- buttonPressedEvent(j,button)
		}
		alreadyPushed = buttonPushed[j]	
		
	}
}

func floorSensorPoll(event chan) 
	currentFloor := ELEV_get_floor_sensor_signal();
	for j:=0; j < N_FLOORS-1; j++{
		floor := ELEV_get_floor_sensor_signal()
		if floor != currentFloor{
			event <- floorReachedEvent(j)
		currentFloor = floor
	}
}
