package event

import(".././heis/"
		"time"
		".././ko/")

var State string //IDLE, MOVING, DOOR_OPEN
var Event string //NEW_ORDER, ORDERED_FLOOR_REACHED, EMPTY_QUEUE
var EventFloor int
var Floor int
var CurrentFloor
var NextFloor int
var Dir string

func StateMachine(){
	startUp()
	for ; true ; {
		NextFloor = ko.NextInQ(Dir,Floor) 
		switch State {
	
		case "IDLE" :
			if Event == "NEW_ORDER"{
				//println("IDLE:NEW_ORDER")
				NextFloor = ko.NextInQ(Dir,Floor)
				if NextFloor == Floor {
					State = "DOOR_OPEN"
				}else{
					State = "MOVING"
					moveToFloor()
				//	println("MOVING")
				}
			}

		case "MOVING" :
			if Event == "NEW_FLOOR" {
				CurrentFloor := driver.ELEV_get_floor_sensor_signal()
				if (temp != -1){
					Floor = temp
				}
				if (driver.ELEV_get_floor_sensor_signal() == NextFloor){
					driver.ELEV_set_motor_direction(0)
					State = "DOOR_OPEN"
					//println("EVENT:DOOR_OPEN")
				}
			}

		case "DOOR_OPEN" :
			if (driver.ELEV_get_floor_sensor_signal() != -1){
				//println("DOOR OPEN ",Floor)
				driver.ELEV_set_door_open_lamp(1)
				time.Sleep(3000*time.Millisecond)
				driver.ELEV_set_door_open_lamp(0)
				//println("DOOR CLOSED")
				ko.RemoveOrder(Floor)
				NextFloor = ko.NextInQ(Dir,Floor)
			}			
			
			//println("NEXT: ", NextFloor)
			//println("FLOOR: ", Floor)
			if (NextFloor != -1){
				State = "MOVING"
				moveToFloor()
			}else if (NextFloor == -1){
				//println("Going to IDLE")
				State = "IDLE"
				Event = ""
			}
		}
	}
}

func ReadEvent() {
	var event driver.FloorEvent

	for ; true ; {
		event = <- driver.ElevChan
		Event = event.Event
		EventFloor = event.Floor
	}
}


func moveToFloor() {
	if (NextFloor == -1) {
		//println("Move to Floor: No order")
		return
	}

	if (Floor < NextFloor) {
		//println("MOVETOFLOOR:UP")
		driver.ELEV_set_motor_direction(1)
		Dir = "U"
	}else if (Floor > NextFloor){
		//println("Moving Down")
		driver.ELEV_set_motor_direction(-1)
		Dir = "D"
	}
}

func startUp() {
	println("STARTUP")
	Dir = "U"
	driver.ELEV_set_motor_direction(-1)
	for ; driver.ELEV_get_floor_sensor_signal() == -1 ; {
	}
	driver.ELEV_set_motor_direction(0)
	Floor = driver.ELEV_get_floor_sensor_signal()
}


