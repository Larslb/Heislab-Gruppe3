package Driver

import(
	"fmt"
	"time"
	".././ElevLib"
)

func ReadElevPanel(buttonChan chan ElevLib.MyOrder){
	for {
		for i:=0;i<ElevLib.N_FLOORS;i++{
			if Elev_get_button_signal(ElevLib.BUTTON_COMMAND,i) != 0 {
				buttonChan <- ElevLib.MyOrder{
					Ip: "",
					ButtonType: ElevLib.BUTTON_COMMAND,
					Floor: i,
				}
			}
		}
		time.Sleep(80*time.Millisecond)
	}
}

func ReadFloorPanel(buttonChan chan ElevLib.MyOrder){
	for{
		for i:=0;i<ElevLib.N_BUTTONS-1;i++{
			for j:=0;j<ElevLib.N_FLOORS;j++{
				if Elev_get_button_signal(i,j) != 0{
					buttonChan <- ElevLib.MyOrder{
						Ip: "",
						ButtonType: i,
						Floor: j,
					}
				}
			}
		}
		time.Sleep(80*time.Millisecond)
	}
}

func ReadSensors(sensorChan chan int){  // ENDRET TIL EXPORT FUNC
	
	for { // INIT ONLY
		tmpVal1 := Elev_get_floor_sensor_signal()
		if tmpVal1 != -1 {
			sensorChan <-tmpVal1
			break
		}
	}

	current_floor := -1

	for {
		tmpVal2 := Elev_get_floor_sensor_signal()
			if tmpVal2 != -1 && tmpVal2 != current_floor {
				current_floor = tmpVal2
				sensorChan <- current_floor
			}

		time.Sleep(time.Millisecond)
	}
}

func SetLights(setLightsOn chan []int, setLightsOff chan []int) {

	for{
		select {
			case lightCommand := <- setLightsOn:
				elev_set_button_lamp(lightCommand[0], lightCommand[1], lightCommand[2])
			case lightCommand := <- setLightsOff:
				elev_set_button_lamp(lightCommand[0], lightCommand[1], lightCommand[2])		
		}
		time.Sleep(time.Millisecond)
	}
}

func floor_reached(floorReached chan ElevLib.NextOrder, updateCurrentFloor chan int,floorSensor chan int, newOrder chan ElevLib.NextOrder){
	
	go2Order :=  ElevLib.NextOrder{
		ButtonType: -2,
		Floor: -1,
		Direction: 0,
	} // TRENGER VI DENNE?


	current_floor := -1
	
	for {
		select {
			case order := <- newOrder:

				if order.Floor == current_floor {
					fmt.Println("floor_reached: go2Order = ", go2Order)
					floorReached <- order
				} else {
					go2Order = order
				}

			case current_floor = <- floorSensor:

				updateCurrentFloor <- current_floor
				elev_set_floor_indicator(current_floor)
				if current_floor == go2Order.Floor {
					floorReached <- go2Order
				}
		}
		time.Sleep(time.Millisecond)
	}
}




func Fsm( rcvChannelsFromQueue chan ElevLib.QM2FSMchannels, setLightsOff chan []int) {

	floorSensor := make(chan int)
	go ReadSensors(floorSensor)

	_,err := Elev_init(floorSensor)
	if err {
		fmt.Println("FSM: Could not initiate elevator")  // Returnere error til MAIN?
	}


	channels  := <- rcvChannelsFromQueue

	breakbool := false
	var reachedFloor ElevLib.NextOrder

	
	floorReached := make(chan ElevLib.NextOrder)
	newNextOrder := make(chan ElevLib.NextOrder)

	go floor_reached(floorReached, channels.Currentfloorupdate, floorSensor, newNextOrder)

	for{
		select{
			case order := <- channels.OrderChan:
				fmt.Println("FSM: Order =  ", order)

				newNextOrder <- order
				elev_set_motor_direction(order.Direction)

				
				for {
					select{
						case reachedFloor = <-floorReached:
							elev_set_motor_direction(0)
							breakbool = true

						case newOrder := <- channels.UpdateOrderChan:
							newNextOrder <- newOrder
							fmt.Println("FSM: UPDATED Order = ", newOrder)
					}
					if breakbool {
						break;
					}
				}
					
				breakbool = false

				channels.DeleteOrder <- reachedFloor

				if reachedFloor.ButtonType == ElevLib.BUTTON_COMMAND {
					setLightsOff <- []int{reachedFloor.ButtonType, reachedFloor.Floor,0}
				} else if reachedFloor.ButtonType == ElevLib.BUTTON_CALL_UP {
					setLightsOff <- []int{ ElevLib.BUTTON_CALL_UP, reachedFloor.Floor, 0 }
					setLightsOff <- []int{reachedFloor.ButtonType, reachedFloor.Floor,0}
				} else if reachedFloor.ButtonType == ElevLib.BUTTON_CALL_DOWN {
					setLightsOff <- []int{ reachedFloor.ButtonType, reachedFloor.Floor, 0 }
					setLightsOff <- []int{ ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0 }
				}
				
				elev_set_door_open_lamp(true)
				time.Sleep(3*time.Second)
				elev_set_door_open_lamp(false)

		}
		time.Sleep(time.Millisecond)
	}
}
