package driver

//import("./driver/heis/")

type Event struct{
	floor int
	event string	
}

var C = make(chan *Event)

/*type floorReachedEvent struct{
	floor int
}

type buttonPressedEvent struct{
	floor int
	button elev_button_type
}*/

func FloorSensor(){
	currentFloor := -1
	var floor int
	for ; true ; {
        floor = ELEV_get_floor_sensor_signal()
        if (floor != -1) && (currentFloor != floor){
            currentFloor = floor
			ELEV_set_floor_indicator(floor)
            var event Event
            event.floor = floor
            event.event = "NEW_FLOOR"
            C <- &event
        }
    }
}  


func ButtonPush() {
	for ; true ; {
		for i := 0 ; i < 4 ; i++{
			if (i != 3) && (ELEV_get_button_signal(0,i)){
				for ; (ELEV_get_button_signal(0,i)) ; {
				}
				println("Knappetrykk: OPP, ", i)
			}else if (i != 0) && (ELEV_get_button_signal(1,i)){
				for ; (ELEV_get_button_signal(1,i)) ; {
				}
				println("Knappetrykk: NED, ", i)
			}else if (ELEV_get_button_signal(2,i)){
				for ; (ELEV_get_button_signal(2,i)) {
				}
				println("Knappetrykk: INNE, ", i)
			}
}
	



/*func buttonPoll(button elev_button_type, event chan) {
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
}*/
