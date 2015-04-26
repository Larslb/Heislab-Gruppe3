package Queue

import(
	"time"
	"fmt"
	".././ElevLib"
)


var localIp string


	
//var InternalOrderChan chan ElevLib.MyOrder
//var ExternalOrderChan chan ElevLib.MyOrder
var SetLightsChan chan []int
var DeleteOrderChan chan []int
var setLightsChan chan []int


func setInternalOrder(iOrders []int, orderfloor, currentFloor, dir int) ([]int) {
	
	for  floor:= 0; floor < len(iOrders) ; floor++ {
		if iOrders[floor] == orderfloor {
			return iOrders
		}
	}

	if dir == 1{
		if currentFloor > orderfloor {
			return append(iOrders, orderfloor)

		}else if currentFloor == orderfloor {
				return insert(iOrders, orderfloor, 0)
		}

		for i := 0; i < len(iOrders); i++ {
			if orderfloor < iOrders[i]{
				return insert(iOrders, orderfloor, i)
			}
		}

	} else if dir == -1 {
		if currentFloor < orderfloor {
			return append(iOrders, orderfloor)

		}else if currentFloor == orderfloor {
				return insert(iOrders, orderfloor, 0)
		}

		for i := 0; i < len(iOrders); i++{
			if orderfloor > iOrders[i] {
				return insert(iOrders, orderfloor, i)
			}
		}
	}

	return append(iOrders, orderfloor)
}

func insert(orders []int ,floor, i int) ([]int) {
	tmp := make([]int, len(orders[:i]), len(orders)+1)
	copy(tmp, orders[:i])
	tmp = append(tmp, floor)
	tmp = append(tmp, orders[i:]...)
	return tmp
}

func setExternalOrder(eOrders [2][ElevLib.N_FLOORS]string, order ElevLib.MyOrder) ([2][ElevLib.N_FLOORS]string) {
	eOrders[order.ButtonType][order.Floor] = order.Ip
	return eOrders
}


func topDownSearch(eOrders [2][ElevLib.N_FLOORS]string, currentFloor int)(int,int) {
	
	var tmpFloor int
	var tmpDir int
	var boolVar bool = false

	for floor := ElevLib.N_FLOORS-2 ; floor > -1 ; floor-- {
		if eOrders[ElevLib.BUTTON_CALL_UP][floor] == localIp {
			if currentFloor < floor {
				tmpDir = 1
				tmpFloor = floor
				boolVar = true
			} else if currentFloor == floor{
				tmpDir = 0
				tmpFloor = floor
				boolVar = true	
			} else if currentFloor > floor && !boolVar{
				tmpDir = -1
				tmpFloor = floor
				boolVar = true
			}
		}
	}
	if boolVar{
		return tmpDir, tmpFloor
	} else {
		return 0, -1
	}
}

func bottomUpSearch(eOrders [2][ElevLib.N_FLOORS]string, currentFloor int)(int,int) {
	
	var tmpFloor int
	var tmpDir int
	var boolVar bool = false

	for floor := 1 ; floor < ElevLib.N_FLOORS; floor++{
		if eOrders[ElevLib.BUTTON_CALL_DOWN][floor] == localIp {
			if currentFloor > floor {
				tmpDir = -1
				tmpFloor = floor
				boolVar = true
			} else if currentFloor == floor {
				tmpDir = 0
				tmpFloor = floor
				boolVar = true
			} else if currentFloor < floor && !boolVar{	
				tmpDir = 1
				tmpFloor = floor
				boolVar = true
			}
		}
	}
	if boolVar{
		return tmpDir, tmpFloor
	} else {
		return 0, -1
	}
}

func search(eOrders [2][ElevLib.N_FLOORS]string, currentFloor int) (ElevLib.NextOrder){
	
	nxtOrder := ElevLib.NextOrder{
		ButtonType: ElevLib.BUTTON_COMMAND,
		Floor: -1,
		Direction: 0,
	}

	for floor := currentFloor ; floor < ElevLib.N_FLOORS ; floor++ {
		for Button:= 0; Button<ElevLib.N_BUTTONS-1; Button ++{
			if eOrders[Button][floor] == localIp{
				if floor == currentFloor {

					nxtOrder = ElevLib.NextOrder{
						ButtonType: Button,
						Floor: floor,
						Direction: 0,
					}
					return nxtOrder
				}

				nxtOrder = ElevLib.NextOrder{
					ButtonType: Button,
					Floor: floor,
					Direction: 1,
				}
				return nxtOrder

			}
		}
	}
	

	for floor := currentFloor; floor > -1; floor-- {
		for Button:= 0; Button<ElevLib.N_BUTTONS-1; Button ++{
			if eOrders[Button][floor] == localIp{
				if floor == currentFloor {
					nxtOrder = ElevLib.NextOrder{
						ButtonType: Button,
						Floor: floor,
						Direction: 0,
					}
					return nxtOrder
				}
				nxtOrder = ElevLib.NextOrder{
					ButtonType: Button,
					Floor: floor,
					Direction: -1,
				}
				return nxtOrder
			}
		}
	}

	return nxtOrder
}

	


