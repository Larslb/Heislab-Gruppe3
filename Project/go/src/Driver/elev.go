package Driver


const (
	BUTTON_CALL_UP int = 0
	BUTTON_CALL_DOWN int = 1
	BUTTON_COMMAND int = 2
)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int {
{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int {
{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}




func Elev_init()int {
	if(io_init() != 1){
		return 0;
	}

	for i := 0;i<N_FLOORS;i++{
		if(i !=0){
			elev_set_button_lamp(BUTTON_CALL_DOWN,i,0)
		}
		
		if(i != N_FLOORS -1){
			elev_set_button_lamp(BUTTON_CALL_UP,i,0)
		}

		elev_set_button_lamp(BUTTON_COMMAND,i,0)
	}

	elev_set_stop_lamp(false)
	elev_set_door_open_lamp(false)
	elev_set_floor_indicator(0)

	return 1

}

func elev_set_motor_direction(dir int){
	if (dir == 0){
		io_write_analog(MOTOR,0)
	} else if (dir > 0){
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR,2800)
	} else if (dir < 0){
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR,2800)
	}
}



func elev_set_door_open_lamp(value bool){
	if(value){
		io_set_bit(LIGHT_DOOR_OPEN)
	} else{
		io_clear_bit(LIGHT_DOOR_OPEN)
	}

}




func elev_get_obstruction_signal(obstrChan chan int){
	obstrChan <- io_read_bit(OBSTRUCTION)

}

func elev_get_stop_signal(stopChan chan int) {
	stopChan <- io_read_bit(STOP)
}

func elev_set_stop_lamp(value int){
	if value{
		io_set_bit(LIGHT_STOP)
	} else{
		io_clear_bit(LIGHT_STOP)
	}
}



func elev_get_floor_sensor_signal(sensorchan chan int) {
	if(io_read_bit(SENSOR_FLOOR1)==1){
		sensorchan <- 0;
	} else if (io_read_bit(SENSOR_FLOOR2)==1){
		sensorchan <- 1;
	} else if (io_read_bit(SENSOR_FLOOR3)==1){
		sensorchan <- 2;
	} else if (io_read_bit(SENSOR_FLOOR4)==1){
		sensorchan <- 3;
	} else{
		sensorchan <- -1;
	}
}

func elev_set_floor_indicator(floor int){
	/*if (floor >= 0) {
		//errorhandling
		return err
	}
	else if (floor < N_FLOORS){
		//errorhandling
		return err
	}*/

	
	if (floor & 0x02) != 0 { 
		io_set_bit(LIGHT_FLOOR_IND1)
	} else{
		io_clear_bit(LIGHT_FLOOR_IND1)
	}	

	if (floor & 0x01) != 0 { 
		io_set_bit(LIGHT_FLOOR_IND2)
	} else{
		io_clear_bit(LIGHT_FLOOR_IND2)
	}	


}

func elev_get_button_signal(button, floor int, buttonChan chan myOrder){
	/*if (floor <0 && floor >N_FLOORS) {
		//errorhandling
		return err
	}	
	
	else if(!(button == C.BUTTON_CALL_UP && floor == N_FLOORS -1)){
		//errorhandling
		return err
	}

	else if(!(button == C.BUTTON_CALL_DOWN && floor == 0)){
		//errorhandling
		return err
	}

	else if (button == C.BUTTON_CALL_UP || button == C.BUTTON_CALL_DOWN || button == C.BUTTON_COMMAND){
		//errorhandling
		return err
	}*/

	var buttonOrder myOrder
	buttonOrder.buttonType = button
	buttonOrder.floor = floor
	
	if(io_read_bit(button_channel_matrix[floor][button]) == 1){
		buttonOrder.value = 1
		
	} else{
		buttonOrder.value = 0
		
	}
	buttonChan <- buttonOrder
}


func elev_set_button_lamp(button int, floor int, value int){
	/*if (floor <0 && floor >N_FLOORS) {
		//errorhandling
		return err
	}	
	
	else if(!(button == C.BUTTON_CALL_UP && floor == N_FLOORS -1)){
		//errorhandling
		return err
	}

	else if(!(button == C.BUTTON_CALL_DOWN && floor == 0)){
		//errorhandling
		return err
	}

	else if (button == C.BUTTON_CALL_UP || button == C.BUTTON_CALL_DOWN || button == C.BUTTON_COMMAND){
		//errorhandling
		return err
	}*/

	if(value == 1){
		io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		io_clear_bit(lamp_channel_matrix[floor][button])
	}
}
