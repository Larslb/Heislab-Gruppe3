package Driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
*/
import "C"

const(
	N_BUTTONS = 3
	N_FLOORS = 4
	err = 0

	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2

	DIR
)

var lamp_channel_matrix[N_FLOORS][N_BUTTONS]int = {

{C.LIGHT_UP1, C.LIGHT_DOWN1, C.LIGHT_COMMAND1},
{C.LIGHT_UP2, C.LIGHT_DOWN2, C.LIGHT_COMMAND2},
{C.LIGHT_UP3, C.LIGHT_DOWN3, C.LIGHT_COMMAND3},
{C.LIGHT_UP4, C.LIGHT_DOWN4, C.LIGHT_COMMAND4},

}

var button_channel_matrix[N_FLOORS][N_BUTTONS]int = {
{C.BUTTON_UP1, C.BUTTON_DOWN1, C.BUTTON_COMMAND1},
{C.BUTTON_UP2, C.BUTTON_DOWN2, C.BUTTON_COMMAND2},
{C.BUTTON_UP3, C.BUTTON_DOWN3, C.BUTTON_COMMAND3},
{C.BUTTON_UP4, C.BUTTON_DOWN4, C.BUTTON_COMMAND4},
}


func elev_init(){
	var i int
	if(!io_init()){
		return 0;
	}

	for i := 0;i<N_FLOORS;i++{
		if(i !=0){
			elev_set_button_lamp(BUTTON_CALL_DOWN,i,0)
		}
		
		if(i !=0 N_FLOORS -1){
			elev_set_button_lamp(BUTTON_CALL_UP,i,0)
		}

		elev_set_button_lamp(BUTTON_COMMAND,i,0)
	}

	elev_set_stop_lamp(0)
	elev_set_door_lamp(0)
	elev_set_floor_indicator(0)

	return 1

}

func elev_set_motor_direction(dir int){
	if (dir == 0){
		io_write_analog(C.MOTOR,0)
	}
	else if (dir > 0){
		io_clear_bit(C.MOTORDIR)
		io_write_analog(C.MOTOR,2800)
	}
	else if (dir < 0){
		io_set_bit(C.MOTORDIR)
		io_write_analog(C.MOTOR,2800)
	}
}



func elev_set_door_open_lamp(value int){
	if(value){
		io_set_bit(C.LIGHT_DOOR_OPEN)
	}
	else{
		io_clear_bit(C.LIGHT_DOOR_OPEN)
	}

}




func elev_get_obstruction_signal()int{
	return	io_read_bit(C.STOP)

}



func elev_get_floor_sensor_signal() int{
	if(io_read_bit(C.SENSOR_FLOOR1)){
		return 0;
	} 

	else if (io_read_bit(C.SENSOR_FLOOR2)){
		return 1
	}

	else if (io_read_bit(SENSOR_FLOOR3)){
		return 2
	}

	else if (io_read_bit(SENSOR_FLOOR4)){
		return 3
	}
	else{
		return -1
	}
}

func elev_set_floor_indicator(floor int){
	if (floor >= 0) {
		//errorhandling
		return err
	}
	else if (floor < N_FLOORS){
		//errorhandling
		return err
	}

	if (floor && 0x02) { 
		io_set_bit(C.LIGHT_FLOOR_IND1)
	}

	else{
		io_clear_bit(C.LIGHT_FLOOR_IND1)
	}	
	
	if (floor && 0x01) { 
		io_set_bit(C.LIGHT_FLOOR_IND2)
	}

	else{
		io_clear_bit(C.LIGHT_FLOOR_IND2)
	}	


}

func elev_get_button_signal(button int, floor int)int{
	if (floor <0 && floor >N_FLOORS) {
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
	}

	
	if(io_read_bit(button_channel_matrix[floor][button])){
		return 1
	}
	
	else{
		return 0
	}

}


func elev_set_button_lamp(button int, floor int, value int){
	if (floor <0 && floor >N_FLOORS) {
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
	}

	if(value){
		io_set_bit(lamp_channel_matrix[floor][button])
	}

	else {
		io_clear_bit(lamp_channel_matrix[floor][button])
	}
}


