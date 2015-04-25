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
	fmt.Println("TopDownSearch is used!")
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
	fmt.Println("bottomUpSearch is used!")
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
	fmt.Println("Search is used!")
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
	
	// MÅ OPPDATERE RETURVARIABEL TIL Å VÆRE NEXTFLOOR TYPE
	var nxtOrder ElevLib.NextOrder
	var eTmpFloor int
	var eTmpDir int

	if currentFloor == -1{
		nxtOrder = ElevLib.NextOrder{
				ButtonType: ElevLib.BUTTON_COMMAND,
				Floor: -1,
				Direction: 0,
		}
		return nxtOrder
	}

	if len(iOrder)==0 {
		fmt.Println("dir ==: ",dir )

		if dir == 1 {
			eTmpDir, eTmpFloor = topDownSearch(eOrders, currentFloor)
			nxtOrder = ElevLib.NextOrder{
				ButtonType: ElevLib.BUTTON_CALL_UP,
				Floor: eTmpFloor,
				Direction: eTmpDir,
			}
			return nxtOrder

		} else if dir == -1 {
			eTmpDir, eTmpFloor = bottomUpSearch(eOrders, currentFloor)
			nxtOrder = ElevLib.NextOrder{
				ButtonType: ElevLib.BUTTON_CALL_DOWN,
				Floor: eTmpFloor,
				Direction: eTmpDir,
			}
			return nxtOrder

		} else if dir == 0 {
			return search(eOrders, currentFloor)

		} 
	}


	nxtOrder.Floor = iOrder[0]

	if currentFloor > nxtOrder.Floor{
		nxtOrder.Direction = -1
	} else if currentFloor < nxtOrder.Floor {
		nxtOrder.Direction = 1
	} else {
		nxtOrder.Direction = 0
	}


	if dir == 1{ 
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
	fmt.Println("TempnextOrder : ", nxtOrder.Floor, "tempDir", nxtOrder.Direction)

	return nxtOrder
}

func deleteOrders(internalOrders []int, externalOrders [2][ElevLib.N_FLOORS]string, order ElevLib.NextOrder) ([]int, [2][ElevLib.N_FLOORS]string){
	fmt.Println("QUEUE: ","DELETING ORDERS")


	// MÅ OPPDATERES I FORHOLD TIL NEXTORDER VARIABELEN

	if len(internalOrders) > 1 {
		internalOrders = internalOrders[1:]
	} else {
		internalOrders = []int{}
	}
	externalOrders[order.ButtonType][order.Floor] = "x"

	return internalOrders, externalOrders
}


// NYTT