func nextOrder(iOrder []int, eOrders [2][ElevLib.N_FLOORS]string, currentFloor int, dir int) (ElevLib.NextOrder) {
	
	


	var nxtOrder ElevLib.NextOrder
	var eTmpFloor int
	var eTmpDir int

	if currentFloor == -1{
		nxtOrder = ElevLib.NextOrder{
				ButtonType: -2,
				Floor: -1,
				Direction: 0,
		}
		return nxtOrder
	}

	if len(iOrder)==0 {
		if dir == 1 {
			fmt.Println("topDownSearch using: currentFloor = ", currentFloor, " direction = ", dir)
			eTmpDir, eTmpFloor = topDownSearch(eOrders, currentFloor)
			nxtOrder = ElevLib.NextOrder{
				ButtonType: ElevLib.BUTTON_CALL_UP,
				Floor: eTmpFloor,
				Direction: eTmpDir,
			}
			fmt.Println("topDownSearch result: nxtOrder.Floor = ", nxtOrder.Floor, " nxtOrder.Direction = ", nxtOrder.Direction)
			return nxtOrder

		} else if dir == -1 {
			fmt.Println("bottomUpSearch using: currentFloor = ", currentFloor, " direction = ", dir)
			eTmpDir, eTmpFloor = bottomUpSearch(eOrders, currentFloor)
			nxtOrder = ElevLib.NextOrder{
				ButtonType: ElevLib.BUTTON_CALL_DOWN,
				Floor: eTmpFloor,
				Direction: eTmpDir,
			}
			return nxtOrder

		} else if dir == 0 {
			fmt.Println("Search using: currentFloor = ", currentFloor, " direction = ", dir)
			nxtOrder = search(eOrders, currentFloor)
			fmt.Println("Result: nxtOrder.Floor = ", nxtOrder.Floor, " nxtOrder.Direction = ", nxtOrder.Direction)
			return nxtOrder

		} 
	}


	nxtOrder.Floor = iOrder[0]
	nxtOrder.ButtonType = ElevLib.BUTTON_COMMAND
	if currentFloor > nxtOrder.Floor{
		nxtOrder.Direction = -1
	} else if currentFloor < nxtOrder.Floor {
		nxtOrder.Direction = 1
	} else {
		nxtOrder.Direction = 0
	}


	if dir == 1{
		fmt.Println("topDownSearch using: currentFloor = ", currentFloor, " direction = ", dir)
		eTmpDir, eTmpFloor = topDownSearch(eOrders, currentFloor)
		if eTmpDir == 1{
			if eTmpFloor < nxtOrder.Floor{
				nxtOrder.Floor = eTmpFloor
				nxtOrder.Direction = eTmpDir
			}

		}else if eTmpDir == 0{
			nxtOrder.Floor = eTmpFloor
			nxtOrder.Direction = eTmpDir
		}
	} else if dir == -1 {
		fmt.Println("bottomUpSearch using: currentFloor = ", currentFloor, " direction = ", dir)
		eTmpDir, eTmpFloor = bottomUpSearch(eOrders, currentFloor)
		if eTmpDir == -1 {
			if eTmpFloor > nxtOrder.Floor {
				nxtOrder.Floor = eTmpFloor
				nxtOrder.Direction = eTmpDir
			}
		}else if eTmpDir == 0 {
			nxtOrder.Floor = eTmpFloor
			nxtOrder.Direction = eTmpDir
		}
	}
	fmt.Println("Result: nxtOrder.Floor = ", nxtOrder.Floor, " nxtOrder.Direction = ", nxtOrder.Direction)
	return nxtOrder
}

