package Driver

import (
	"fmt"
	"time"
	"ElevLib"
)

func ReadElevPanel(buttonChan chan ElevLib.MyOrder){
	for {
		for i:=0;i<ElevLib.N_FLOORS;i++{
			elev_get_button_signal(ElevLib.BUTTON_COMMAND,i, buttonChan)
		}
	}
}

func ReadFloorPanel(buttonChan chan ElevLib.MyOrder){
	for{
		for i:=0;i<ElevLib.N_BUTTONS-1;i++{
			for j:=0;j<ElevLib.N_FLOORS;j++{
				elev_get_button_signal(i, j, buttonChan)
			}
		}
	}
}

func readSensors(sensorChan chan int){
	
	for {
		tmpVal := elev_get_floor_sensor_signal()
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

func Fsm(nextFloorChan chan int, deleteOrderOnFloorChan chan int, currentFloorChan chan int, directionChan chan int, setLightsChan chan []int) {

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
	
	STATE := ElevLib.WAIT
	
	for {
		switch STATE {
		
		
			case ElevLib.WAIT:
				fmt.Println("FSM: ", "STATE = WAIT")
				currentFloorChan <- current_floor
				next_floor = <-nextFloorChan
				
				if next_floor == -1 {
					fmt.Println("FSM: ","current_floor = ", current_floor, "next_floor = ", next_floor)
					direction = 0
					directionChan <- direction
					
					time.Sleep(30*time.Millisecond)
				} else if next_floor < current_floor {
					fmt.Println("FSM: ","current_floor = ", current_floor, "next_floor = ", next_floor)
					direction = -1
					fmt.Println(direction)
					directionChan <- direction
					elev_set_motor_direction(direction)
					STATE = ElevLib.MOVING
					fmt.Println("FSM: ","MOVING DOWN")
					time.Sleep(30*time.Millisecond)
					
				} else if next_floor > current_floor {
					fmt.Println("FSM: ","current_floor = ", current_floor, "next_floor = ", next_floor)
					direction = 1
					directionChan <- direction
					elev_set_motor_direction(direction)
					STATE = ElevLib.MOVING
					fmt.Println("FSM: ","MOVING UP")
					time.Sleep(30*time.Millisecond)
					
				} else if next_floor == current_floor {
					fmt.Println("FSM: ","current_floor = ", current_floor, "next_floor = ", next_floor)
					directionChan <- direction
					STATE = ElevLib.OPEN_DOOR
					time.Sleep(30*time.Millisecond)
					
				}
				
			case ElevLib.MOVING:
				fmt.Println("FSM: ","STATE = MOVING")
				current_floor = <- sensorChan
	
				currentFloorChan <- current_floor
				next_floor = <- nextFloorChan
				directionChan <- direction
				
				if current_floor == next_floor {
					STATE = ElevLib.OPEN_DOOR
				}
				
			case ElevLib.OPEN_DOOR:
				fmt.Println("FSM: ","STATE = OPEN_DOOR")
				elev_set_motor_direction(0)
				elev_set_door_open_lamp(true)
				t := time.Now()
				t2 := t.Add(3*time.Second)
				for !t.After(t2) {
					currentFloorChan <- current_floor
					next_floor = <- nextFloorChan
					directionChan <- direction

					if current_floor == next_floor{
						t = time.Now()
						t2 =t.Add(3*time.Second)
					}
				}
				elev_set_door_open_lamp(false)				

				deleteOrderOnFloorChan <- current_floor
				STATE = ElevLib.WAIT				
		}
	}
}

