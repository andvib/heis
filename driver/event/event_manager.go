package event

import(".././heis/"
		"time"
		".././ko/"
		".././network/")

var State string //IDLE, MOVING, DOOR_OPEN
var Event string //NEW_ORDER, ORDERED_FLOOR_REACHED, EMPTY_QUEUE
var EventFloor int
var Floor int
var NextFloor int
var Dir string
var TimerMoving time.Time
var TimerDoor time.Time

func StateMachine(){
	startUp()
	for ; true ; {
		NextFloor = ko.NextInQ(Dir,Floor) 
		switch State {
	
		case "IDLE" :
			if Event == "NEW_ORDER"{
				println("IDLE:NEW_ORDER")
				NextFloor = ko.NextInQ(Dir,Floor)
				println(NextFloor)
				if NextFloor == Floor {
					TimerDoor = time.Now()
					State = "DOOR_OPEN"
				}else{
					TimerMoving = time.Now()
					State = "MOVING"
					moveToFloor()
				Event = ""
				}
			}

		case "MOVING" :
			if Event == "NEW_FLOOR" {
				temp := driver.ELEV_get_floor_sensor_signal()	
				if (temp != -1){
					Floor = temp
				}
				if (temp == NextFloor){
					driver.ELEV_set_motor_direction(0)
					TimerDoor = time.Now()
					State = "DOOR_OPEN"
				}
			}

		case "DOOR_OPEN" :
			if (driver.ELEV_get_floor_sensor_signal() != -1){
				driver.ELEV_set_door_open_lamp(1)
				time.Sleep(3000*time.Millisecond)
				driver.ELEV_set_door_open_lamp(0)
				ko.RemoveOrder(Floor)
				NextFloor = ko.NextInQ(Dir,Floor)
				println(NextFloor)
			}			
			if (NextFloor != -1){
				TimerMoving = time.Now()
				State = "MOVING"
				moveToFloor()
			}else if (NextFloor == -1){
				State = "IDLE"
				Event = ""
			}
		}
	}
}


func ReadTimer(){
	for {
		if (time.Since(TimerMoving) > 10*time.Second) && (State == "MOVING"){
			println("Elevator not moving!")
			network.Alive = false
		}else if (time.Since(TimerDoor) > 10*time.Second) && (State == "DOOR_OPEN"){
			println("Door stuck!")
			network.Alive = false
		}else{
			network.Alive = true
		}
		time.Sleep(5*time.Second)
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
		return
	}

	if (Floor < NextFloor) {
		driver.ELEV_set_motor_direction(1)
		Dir = "U"
	}else if (Floor > NextFloor){
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


