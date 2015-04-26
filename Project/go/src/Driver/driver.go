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
	}
	
	for {
		select {
			case go2Order = <- newOrder:
				fmt.Println("floor_reached: go2Order = ", go2Order)
			case current_floor := <- floorSensor:
				updateCurrentFloor <- current_floor

				elev_set_floor_indicator(current_floor)
				if current_floor == go2Order.Floor {
					floorReached <- go2Order
				}
		}
		time.Sleep(time.Millisecond)
	}
}

// NYTT




func Fsm( rcvChannelsFromQueue chan ElevLib.QM2FSMchannels, setLightsOff chan []int) {

	floorSensor := make(chan int)
	go ReadSensors(floorSensor)

	_,err := Elev_init(floorSensor)  // ER DET VITS Å HA RETUR PÅ CURRENT_FLOOR HER? SKRIVE OM INIT?
	if err {
		fmt.Println("FSM: Could not initiate elevator")
	}


	channels  := <- rcvChannelsFromQueue

	breakbool := false
	var reachedFloor ElevLib.NextOrder

	
	floorReached := make(chan ElevLib.NextOrder)  // Brukes i floor_reached func til å sende ut til drive
	newNextOrder := make(chan ElevLib.NextOrder)  // Brukes i floor_reached func til å oppdatere nxtOrder

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

				if reachedFloor.Direction == 1 && reachedFloor.ButtonType == ElevLib.BUTTON_CALL_UP{
					setLightsOff <- []int{ ElevLib.BUTTON_CALL_UP, reachedFloor.Floor, 0 }
					setLightsOff <- []int{ ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0 }
				} else if reachedFloor.Direction == -1 && reachedFloor.ButtonType == ElevLib.BUTTON_CALL_DOWN{
					setLightsOff <- []int{ ElevLib.BUTTON_CALL_DOWN, reachedFloor.Floor, 0 }
					setLightsOff <- []int{ ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0 }
				} else {
					setLightsOff <- []int{ ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0 }
				}

				fmt.Println("FSM: Door opening")
				elev_set_door_open_lamp(true)
				time.Sleep(3*time.Second)
				elev_set_door_open_lamp(false)
				fmt.Println("FSM: Door closing")

			// case ???

		}
		time.Sleep(time.Millisecond)
	}
}

/*func driver(channels ElevLib.OrderHandler2FSMchannels, killDrive chan bool, setLightsOff chan []int) {


	breakbool := false
	var reachedFloor ElevLib.NextOrder

	
	floorReached := make(chan ElevLib.NextOrder)  // Brukes i floor_reached func til å sende ut til drive
	newNextOrder := make(chan ElevLib.NextOrder)  // Brukes i floor_reached func til å oppdatere nxtOrder

	for{
		select{
			case order := <- channels.OrderChan:



				fmt.Println("driver: Order = ", order, " received")
				
				// Kan oppgraderes til å kjøre hele tiden ??
				go floor_reached(floorReached, channels.Currentfloorupdate, newNextOrder, order)

				elev_set_motor_direction(order.Direction)
				
				for {
					select{
						case reachedFloor = <-floorReached:
							elev_set_motor_direction(0)
							channels.FloorReachedChan <- true
							breakbool = true
						case newOrder := <- channels.UpdateOrderChan:
							newNextOrder <- newOrder
							fmt.Println("driver: UPDATED Order = ", order, " received")
					}
					if breakbool {
						break;
					}
				}
					
				breakbool = false

				channels.DeleteOrder <- reachedFloor

				if reachedFloor.Direction == 1 && reachedFloor.ButtonType == ElevLib.BUTTON_CALL_UP{
					setLightsOff <- []int{ ElevLib.BUTTON_CALL_UP, reachedFloor.Floor, 0 }
					setLightsOff <- []int{ ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0 }
				} else if reachedFloor.Direction == -1 && reachedFloor.ButtonType == ElevLib.BUTTON_CALL_DOWN{
					setLightsOff <- []int{ ElevLib.BUTTON_CALL_DOWN, reachedFloor.Floor, 0 }
					setLightsOff <- []int{ ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0 }
				} else {
					setLightsOff <- []int{ ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0 }
				}

				fmt.Println("driver: Door opening")
				elev_set_door_open_lamp(true)
				time.Sleep(3*time.Second)
				elev_set_door_open_lamp(false)

		}
		time.Sleep(time.Millisecond)
	}
}*/


// GAMMELT

/*
func FSM(sendReq2EM chan ElevLib.NewReqFSM, orderHandledChan chan ElevLib.NextOrder, setLightsOff chan []int, setlights chan bool, currentfloorupdate chan int) {

	rcvFromQueue := make(chan ElevLib.NextOrder)
	updFromQueue := make(chan ElevLib.NextOrder)
	killThreadChan := make(chan bool)
	var askNewOrder bool = true
	var breakbool bool = false
	
	// Used in goroutine func floorReached()
	newNextOrder     := make(chan ElevLib.NextOrder)
	floorReached := make(chan ElevLib.NextOrder)
	reachedFloor := ElevLib.NextOrder{
		ButtonType: ElevLib.BUTTON_COMMAND,
		Floor: -1,
		Direction: 0,
		}

	time.Sleep(10*time.Millisecond)
	fmt.Println("FSM: ", "Starting For Select Routine")


	for {
		select {
			case order := <-rcvFromQueue:

				if order.Floor != -1 {
					//fmt.Println("FSM: Driving ", order.Direction)
					
					elev_set_motor_direction(order.Direction)

					go floor_reached(floorReached, currentfloorupdate, newNextOrder, order)
				
					for {
						select{
							case reachedFloor = <-floorReached:
								elev_set_motor_direction(0)
								killThreadChan <- true // terminate Update routine in queue
								breakbool = true

							case newOrder := <- updFromQueue:
								newNextOrder<- newOrder
						}
						if breakbool {
							break;
						}
					}
					breakbool = false
					fmt.Println("FSM: Door opening")
					//elev_set_motor_direction(0)


					orderHandledChan <- ElevLib.NextOrder{
						ButtonType: reachedFloor.ButtonType,
						Floor:	reachedFloor.Floor,
						Direction: reachedFloor.Direction,
					}
					elev_set_door_open_lamp(true)  // MÅ FIKSES PÅ! HOLDES ÅPEN I 3 SEK
					time.Sleep(3*time.Second)
					elev_set_door_open_lamp(false)

				
					

					<-setlights
					if reachedFloor.Direction == 1 {
						setLightsOff <- []int{ElevLib.BUTTON_CALL_UP, reachedFloor.Floor, 0}  // MÅ ENDRE PÅ ELEV_SET_LIGHTS fra int til bool
						setLightsOff <- []int{ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0}
					} else if reachedFloor.Direction == -1 {
						setLightsOff <- []int{ElevLib.BUTTON_CALL_DOWN, reachedFloor.Floor, 0}
						setLightsOff <- []int{ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0}
					} else {
						setLightsOff <- []int{ElevLib.BUTTON_COMMAND, reachedFloor.Floor, 0}
					}

					//rdy2rcv <- true
					askNewOrder = true
				} else {
					askNewOrder = true
					time.Sleep(1*time.Second)
				}

			default:
				if askNewOrder {
					sendReq2EM <- ElevLib.NewReqFSM{
					OrderChan: rcvFromQueue,
					UpdateOrderChan: updFromQueue,
					KillThread: killThreadChan,
					//Current_floor: 0,  //no-care?
					//Direction: 0,      //no-care?
					}
					askNewOrder = false
				}
				
		}
		time.Sleep(time.Millisecond)
	}
}

*/