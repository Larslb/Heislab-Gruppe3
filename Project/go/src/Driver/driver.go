package Driver

import (
	"time"
)

// SKAL ELEV_INIT KJØRE TIL NÆRMESTE ETASJE?

func ReadElevPanel(buttonChan chan Queue.MyOrder){
	for {
		for i:=0;i<N_FLOORS;i++{
			elev_get_button_signal(BUTTON_COMMAND,i, buttonChan)
		}
	}
}

func ReadFloorPanel(buttonChan chan Queue.MyOrder){
	for{
		for i:=0;i<N_BUTTONS-1;i++{
			for j:=0;j<N_FLOORS;j++{
				elev_get_button_signal(i, j, buttonChan)
			}
		}
	}
}

func readSensors(sensorChan chan int){
	tmpSensorChan := make(chan int, 1) // Nødvendig med buffer?
	for {
		elev_get_floor_sensor_signal(tmpSensorChan)
		tmpVal := <- tmpSensorChan
		if tmpVal != -1 {
			sensorChan <- tmpVal		
		}
	}
}

func Fsm(dirOrNextFloor chan int, deleteOrderOnFloor chan int, currentFloor chan int) { // skal være i drivermodulen

	current_floor := -1 // initielt (må settes av ELEV_INIT)
	direction     := 0  // initielt
	next_floor    := -1 // initielt ingen bestillinger
	
	sensorChan := make(chan int)
	go readSensors(sensorChan)
	
	// ELEV_INIT HER? -> må ta inn sensorChan og sette current_floor
	
	buttonFloorChan := make(chan Queue.MyOrder) // Skal ikke brukes her
	buttonElevChan  := make(chan Queue.MyOrder) // Skal ikke brukes her
	go readElevPanel(buttonElevChan)		  // Skal ikke brukes her
	go readFloorPanel(buttonFloorChan)		  // Skal ikke brukes her
	
	
	STATE := WAIT
	
	for {
		switch STATE {
		
		
			case WAIT:
				
				currentFloor <- current_floor
				next_floor = <-dirOrNextFloor
				
				if next_floor < current_floor {
				
					direction = -1
					dirOrNextFloor <- direction
					elev_set_motor_direction(direction)
					STATE = MOVING
					
				} else if next_floor > current_floor {
				
					direction = 1
					dirOrNextFloor <- direction
					elev_set_motor_direction(direction)
					STATE = MOVING
					
				} else if next_floor == current_floor {
					
					dirOrNextFloor <- direction
					STATE = OPEN_DOOR
					
				} else if next_floor == -1 {
					
					direction = 0
					dirOrNextFloor <- direction
					
					//time.Sleep() for å ikke overbelaste QM med requests
				
				}
				
			case MOVING:
			
				current_floor = <- sensorChan
	
				currentFloor <- current_floor
				next_floor <- dirOrnextFloor
				dirOrNextFloor <- direction
				
				if current_floor == next_floor {
					STATE = DOOR_OPEN
				}
				
			case OPEN_DOOR:
				elev_set_motor_direction(0)
				elev_set_door_open_lamp(true)
				t := time.Now()
				for(!t.After(3*time.Seconds){
					currentFloor <- current_floor
					next_floor <- dirOrnextFloor
					dirOrNextFloor <- direction

					if current_floor == next_floor{
						t = time.Now()
					}
				}
				elev_set_door_open_lamp(false)				

				deleteOrderOnFloor <- current_floor
				ordersDeleted := <-deleteOrderOnFloor
				
				// SLETTE LYS
				if direction == 1{
					elev_set_button_lamp(BUTTON_CALL_UP, current_floor,0)
					elev_set_button_lamp(BUTTON_COMMAND, current_floor,0)
				} else {
					elev_set_button_lamp(BUTTON_CALL_DOWN, current_floor, 0)
					elev_set_button_lamp(BUTTON_COMMAND, current_floor,0)
				}
				
				STATE = WAIT				
		}
	}
}
