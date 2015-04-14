package ko

import (.".././heis/")

var Q_up [4]int
var Q_down [4]int
var Q_cmd [4]int


func Q_init() {
	for i := 0 ; i < 4 ; i++ {
		Q_up[i] = 0
		Q_down[i] = 0
		Q_cmd[i] = 0
	}
}


func ButtonHandle(){
	var buttonEvent ButtonEvent

	for ; true ; {
		buttonEvent = <- ButtonChan
		addOrder(buttonEvent.Floor, buttonEvent.Button)
		println("NY BESTILLING: ", buttonEvent.Floor, buttonEvent.Button)
	}	

}


func addOrder(floor int, dir string) {
	//dir is UP, DOWN or CMD
	var event FloorEvent
    q := emptyQ()

	switch dir {
	case "UP" :
		println("NEW ORDER UP")
		if (Q_up[floor] == 0){
			Q_up[floor] = 1
			ELEV_set_button_lamp(0,floor,1)
		}

	case "DOWN" :
		println("NEW ORDER DOWN")
		if (Q_down[floor] == 0){
			Q_down[floor] = 1
			ELEV_set_button_lamp(1,floor,1)
		}

	case "CMD" :
		println("NEW ORDER CMD")
		if (Q_cmd[floor] == 0){
			Q_cmd[floor] = 1
			ELEV_set_button_lamp(2,floor,1)
		}
	}

    if (q == 1) {
        event.Event = "NEW_ORDER"
		event.Floor = floor
		ElevChan <- event
    }
}


func emptyQ()(int){
	for i := 0 ; i < 4; i++{
        if (Q_up[i] == 1) || (Q_down[i] == 1) || (Q_cmd[i] == 1){
            return 0
        }
	}
    return 1
}


func NextInQ(dir string, floor int) (int) {
	switch dir {
	case "UP" :
		for i := floor ; i < 4 ; i++ {
			if (Q_up[i] == 1) || (Q_cmd[i] == 1) {
				return i
			}
		}
		
		for j := 3 ; j > floor ; j-- {
			if (Q_down[j] == 1) {
				return j
			}
		}


	case "DOWN" :
		for i := floor ; i > -1 ; i-- {
			if (Q_down[i] == 1) || (Q_cmd[i] == 1) {
				return i
			}
		}
	
		for j := 0 ; j < floor ; j++ {
			if (Q_up[j] == 1){
				return j
			}
		}
	}

	for i := 0 ; i < 4 ; i++{
		if (Q_up[i] == 1) || (Q_down[i] == 1) || (Q_cmd[i] == 1){
			return i
		}
	}

	return -1
}


func RemoveOrder(floor int) {
	Q_up[floor] = 0
	Q_down[floor] = 0
	Q_cmd[floor] = 0

	switch floor {
		case 0:
			ELEV_set_button_lamp(0,floor,0)
			ELEV_set_button_lamp(2,floor,0)

		case 3:
			ELEV_set_button_lamp(1,floor,0)
			ELEV_set_button_lamp(2,floor,0)

		case 1, 2:
			ELEV_set_button_lamp(0,floor,0)
			ELEV_set_button_lamp(1,floor,0)
			ELEV_set_button_lamp(2,floor,0)
	}
}
