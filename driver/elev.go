package driver

var LAMP_CHANNEL = [4][3]int{{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
			{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    			{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    			{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

var BUTTON_CHANNEL = [4][3]int{{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    			{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    			{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    			{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}}


const N_FLOORS = 4
const N_BUTTONS = 3


func elev_init() (int){
	var i int
	
	//Init hardware
	if IO_init() == 0 {
		return 0;
	}
	
	//Zero all floor button lamps
	for i = 0 ; i < N_FLOORS ; i++{
		if i != 0{
			elev_set_button_lamp(1, i, 0)
		}
		
		if i != N_FLOORS - 1 {
			elev_set_button_lamp(0, i, 0)
		}

		elev_set_button_lamp(2, i, 0)
	}

	//Clear stop lamp, door open lamp, and set floor indicator to ground floor
	elev_set_stop_lamp(0)
    elev_set_door_open_lamp(0)
    elev_set_floor_indicator(0)

	return 1;
}
