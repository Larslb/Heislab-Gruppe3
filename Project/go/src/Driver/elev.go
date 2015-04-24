package Driver
import (
	"ElevLib"
	//"fmt"
)


var lamp_channel_matrix = [ElevLib.N_FLOORS][ElevLib.N_BUTTONS]int {
{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [ElevLib.N_FLOORS][ElevLib.N_BUTTONS]int {
{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}




func Elev_init(sensorChan chan int) (int, bool) {
	elev_set_motor_direction(0)
	
	if(io_init() != 1){
		return -1, true;
	}

	elev_set_door_open_lamp(false)

	for i := 0;i<ElevLib.N_FLOORS;i++{
		if(i !=0){
			elev_set_button_lamp(ElevLib.BUTTON_CALL_DOWN ,i,0)
		}
		
		if(i != ElevLib.N_FLOORS-1){
			elev_set_button_lamp(ElevLib.BUTTON_CALL_UP,i,0)
		}
		elev_set_button_lamp(ElevLib.BUTTON_COMMAND,i,0)
	}
	
	elev_set_motor_direction(-1)

	current_floor := <- sensorChan
	elev_set_motor_direction(0)
	elev_set_floor_indicator(current_floor)

	return current_floor, false
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

func elev_get_floor_sensor_signal() int {
	if(io_read_bit(SENSOR_FLOOR1)==1){
		return 0;
	} else if (io_read_bit(SENSOR_FLOOR2)==1){
		return 1;
	} else if (io_read_bit(SENSOR_FLOOR3)==1){
		return 2;
	} else if (io_read_bit(SENSOR_FLOOR4)==1){
		return 3;
	} else{
		return -1;
	}
}

func elev_set_floor_indicator(floor int){
	/*if (floor >= 0) {
		//errorhandling
		return err
	}
	else if (floor < ElevLib.N_FLOORS){
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

func Elev_get_button_signal(button, floor int) int {
	/*if (floor <0 && floor >ElevLib.N_FLOORS) {
		//errorhandling
		return err
	}	
	
	else if(!(button == C.ElevLib.BUTTON_CALL_UP && floor == ElevLib.N_FLOORS -1)){
		//errorhandling
		return err
	}

	else if(!(button == C.ElevLib.BUTTON_CALL_DOWN  && floor == 0)){
		//errorhandling
		return err
	}

	else if (button == C.ElevLib.BUTTON_CALL_UP || button == C.ElevLib.BUTTON_CALL_DOWN  || button == C.ElevLib.BUTTON_COMMAND){
		//errorhandling
		return err
	}
	fmt.Println(button_channel_matrix[floor][button])*/
	if (io_read_bit(button_channel_matrix[floor][button]) != 0){
		return 1
	}else{
		return 0
	}
}


func elev_set_button_lamp(button int, floor int, value int){
	/*if (floor <0 && floor >ElevLib.N_FLOORS) {
		//errorhandling
		return err
	}	
	
	else if(!(button == C.ElevLib.BUTTON_CALL_UP && floor == ElevLib.N_FLOORS -1)){
		//errorhandling
		return err
	}

	else if(!(button == C.ElevLib.BUTTON_CALL_DOWN  && floor == 0)){
		//errorhandling
		return err
	}

	else if (button == C.ElevLib.BUTTON_CALL_UP || button == C.ElevLib.BUTTON_CALL_DOWN  || button == C.ElevLib.BUTTON_COMMAND){
		//errorhandling
		return err
	}*/

	if(value == 1){
		io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		io_clear_bit(lamp_channel_matrix[floor][button])
	}
}