func Queue_Manager(sendChannels2fsm ElevLib.OrderHandler2FSMchannels, setLightsOn chan []int, currentFloorChan chan int, localIpsent string){  // TA INN NØDVENDIGE CHANNELS

	localIp = localIpsent

	// CHANNELS COMMUNICATING WITH ORDERHANDLER
	orderHandlerIsAliveChan := make(chan bool)
	orderHandlerIsAlive := false
	iOrder2orderHandler := make(chan ElevLib.MyOrder)
	eOrder2orderHandler := make(chan ElevLib.MyOrder)

	// CHANNELS COMMUNICATING WITH DRIVER
	orderChan     := make(chan ElevLib.NextOrder)
	updOrderChan  := make(chan ElevLib.NextOrder)
	killGoRoutine := make(chan bool)
	fsmRdy4nextOrder := make(chan bool)
	floorReachedChan := make(chan bool)
	deleteOrderChan  := make(chan ElevLib.NextOrder)
	currentFloorUpdateChan := make(chan int)


	// DEFER CLOSE CHANNELS ?

	for {
		select {

			case iOrder := <- internalOrdersFromSensor:
				iOrder.Ip = localIp
				setLightsOn <- []int{iOrder.ButtonType, iOrder.Floor, 1}

				if orderHandlerIsAlive {
					iOrder2orderHandler <- iOrder

				} else {

					oh2fsmChannels := ElevLib.OrderHandler2FSMchannels{
						OrderChan: orderChan,
						UpdateOrderChan: updOrderChan,
						KillGoRoutine: killGoRoutine,
						Rdy4nextOrder: fsmRdy4nextOrder,
						FloorReachedChan: floorReachedChan,
						DeleteOrder: deleteOrderChan,
						Currentfloorupdate: currentFloorUpdateChan,
					}
					
					sendChannels2fsm <- oh2fsmChannels

					QM2orderHandlerChannels := ElevLib.Queue2OrderHandlerchannels{
							IOrdersChan: iOrder2orderHandler,
							EOrdersChan: eOrder2orderHandler,
							IsAliveChan: orderHandlerIsAliveChan,
					}

					go orderHandler(oh2fsmChannels, QM2orderHandlerChannels, currentFloorChan)
				}



			case eOrder := <- externalOrdersFromSensor:
				setLightsOn <- []int{eOrder.ButtonType, eOrder.Floor, 1}

				if orderHandlerIsAlive {
					eOrder2orderHandler <- eOrder

				} else {
					oh2fsmChannels := ElevLib.OrderHandler2FSMchannels{
						OrderChan: orderChan,
						UpdateOrderChan: updOrderChan,
						KillGoRoutine: killGoRoutine,
						Rdy4nextOrder: fsmRdy4nextOrder,
						FloorReachedChan: floorReachedChan,
						DeleteOrder: deleteOrderChan,
						Currentfloorupdate: currentFloorUpdateChan,
					}

					sendChannels2fsm <- oh2fsmChannels

					QM2orderHandlerChannels := ElevLib.Queue2OrderHandlerchannels{
							IOrdersChan: iOrder2orderHandler,
							EOrdersChan: eOrder2orderHandler,
							IsAliveChan: orderHandlerIsAliveChan,
					}

					go orderHandler(oh2fsmChannels, QM2orderHandlerChannels, currentFloorChan)
				}

			case orderHandlerIsAlive = <- orderHandlerIsAliveChan:

		}
		time.Sleep(time.Millisecond)
	}
}

func orderHandler(channels2fsm chan ElevLib.OrderHandler2FSMchannels, channelsFromQM ElevLib.Queue2OrderHandlerchannels, currentFloorChan chan int) {
	
	currentFloor := -1
	direction    :=  0

	internalOrders = []int{}
	externalOrders = [2][ElevLib.N_FLOORS]string{}

	lastOrder    := ElevLib.NextOrder{}
	floorReached := false
	killOrderHandler := false


	// DEFER CLOSE CHANNELS ?

	for {
		select {
			case iOrder := <- channelsFromQM.iOrderChan:

				if currentFloor != -1{
					internalOrders := setInternalOrder(internalOrders, iOrder.Floor , currentFloor , direction)
					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

					if nxtOrder != lastOrder && !floorReached {
						channels2fsm.UpdateOrderChan <- nxtOrder
						lastOrder = nxtOrder
					}
					// SEND UPDATE INFO TO MASTER
				}


			case eOrder := <- channelsFromQM.eOrdersChan:

				if currentFloor != -1 {
					externalOrders = setExternalOrder(externalOrders, eOrder)
					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

					if nxtOrder != lastOrder && !floorReached{
						channels2fsm.UpdateOrderChan <- nxtOrder
						lastOrder = nxtOrder
					}
				}


			case <-channelsFromQM.fsmRdy4NextOrder:

				nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, direction)

				if nxtOrder.Floor == -1 {  // NO PENDING ORDERS TO HANDLE
					killOrderHandler = true
					channels2fsm.KillGoRoutine <- true
				} else {
					channels2fsm.OrderChan <- nxtOrder
				}

			case delOrder := channelsFromQM.DeleteOrderChan:

				if delOrder == lastOrder {
					internalOrders, externalOrders = deleteOrders(internalOrders, externalOrders, delOrder)
					// SEND UPDATE INFO TO MASTER

				} else {
					// ERROR
				}

			case floorReached = <- channelsFromQM.FloorReachedChan:

			case currentFloor = <- currentFloorChan:
				channels2fsm.Currentfloorupdate <- currentFloor
				// SEND UPDATE INFO TO MASTER

		}

		if killOrderHandler {
			channelsFromQM.IsAliveChan <- false
			// SEND UPDATE INFO TO MASTER
		}
	}
	
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

				*/
				
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
			
			/*
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
			
			*/
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
}