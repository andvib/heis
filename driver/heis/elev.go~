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

//var elev_motor_direction int;

type elev_motor_direction int
const (
	DIRN_DOWN = -1
	DIRN_STOP = 0
	DIRN_UP = 1
)


//var elev_button_type int;
type elev_button_type int
const (
	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2
)

func ELEV_init() (int){
	var i int
	
	//Init hardware
	if IO_init() == 0 {
		return 0
	}
	
	//Zero all floor button lamps
	for i = 0 ; i < N_FLOORS ; i++{
		if i != 0{
			ELEV_set_button_lamp(1, i, 0)
		}
		
		if i != N_FLOORS - 1 {
			ELEV_set_button_lamp(0, i, 0)
		}

		ELEV_set_button_lamp(2, i, 0)
	}

	//Clear stop lamp, door open lamp, and set floor indicator to ground floor
	ELEV_set_stop_lamp(0)
    ELEV_set_door_open_lamp(0)
   	ELEV_set_floor_indicator(0)

	return 1
}


func ELEV_set_motor_direction(direction elev_motor_direction) {
	if direction == 0{
		IO_write_analog(MOTOR, 0)	
	}else if direction > 0 {
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR,2800)
	}else if direction < 0 {
		IO_set_bit(MOTORDIR)
		IO_write_analog(MOTOR, 2800)
	}
	
}


func ELEV_set_door_open_lamp(value int) {
	if value != 0 {
		IO_set_bit(LIGHT_DOOR_OPEN)
	}else{
		IO_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func ELEV_get_obstruction_signal() (int){
	return IO_read_bit(OBSTRUCTION)
}

func ELEV_get_stop_signal() (int){
	return IO_read_bit(STOP)
}


func ELEV_set_stop_lamp(value int) {
	if value != 0 {	
		IO_set_bit(LIGHT_STOP)
	}else {
		IO_clear_bit(LIGHT_STOP)
	}
}

func ELEV_get_floor_sensor_signal() (int){
	if IO_read_bit(SENSOR_FLOOR1) == 1{
		return 0
	}else if IO_read_bit(SENSOR_FLOOR2) == 1{
		return 1
	}else if IO_read_bit(SENSOR_FLOOR3) == 1{
		return 2
	}else if IO_read_bit(SENSOR_FLOOR4) == 1{
		return 3
	}else {
		return -1
	}
}


func ELEV_set_floor_indicator(floor int) {
	if floor >= 0 || floor < N_FLOORS {
		return 
	}
	
	/*tall := 1
	test := floor & tall
	println(test)
	println(floor)*/

	
   	 // Binary encoding. One light must always be on.
	if (floor & 0x02) != 0 {
		IO_set_bit(LIGHT_FLOOR_IND1)
	}else {
		IO_clear_bit(LIGHT_FLOOR_IND1)
	}
	if (floor & 0x01) != 0 {
		IO_set_bit(LIGHT_FLOOR_IND2)
	}else {
		IO_clear_bit(LIGHT_FLOOR_IND2)
	}

}


func ELEV_get_button_signal(button elev_button_type, floor int) (int){

	if !(floor >= 0 || floor < N_FLOORS) {
		return 0
	}
	if button == BUTTON_CALL_UP && floor == N_FLOORS -1 {
		return 0
	}
	if button == BUTTON_CALL_DOWN && floor == 0 {
		return 0
	}
	if !(button == BUTTON_CALL_UP || button == BUTTON_CALL_DOWN || button == BUTTON_COMMAND) {
		return 0
	}
	
	if IO_read_bit(BUTTON_CHANNEL[floor][button]) == 1{
		return 1
	} else {
		return 0
	}
	
}

func ELEV_set_button_lamp(button elev_button_type, floor int, value int) {
	
	if !(floor >= 0 || floor < N_FLOORS) {
		return
	}
	if button == BUTTON_CALL_UP && floor == N_FLOORS -1 {
		return
	}
	if button == BUTTON_CALL_DOWN && floor == 0 {
		return
	}
	if !(button == BUTTON_CALL_UP || button == BUTTON_CALL_DOWN || button == BUTTON_COMMAND) {
		return
	}

	println(button)
	println(floor)
	println(value)

	if value == 1{
		IO_set_bit(LAMP_CHANNEL[floor][button])
	}else {
		IO_clear_bit(LAMP_CHANNEL[floor][button])
	}
}

/*
func poll() {
	
	
	for i:=0; i < N_FLOORS; i++ {
		for j:=0; j < 2; j++{
			BUTTON_CHANNEL[j][i] := elev_get_button_signal(j ,i)
		}
	}
}
*/















