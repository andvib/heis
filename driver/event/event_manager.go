package event

import(".././heis/"
		"time"
		".././ko/")

var State string //IDLE, MOVING, DOOR_OPEN
var Event string //NEW_ORDER, ORDERED_FLOOR_REACHED, EMPTY_QUEUE
var EventFloor int
var Floor int
var CurrentFloor int
var NextFloor int
var Dir string

func StateMachine(){
	startUp()
	for ; true ; {
		NextFloor = ko.NextInQ(Dir,CurrentFloor) 
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
				Event = ""
				}
			}

		case "MOVING" :
			if Event == "NEW_FLOOR" {
				println("NEW_FLOOR")
				println("NextFloor: ", NextFloor)
				CurrentFloor := driver.ELEV_get_floor_sensor_signal()
				if (CurrentFloor != -1){
					Floor = CurrentFloor
				}
				if (driver.ELEV_get_floor_sensor_signal() == NextFloor){
					driver.ELEV_set_motor_direction(0)
					State = "DOOR_OPEN"
				}
				//Event = ""
			}/*else if Event == "NEW_ORDER" {
				NextFloor = ko.NextInQ(Dir,CurrentFloor)
				moveToFloor()
				Event = ""
			}*/

		case "DOOR_OPEN" :
			if (driver.ELEV_get_floor_sensor_signal() != -1){
				driver.ELEV_set_door_open_lamp(1)
				time.Sleep(1000*time.Millisecond)
				driver.ELEV_set_door_open_lamp(0)
				ko.RemoveOrder(CurrentFloor)
				NextFloor = ko.NextInQ(Dir,Floor)
			}			
			if (NextFloor != -1){
				State = "MOVING"
				moveToFloor()
			}else if (NextFloor == -1){
				State = "IDLE"
				Event = ""
			}
		}
	}
}

func ReadEvent() {
	var event driver.FloorEvent

	for ; true ; {
		//if Event == "" {
			event = <- driver.ElevChan
			Event = event.Event
			EventFloor = event.Floor
		//}
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


