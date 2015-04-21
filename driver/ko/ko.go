package ko

import (.".././heis/"
		".././network/"
		"strconv")

type Queue struct {
	UP [4]int
	DOWN [4]int
	CMD [4]int
}

var Q Queue


func Q_init() {
	for i := 0 ; i < 4 ; i++ {
		Q.UP[i] = 0
		Q.DOWN[i] = 0
		Q.CMD[i] = 0
	}
}


func ButtonHandle(){
	var buttonEvent ButtonEvent

	for ; true ; {
		buttonEvent = <- ButtonChan
		AddOrder(buttonEvent.Floor, buttonEvent.Button)
		println("NY BESTILLING: ", buttonEvent.Floor, buttonEvent.Button)
	}	

}


func AddOrder(floor int, dir string) {
	//dir is U, D or C
	var event FloorEvent
    q := EmptyQ()

	switch dir {
	case "U" :
		println("NEW ORDER UP")
		if (Q.UP[floor] == 0){
			Q.UP[floor] = 1
			ELEV_set_button_lamp(0,floor,1)
		}

	case "D" :
		println("NEW ORDER DOWN")
		if (Q.DOWN[floor] == 0){
			Q.DOWN[floor] = 1
			ELEV_set_button_lamp(1,floor,1)
		}

	case "C" :
		println("NEW ORDER CMD")
		if (Q.CMD[floor] == 0){
			Q.CMD[floor] = 1
			ELEV_set_button_lamp(2,floor,1)
		}
	}

    if (q == 1) {
        event.Event = "NEW_ORDER"
		event.Floor = floor
		ElevChan <- event
    }
}


func EmptyQ()(int){
	for i := 0 ; i < 4; i++{
        if (Q.UP[i] == 1) || (Q.DOWN[i] == 1) || (Q.CMD[i] == 1){
            return 0
        }
	}
    return 1
}


func NextInQ(dir string, floor int) (int) {
	switch dir {
	case "U" :
		for i := floor ; i < 4 ; i++ {
			if (Q.UP[i] == 1) || (Q.CMD[i] == 1) {
				return i
			}
		}
		
		for j := 3 ; j > floor ; j-- {
			if (Q.DOWN[j] == 1) {
				return j
			}
		}

	case "D" :
		for i := floor ; i > -1 ; i-- {
			if (Q.DOWN[i] == 1) || (Q.CMD[i] == 1) {
				return i
			}
		}
	
		for j := 0 ; j < floor ; j++ {
			if (Q.UP[j] == 1){
				return j
			}
		}
	}

	for i := 0 ; i < 4 ; i++{
		if (Q.UP[i] == 1) || (Q.DOWN[i] == 1) || (Q.CMD[i] == 1){
			return i
		}
	}

	return -1
}


func RemoveOrder(floor int) {
	Q.UP[floor] = 0
	Q.DOWN[floor] = 0
	Q.CMD[floor] = 0

	if (network.Master){
		var temp network.Message
		temp.From = network.IP
		temp.Message = "rm" + strconv.Itoa(floor)
		network.NewMessage <- &temp
	}else{
		temp := "rm" + strconv.Itoa(floor)
		network.SendMessage(temp, network.Broadcast.Conn)
	}

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