func deleteOrders(internalOrders []int, externalOrders [2][ElevLib.N_FLOORS]string, order ElevLib.NextOrder) ([]int, [2][ElevLib.N_FLOORS]string){

	if len(internalOrders) > 1 {
		internalOrders = internalOrders[1:]
	} else {
		internalOrders = []int{}
	}
	
	if order.ButtonType != ElevLib.BUTTON_COMMAND {
		externalOrders[order.ButtonType][order.Floor] = "x"
	}

	return internalOrders, externalOrders
}


// NYTT

func Queue_Manager(channels2fsm chan ElevLib.QM2FSMchannels, internalOrdersFromSensor chan ElevLib.MyOrder, externalOrdersFromSensor chan ElevLib.MyOrder, setLightsOn chan []int, localIpsent string){

	localIp = localIpsent

	currentFloor := -1
	direction    :=  0

	internalOrders := []int{}
	externalOrders := [2][ElevLib.N_FLOORS]string{}

	lastOrder         := ElevLib.NextOrder{}
	lastOrderFinished := true

	// COMMUNICATION WITH FSM
	orderChan     := make(chan ElevLib.NextOrder)      // SEND
	updOrderChan  := make(chan ElevLib.NextOrder)      // SEND
	deleteOrderChan  := make(chan ElevLib.NextOrder)   // MOTTA
	currentFloorUpdateChan := make(chan int)		   // MOTTA
	
	qm2fsmChannels := ElevLib.QM2FSMchannels{
						OrderChan: orderChan,
						UpdateOrderChan: updOrderChan,
						DeleteOrder: deleteOrderChan,
						Currentfloorupdate: currentFloorUpdateChan,
					}
	


	channels2fsm <- qm2fsmChannels

	for {
		select {
			case iOrder := <- internalOrdersFromSensor:
				
				if currentFloor != -1{
					iOrder.Ip = localIp
					setLightsOn <- []int{iOrder.ButtonType, iOrder.Floor, 1}


					internalOrders = setInternalOrder(internalOrders, iOrder.Floor , currentFloor , direction)

					fmt.Println("QUEUE: currentFloor = ", currentFloor, ", direction = ", direction)
					fmt.Println(" ")

					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

					fmt.Println("QUEUE: internalOrders = ", internalOrders, ", externalOrders = ", externalOrders)
					fmt.Println(" ")
					if lastOrderFinished {
						orderChan <- nxtOrder
						lastOrder = nxtOrder
						direction = nxtOrder.Direction
						lastOrderFinished = false
					} else if sendUpdate(lastOrder, nxtOrder) {   // MÅ HA EN CHECK UPDATE HER
						updOrderChan <- nxtOrder
						lastOrder = nxtOrder
						direction = nxtOrder.Direction
					}

					// SEND UPDATE INFO TO MASTER
				}


			case eOrder := <- externalOrdersFromSensor:
				

				if currentFloor != -1 {
					setLightsOn <- []int{eOrder.ButtonType, eOrder.Floor, 1}
					externalOrders = setExternalOrder(externalOrders, eOrder)

					fmt.Println("QUEUE: currentFloor = ", currentFloor, ", direction = ", direction)
					fmt.Println(" ")

					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

					fmt.Println("QUEUE: internalOrders = ", internalOrders, ", externalOrders = ", externalOrders)
					fmt.Println(" ")
					if lastOrderFinished {
						orderChan <- nxtOrder
						lastOrder = nxtOrder
						lastOrderFinished = false
						direction = nxtOrder.Direction

					} else if sendUpdate(lastOrder, nxtOrder) {  // MÅ HA EN CHECK UPDATE HER
						updOrderChan <- nxtOrder
						lastOrder = nxtOrder
						direction = nxtOrder.Direction
					}
				}
				

			case delOrder := <-deleteOrderChan:

				if delOrder == lastOrder {
					internalOrders, externalOrders = deleteOrders(internalOrders, externalOrders, delOrder)
					lastOrderFinished = true

					fmt.Println("QUEUE: internalOrders = ", internalOrders, ", externalOrders = ", externalOrders)
					fmt.Println(" ")
					// SEND UPDATE INFO TO MASTER
				} else {
					// ERROR
				}

				nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)
							
				if nxtOrder.Floor != -1 {
					orderChan <- nxtOrder
					direction = nxtOrder.Direction
					lastOrderFinished = false
				} else {
					fmt.Println("QUEUE: No more orders")
					fmt.Println(" ")
					direction = 0
				}

			case currentFloor = <- currentFloorUpdateChan:
				fmt.Println("QUEUE: currentFloor = ", currentFloor)
				fmt.Println(" ")
				// SEND UPDATE INFO TO MASTER
		}
		time.Sleep(time.Millisecond)
	}
}


