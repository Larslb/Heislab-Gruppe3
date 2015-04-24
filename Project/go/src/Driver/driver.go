package Driver

import (
	"fmt"
	"time"
	"ElevLib"
)

func ReadElevPanel(buttonChan chan ElevLib.MyOrder, read chan bool){
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

func ReadSensors(sensorChan1 chan int,sensorChan2 chan int,sensorChan3 chan int){  // ENDRET TIL EXPORT FUNC
	
	current_floor := -1

	for {

		select{
			case sensorChan1 <- current_floor:
			case sensorChan2 <- current_floor:
			default:
				tmpVal := elev_get_floor_sensor_signal()
				if tmpVal != -1 {
				sensorChan <- tmpVal		
		}
		}
		
	}
}

func setLights(setLightsOn chan []int, setLightsOff chan []int) {

	for{
		select {
			case lightCommand := <- setLightsOn:
				

			case lightCommand := <- setLightsOff:


				
		}
	}
}

/*
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
}*/


func ready2receive(rdy2rcv chan bool, reqNewOrder chan int) {
	for {
		if <-rdy2rcv{
			reqNewOrder <- 1
		}
	}
}

func floor_reached(floorReached chan bool, floorSensor chan int, newFloor chan int, floor){
	go2Floor := floor
	
	for {
		case go2Floor = <- newFloor: // funker dette?
		case current_floor := <- floorSensor:
			elev_set_floor_indicator(current_floor)
			if current_floor == go2Floor {
				floorReached <- current_floor
				return
			}
	}
}

func FSM(sendReq2EM chan ElevLib.NewReqFSM, status chan int, orderHandledChan chan int, setLightsOn chan []int, init chan chan int, setlights chan bool) {
	
	rcvFromQueue := make(chan [2]int)
	updFromQueue := make(chan int)
	
	// Used in func ready2receive
	rdy2rcv	 := make(chan bool)
	reqNewOrder  := make(chan int)
	
	go ready2receive(rdy2rcv, reqNewOrder)
	
	// Used in goroutine func floorReached()
	newFloor     := make(chan int)
	floorReached := make(chan int)
	
	floorSensor   := make(chan int)
	go Driver.ReadSensor(floorSensor)
	go setLights(setLightsOn, setLightsOff)
	
	for {
		select {
			case order := <-rcvFromQueue:

				if order[1] != -1 {
					
					elev_set_motor_direction(order[0])

					go floor_reached(floorReached, floorSensor, newFloor, order[1])
				
					for {
						select{
							case floor := <-floorReached:
								updFromQueue <- 1 // terminate Update routine in queue
								break;

							case newOrder := <- updFromQueue:
								newFloor<- newOrder.Floor
						}
					}
				
					elev_set_motor_direction(0)

					elev_set_door_open_lamp(value bool)  // MÅ FIKSES PÅ! HOLDES ÅPEN I 3 SEK



					orderHandledChan <- floor

					// send false til setLights go-routine

					<-setlights
					if dir == 1 {
						setLightsOff <- []int{ElevLib.BUTTON_CALL_UP, floor, false}  // MÅ ENDRE PÅ ELEV_SET_LIGHTS fra int til bool
						setLightsOff <- []int{ElevLib.BUTTON_COMMAND, floor, false}
					} else if dir == -1 {
						setLightsOff <- []int{ElevLib.BUTTON_CALL_DOWN, floor, false}
						setLightsOff <- []int{ElevLib.BUTTON_COMMAND, floor, false}
					} else {
						setLightsOff <- []int{ElevLib.BUTTON_COMMAND, floor, false}
					}

					rdy2rcv <- true
				}else{
					time.Sleep(1*time.Second)
				}
			case <-reqNewOrder:
				sendReq2EM <- ElevLib.NewReqFSM{
					OrderChan: rcvFromQueue
					UpdateOrderChan: updFromQueue
					Current_floor: 0  //no-care?
					Direction: 0      //no-care?
				}
				
			case <-status: // Brukes når heisen har stått stille en stund
				sendReq2EM <- ElevLib.NewReqFSM{
					OrderChan: rcvFromQueue
					UpdateOrderChan: updFromQueue
					Current_floor: 0  //no-care?
					Direction: 0      //no-care?
				}

			case c := <-init:
				current_floor, err := Driver.Elev_init(floorSensor)
				if err {
					fmt.Println("Unable to initialize elevator!")
					//do something: return
				c <- current_floor
	}
				
		}
	}
}

