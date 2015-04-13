package ko

import ()

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


func addOrder(floor int, dir string) {
	//dir is UP, DOWN or CMD

    q = emptyQ()

	switch dir {
	case "UP" :
		if (Q_up[floor] == 0){
			Q_up[floor] = 1
			ELEV_set_button_lamp(0,floor,1)
		}

	case "DOWN" :
		if (Q_down[floor] == 0){
			Q_down[floor] = 1
			ELEV_set_button_lamp(1,floor,1)
		}

	case "CMD" :
		if (Q_cmd[floor] == 0){
			Q_cmd[floor] = 1
			ELEV_set_button_lamp(2,floor,1)
		}
	}

    if (q == 1) {
        //Send NEW_ORDER event
    }

}


func emptyQ()(int){
	for i := 0 ; i < 4; i++{
        if (Q_up[i] == 1) || (Q_down[i] == 1) || (Q_cmd[i] == 1){
            return 0
        }
    return 1
}
