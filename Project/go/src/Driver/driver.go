package Driver

import (
	"time"
	"ElevLib"
)

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
	tmpSensorChan := make(chan int)
	// defer close tmpSensorChan ??
	for {
		elev_get_floor_sensor_signal(tmpSensorChan)
		tmpVal := <- tmpSensorChan
		if tmpVal != -1 {
			sensorChan <- tmpVal		
		}
	}
}

func setLights(setLightsChan chan []int) {

	for{
		lightCommand := <- setLightsChan

		if lightCommand[0] == ElevLib.BUTTON_COMMAND {
			elev_set_button_lamp(lightCommand[0], lightCommand[1], lightCommand[2])
		}else if lightCommand[0] == ElevLib.BUTTON_CALL_UP {
			elev_set_button_lamp(lightCommand[0], lightCommand[1], lightCommand[2])
		}else if lightCommand[0] == ElevLib.BUTTON_CALL_DOWN {
			elev_set_button_lamp(lightCommand[0], lightCommand[1], lightCommand[2])
		}
	}
}

func Fsm(nextFloorChan chan int, deleteOrderOnFloorChan chan int, currentFloorAndDirChan chan int, setLightsChan chan []int) {

	current_floor := -1
	direction     := 0
	next_floor    := -1
	errorVar      := false
	
	sensorChan := make(chan int)
	go readSensors(sensorChan)
	go setLights(setLightsChan)

	// defer close sensorChan	

	current_floor, errorVar = elev_init(sensorChan)
	
	if !errorVar {
		fmt.Println("ERROR: elev.init() did not succeed ")
		// ERRORHANDLING - INIT DID NOT SUCCEED
	}


	// VI HAR INGEN SET LIGHTS ON/OFF HÅNDTERING PÅ TVERS AV HEISENE
	
	STATE := WAIT
	
	for {
		switch STATE {
		
		
			case WAIT:
				
				currentFloorAndDirChan <- current_floor
				next_floor = <-nextFloorChan
				
				if next_floor < current_floor {
				
					direction = -1
					currentFloorAndDirChan <- direction
					elev_set_motor_direction(direction)
					STATE = MOVING
					
				} else if next_floor > current_floor {
				
					direction = 1
					currentFloorAndDirChan <- direction
					elev_set_motor_direction(direction)
					STATE = MOVING
					
				} else if next_floor == current_floor {
					
					currentFloorAndDirChan <- direction
					STATE = OPEN_DOOR
					
				} else if next_floor == -1 {
					
					direction = 0
					currentFloorAndDirChan <- direction
					
					time.Sleep(300*time.Millisecond)
				}
				
			case MOVING:
			
				current_floor = <- sensorChan
	
				currentFloorAndDirChan <- current_floor
				next_floor <- nextFloorChan
				currentFloorAndDirChan <- direction
				
				if current_floor == next_floor {
					STATE = DOOR_OPEN
				}
				
			case OPEN_DOOR:
				elev_set_motor_direction(0)
				elev_set_door_open_lamp(true)
				t := time.Now()
				for(!t.After(3*time.Seconds){
					currentFloorAndDirChan <- current_floor
					next_floor <- nextFloorChan
					currentFloorAndDirChan <- direction

					if current_floor == next_floor{
						t = time.Now()
					}
				}
				elev_set_door_open_lamp(false)				

				deleteOrderOnFloorChan <- current_floor
				//ordersDeleted := <-deleteOrderOnFloorChan // Foreløpig sender vi ikke bekreftelse fra queue om at etasjen er slettet
				
				STATE = WAIT				
		}
	}
}

}