func sendUpdate(lastOrder, newOrder ElevLib.NextOrder) bool {

	if lastOrder.Direction == 1  && newOrder.Direction == lastOrder.Direction {
		if newOrder.Floor < lastOrder.Floor {
			return true
		}
	} else if lastOrder.Direction == -1 && newOrder.Direction == lastOrder.Direction {
		if newOrder.Floor > lastOrder.Floor {
			return true
		}
	}

	return false
}


/*
	for {
		select {

			case iOrder := <- internalOrdersFromSensor:
				iOrder.Ip = localIp
				setLightsOn <- []int{iOrder.ButtonType, iOrder.Floor, 1}

				if orderHandlerIsAlive {
					iOrder2orderHandler <- iOrder

				} else {
					go orderHandler(sendChannels2fsm, sendChannels2oh, floorSensor)

					

					QM2orderHandlerChannels := ElevLib.Queue2OrderHandlerchannels{
							IOrdersChan: iOrder2orderHandler,
							EOrdersChan: eOrder2orderHandler,
							IsAliveChan: orderHandlerIsAliveChan,
					}

					fmt.Println("Queue: sending channels to orderHandler")      // NÅ HENGER orderHANDLER
					sendChannels2oh <- QM2orderHandlerChannels

					iOrder2orderHandler <- iOrder
				}



			case eOrder := <- externalOrdersFromSensor:
				setLightsOn <- []int{eOrder.ButtonType, eOrder.Floor, 1}

				if orderHandlerIsAlive {
					eOrder2orderHandler <- eOrder

				} else {
					go orderHandler(sendChannels2fsm, sendChannels2oh, floorSensor)

					QM2orderHandlerChannels := ElevLib.Queue2OrderHandlerchannels{
							IOrdersChan: iOrder2orderHandler,
							EOrdersChan: eOrder2orderHandler,
							IsAliveChan: orderHandlerIsAliveChan,
					}

					sendChannels2oh <- QM2orderHandlerChannels

					eOrder2orderHandler <- eOrder

				}

			case orderHandlerIsAlive = <- orderHandlerIsAliveChan:

		}
		time.Sleep(time.Millisecond)
	}

}

func orderHandler(channels2fsm chan ElevLib.OrderHandler2FSMchannels, channelsFromQM chan ElevLib.Queue2OrderHandlerchannels, floorSensor chan int) {

	
	fromQM := <- channelsFromQM
	fromQM.IsAliveChan <- true

	currentFloor := -1
	direction    :=  0

	internalOrders := []int{}
	externalOrders := [2][ElevLib.N_FLOORS]string{}

	lastOrder         := ElevLib.NextOrder{}
	lastOrderFinished := false
	floorReached      := false
	killOrderHandler  := false
	floor_reached_isAlive := false
	driver_isAlive    := false

	orderChan     := make(chan ElevLib.NextOrder)
	updOrderChan  := make(chan ElevLib.NextOrder)
	killGoRoutine := make(chan bool)
	fsmRdy4nextOrder := make(chan bool)
	floorReachedChan := make(chan bool)
	deleteOrderChan  := make(chan ElevLib.NextOrder)
	currentFloorUpdateChan := make(chan int)
	
	oh2fsmChannels := ElevLib.OrderHandler2FSMchannels{
						OrderChan: orderChan,
						UpdateOrderChan: updOrderChan,
						KillGoRoutine: killGoRoutine,
						FsmRdy4nextOrder: fsmRdy4nextOrder,
						FloorReachedChan: floorReachedChan,
						DeleteOrder: deleteOrderChan,
						Currentfloorupdate: currentFloorUpdateChan,
					}
	

	channels2fsm <- oh2fsmChannels

	for {
		select {
			case iOrder := <- fromQM.IOrdersChan:

				// FIKSE SET INTERNAL ORDERS OG SET EXTERNAL ORDERS 

				fmt.Println(" sagfhjhgfdsafjk      ", currentFloor)
				if currentFloor != -1{
					internalOrders := setInternalOrder(internalOrders, iOrder.Floor , currentFloor , direction)

					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

					if nxtOrder != lastOrder && !floorReached && driver_isAlive {  // SENDER NY NÅR VI TRYKKER PÅ SAMME KNAPP 2 GANGER. MEN IKKE 3
						updOrderChan <- nxtOrder
						lastOrder = nxtOrder
					}
					// SEND UPDATE INFO TO MASTER
				}


			case eOrder := <- fromQM.EOrdersChan:

				if currentFloor != -1 {
					externalOrders = setExternalOrder(externalOrders, eOrder)
					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

					if nxtOrder != lastOrder && !floorReached && driver_isAlive {
						updOrderChan <- nxtOrder
						lastOrder = nxtOrder
					}
				}

			case <-fsmRdy4nextOrder:
				floor_reached_isAlive = true
				driver_isAlive = true

				fmt.Println("orderHandler: Current Floor = ", currentFloor)
				nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

				fmt.Println("orderHandler: Direction = ", nxtOrder.Direction, ", Next Floor = ", nxtOrder.Floor, ", ButtonType = ", nxtOrder.ButtonType)
				fmt.Println(" ")
				if nxtOrder.Floor == -1 && lastOrderFinished {  // NO PENDING ORDERS TO HANDLE
					killOrderHandler = true
					killGoRoutine <- true
					driver_isAlive = false
				} else {
					orderChan <- nxtOrder
				}

			case delOrder := <-deleteOrderChan:

				if delOrder == lastOrder {
					internalOrders, externalOrders = deleteOrders(internalOrders, externalOrders, delOrder)
					lastOrderFinished = true
					// SEND UPDATE INFO TO MASTER
				} else {
					// ERROR
				}
				floor_reached_isAlive = false

			case floorReached = <- floorReachedChan:

			case currentFloor = <- floorSensor:
				fmt.Println("orderHandler: Current Floor = ", currentFloor)
				// MÅ BARE SENDE DERSOM floor_reached KJØRER
				if floor_reached_isAlive {
				fmt.Println("orderHandler: Sending update on current_floor")
				currentFloorUpdateChan <- currentFloor
				fmt.Println("UPDATE COMPLETE")
				fmt.Println(" ")

			}
				// SEND UPDATE INFO TO MASTER

		}

		if killOrderHandler {
			fromQM.IsAliveChan <- false
			// SEND UPDATE INFO TO MASTER
		}
	}
	fmt.Println("Queue: orderHandler killed")
	fmt.Println(" ")
}




// GAMMELT

func Queue_manager(rcvFromEMChan chan ElevLib.NewReqFSM, sendReceipt2EM chan int, localIpsent string, setLightsOn chan []int, updateCurrentFloor chan int, InternalOrderChan chan ElevLib.MyOrder, ExternalOrderChan chan ElevLib.MyOrder, deleteOrderChan chan ElevLib.NextOrder ) {
	
	currentFloor := -1
	direction := 0
	internalOrders := []int{}
	externalOrders := [2][ElevLib.N_FLOORS]string{}
	
	// Used by func CheckForUpdOrders
	localIp = localIpsent
	//rdy2rcvUpdate  := make(chan bool)
	//rdy2rcvUpdateChan := make(chan bool)
	newInternalOrder2check := make(chan ElevLib.MyOrder)
	newExternalOrder2check := make(chan ElevLib.MyOrder)

	//internalOrders, externalOrders := initializeOrders()
	fmt.Println("QUEUE:", "Going on")

	time.Sleep(10*time.Millisecond)
	for{
		select{


			case order := <- InternalOrderChan:

				
				select{
					case newInternalOrder2check <- order:
					case <-time.After(10*time.Millisecond):
						fmt.Println("newInternalOrder2check TIMEOUT")
				}



				fmt.Println(internalOrders, externalOrders)

				time.Sleep(time.Millisecond)
				fmt.Println("QUEUE: ","Internal order received on floor: ", order.Floor)

				internalOrders = setInternalOrder(internalOrders, order.Floor , currentFloor , direction)

				setLightsOn <- []int{ElevLib.BUTTON_COMMAND, order.Floor, 1}
				fmt.Println(" ")
				fmt.Print("QUEUE: ", "Sending new info on internal orders to MASTER")
				/*sendInfoChan <- ElevLib.MyInfo{
					Ip: localIp,
					Dir: direction,
					CurrentFloor: currentFloor,
					InternalOrders: internalOrders,
					}

				
				
			case order := <- ExternalOrderChan:

				order.Ip = localIp
				select{
					case newExternalOrder2check <- order:
					case <-time.After(10*time.Millisecond):
				//		fmt.Println("newExternalOrder2check TIMEOUT")
				}

				
				//fmt.Println("QUEUE: ","External order received on floor: ", order.Floor)
				externalOrders = setExternalOrder(externalOrders, order)
				setLightsOn <- []int{order.ButtonType, order.Floor, 1}
				fmt.Println(" ")
			
			
			case newExtorder := <- newExternalOrderChan:
				select{
					case newExternalOrder2check <- newExtorder:
					case <-time.After(10*time.Millisecond):
				//		fmt.Println("newExternalOrder2check TIMEOUT")
				}

				
				fmt.Println("QUEUE: ","External order received on floor: ", newExtorder.Floor)
				externalOrders = setExternalOrder(externalOrders, newExtorder)
				setLightsOn <- []int{newExtorder.ButtonType, newExtorder.Floor, 1}
				fmt.Println(" ")
			
			
			case reqNewOrderFSM := <-rcvFromEMChan:
				fmt.Println("Queue: currentFloor = ", currentFloor)


				// MÅ ENDRE PÅ RETURVARIABELEN
				nxtOrder  := nextOrder(internalOrders, externalOrders, currentFloor, direction)


				fmt.Println("Queue: Detected new request. Sending FSM to floor ", nxtOrder.Floor, " in direction ", nxtOrder.Direction, " Current floor is: ", currentFloor)
				fmt.Println("INternalORders: ",internalOrders, "externalOrders: ", externalOrders)


				// MÅ ENDRE HER OGSÅ
				direction = nxtOrder.Direction


				if nxtOrder.Floor != -1 {
					go checkForUpdOrders(reqNewOrderFSM.UpdateOrderChan,reqNewOrderFSM.KillThread, newInternalOrder2check, newExternalOrder2check, nxtOrder, localIp)

					
					reqNewOrderFSM.OrderChan <- nxtOrder
					sendReceipt2EM <- direction
					

				} else {
					reqNewOrderFSM.OrderChan <- nxtOrder

					sendReceipt2EM <- direction
					
				}
				fmt.Println(" ")
				

			case delOrder := <-deleteOrderChan:
				internalOrders, externalOrders = deleteOrders(internalOrders, externalOrders, delOrder)

				
				//sendReceipt2EM <- delOrder[0]
				fmt.Println(" ")


			case currentFloor = <-updateCurrentFloor:
				fmt.Println("CurrentFloor Update to: ", currentFloor)
				fmt.Println(" ")

		}
		time.Sleep(time.Millisecond)
	}
	
}

func checkForUpdOrders(updateOrderChan chan ElevLib.NextOrder, killThread chan bool, iOrder chan ElevLib.MyOrder, eOrder chan ElevLib.MyOrder, nxtOrder ElevLib.NextOrder, localIp string){

	currentNextOrder := nxtOrder

	fmt.Println("checkForUpdOrders start")
	for {
		select{
			case <- killThread: // FSM tell checkForUpdOrders to terminate
				fmt.Println("checkForUpdOrders closing")
				fmt.Println(" ")
				return


			case order := <- iOrder:
				if order.Floor > currentNextOrder.Floor  && currentNextOrder.Direction == 1 {
					currentNextOrder.Floor = order.Floor
					updateOrderChan <- currentNextOrder
				}else if order.Floor < currentNextOrder.Floor && currentNextOrder.Direction == -1 {
					currentNextOrder.Floor = order.Floor
					updateOrderChan <- currentNextOrder
				}

			case order := <- eOrder:

				if order.Ip == localIp{
					if order.Floor > currentNextOrder.Floor && currentNextOrder.Direction == 1 && order.ButtonType == ElevLib.BUTTON_CALL_UP {
						currentNextOrder.Floor = order.Floor
						currentNextOrder.ButtonType = order.ButtonType
						updateOrderChan <- currentNextOrder
					}else if order.Floor < currentNextOrder.Floor && currentNextOrder.Direction == -1 && order.ButtonType == ElevLib.BUTTON_CALL_DOWN{
						currentNextOrder.Floor = order.Floor
						currentNextOrder.ButtonType = order.ButtonType
						updateOrderChan <- currentNextOrder
					}
				}
		}
		time.Sleep(time.Millisecond)
	}
}

func initializeOrders() ([]int, [2][ElevLib.N_FLOORS]string){
	internalOrders := []int{}
	externalOrders := [2][ElevLib.N_FLOORS]string{}

	for i := 0; i < ElevLib.N_FLOORS; i++ {
		externalOrders[0][i] = ""
		externalOrders[1][i] = ""
		
	}

	return internalOrders, externalOrders
} */