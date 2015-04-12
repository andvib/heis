package event_manager

import("./driver/heis/")

var State string //IDLE, MOVING, DOOR_OPEN
var Event string //NEW_ORDER, ORDERED_FLOOR_REACHED, EMPTY_QUEUE
var Floor int
var NextFloor int

func stateMachine(){
	switch State {
	
	case "IDLE" :
		if Event == "NEW_ORDER"{	
			if NextFloor == Floor {
				State = "DOOR_OPEN"
				heis.ELEV_set_door_open_lamp(1)
				return
			}else{
				State = "MOVING"
				//Function for calculating dir and
				//and setting direction
				return
			}
		}

	case "MOVING" :
		if Event == ORDERED_FLOOR_REACHED {
			//Stop motor
			State = "DOOR_OPEN"
			heis.ELEV_set_door_open_lamp(1)
			return
		}
		return

	case "DOOR_OPEN" :
		
				


}

/*type event
const (
	floorReachedEvent
	buttonPressedEvent
	
)

func stateMachine(Event event){
	switch(Event){
		case floorReachedEvent{
			handleFloorReachedEvent()
		}
		case buttonPressedEvent{
			handleButtonPressedEvent()
		}
	}
}

func handleFloorReachedEvent(){


}

func handleButtonPressedEvent(){

}
*/
