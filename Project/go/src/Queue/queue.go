package Queue

import(
	"time"
	//"fmt"
	".././ElevLib"
)

var localIp string

func Queue_Manager( channels2fsm chan ElevLib.QM2FSMchannels, internalOrdersFromSensor chan ElevLib.MyOrder, externalOrdersFromMaster chan ElevLib.MyOrder, setLightsOn chan []int, localIpsent string, newInfo chan ElevLib.MyInfo, orderdeletion chan ElevLib.MyOrder, orderDelFromMaster chan ElevLib.MyOrder){

	localIp = localIpsent

	currentFloor := -1
	orderdirection    :=  0

	internalOrders := []int{}
	externalOrders := [2][ElevLib.N_FLOORS]string{}

	lastOrder         := ElevLib.NextOrder{}
	lastOrderFinished := true

	// COMMUNICATION WITH FSM
	orderChan     := make(chan ElevLib.NextOrder)    
	updOrderChan  := make(chan ElevLib.NextOrder)     
	deleteOrderChan  := make(chan ElevLib.NextOrder)  
	currentFloorUpdateChan := make(chan int)		  
	
	qm2fsmChannels := ElevLib.QM2FSMchannels{
						OrderChan: orderChan,
						UpdateOrderChan: updOrderChan,
						DeleteOrder: deleteOrderChan,
						Currentfloorupdate: currentFloorUpdateChan,
					}
	
	channels2fsm<- qm2fsmChannels
	
	sendInfo := ElevLib.MyInfo {
		Ip: localIp,
		Dir: orderdirection,
		CurrentFloor: currentFloor,
		InternalOrders: internalOrders,
	}
	
	

	for {
		select {
			case iOrder := <- internalOrdersFromSensor:
				
				if currentFloor != -1 && notInInternalOrders(internalOrders, iOrder.Floor){
					setLightsOn <- []int{iOrder.ButtonType, iOrder.Floor, 1}
					internalOrders = setInternalOrder(internalOrders, iOrder.Floor , currentFloor , orderdirection)


					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, orderdirection)

					if lastOrderFinished {
						orderChan <- nxtOrder
						lastOrder = nxtOrder
						if orderdirection != nxtOrder.Direction && len(internalOrders) > 1{
							internalOrders = sortInDirection(internalOrders, currentFloor, nxtOrder.Direction)
						}
						orderdirection = nxtOrder.Direction
						lastOrderFinished = false
					} else if sendUpdate(lastOrder, nxtOrder) {
						updOrderChan <- nxtOrder
						lastOrder = nxtOrder
						if orderdirection != nxtOrder.Direction && len(internalOrders) > 1 {
						internalOrders = sortInDirection(internalOrders, currentFloor, nxtOrder.Direction)
						}
						orderdirection = nxtOrder.Direction
					}

					sendInfo.Ip = localIp
					sendInfo.Dir = orderdirection
					sendInfo.CurrentFloor = currentFloor
					sendInfo.InternalOrders = internalOrders
					newInfo <- sendInfo
				}


			case eOrder := <- externalOrdersFromMaster:

				if currentFloor != -1 && notInExternalOrders(externalOrders, eOrder) {
					setLightsOn <- []int{eOrder.ButtonType, eOrder.Floor, 1}
					externalOrders = setExternalOrder(externalOrders, eOrder)

					nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, orderdirection)

					if lastOrderFinished {
						orderChan <- nxtOrder
						lastOrder = nxtOrder
						lastOrderFinished = false
						if orderdirection != nxtOrder.Direction && len(internalOrders) > 1{
							internalOrders = sortInDirection(internalOrders, currentFloor, nxtOrder.Direction)
						}
						orderdirection = nxtOrder.Direction

					} else if sendUpdate(lastOrder, nxtOrder) {  
						updOrderChan <- nxtOrder
						lastOrder = nxtOrder
						if orderdirection != nxtOrder.Direction && len(internalOrders) > 1{
							internalOrders = sortInDirection(internalOrders, currentFloor, nxtOrder.Direction)
							}
						}
						orderdirection = nxtOrder.Direction
					}


			case delOrder := <-deleteOrderChan:

				if delOrder == lastOrder {

					internalOrders, externalOrders = deleteOrders(internalOrders, externalOrders, delOrder)
					lastOrderFinished = true

					if delOrder.ButtonType != ElevLib.BUTTON_COMMAND {
						order := ElevLib.MyOrder{
							Ip: localIp,
							ButtonType: delOrder.ButtonType,
							Floor: delOrder.Floor,
							Set: false,
						}
						orderdeletion <- order
					} 
					sendInfo.Ip = localIp
					sendInfo.Dir = orderdirection
					sendInfo.CurrentFloor = currentFloor
					sendInfo.InternalOrders = internalOrders
					newInfo <- sendInfo
				}

				nxtOrder := nextOrder(internalOrders, externalOrders, currentFloor, orderdirection)

				if nxtOrder.Floor != -1 {
					orderChan <- nxtOrder
					if orderdirection != nxtOrder.Direction && len(internalOrders) > 1 {
						internalOrders = sortInDirection(internalOrders, currentFloor, nxtOrder.Direction)
					}
					orderdirection = nxtOrder.Direction
					lastOrder = nxtOrder
					lastOrderFinished = false
				} else {
					orderdirection = 0
				}

			case orderdelete := <- orderDelFromMaster:
				externalOrders[orderdelete.ButtonType][orderdelete.Floor] = ""

			case currentFloor = <- currentFloorUpdateChan:
				sendInfo.Ip = localIp
				sendInfo.Dir = orderdirection
				sendInfo.CurrentFloor = currentFloor
				sendInfo.InternalOrders = internalOrders
				newInfo <- sendInfo
		}
		time.Sleep(time.Millisecond)
	}
}


func sortInDirection(iOrders []int, currentFloor int, direction int) []int {

	tmpOrders := []int{}
	tmpOrders = append(tmpOrders, iOrders[0])
	for i := 1; i < len(iOrders) ; i++ {
		tmpOrders = setInternalOrder(tmpOrders, iOrders[i], currentFloor, direction)
	}

	return tmpOrders
}

func setInternalOrder(iOrders []int, orderfloor, currentFloor, dir int) ([]int) {
	
	if dir == 1{
		if currentFloor > orderfloor {
			return append(iOrders, orderfloor)
		}else if currentFloor == orderfloor {
				return insert(iOrders, orderfloor, 0)
		}
		for i := 0; i < len(iOrders); i++ {
			if orderfloor < iOrders[i] || currentFloor > iOrders[i] {
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
			if orderfloor > iOrders[i] || currentFloor < iOrders[i] {
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

func topDownSearch(eOrders [2][ElevLib.N_FLOORS]string, currentFloor int) ElevLib.NextOrder  {
	
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
		nxtOrder := ElevLib.NextOrder{
			ButtonType: ElevLib.BUTTON_CALL_UP,
			Floor: tmpFloor,
			Direction: tmpDir,
		}

		return nxtOrder
	} else {
		nxtOrder := ElevLib.NextOrder{
			ButtonType: ElevLib.BUTTON_CALL_UP,
			Floor: -1,
			Direction: 0,
		}
		return nxtOrder
	}
}

func bottomUpSearch(eOrders [2][ElevLib.N_FLOORS]string, currentFloor int) ElevLib.NextOrder {
	
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
		nxtOrder := ElevLib.NextOrder{
			ButtonType: ElevLib.BUTTON_CALL_DOWN,
			Floor: tmpFloor,
			Direction: tmpDir,
		}
		return nxtOrder

	} else {
		nxtOrder := ElevLib.NextOrder{
			ButtonType: ElevLib.BUTTON_CALL_DOWN,
			Floor: -1,
			Direction: 0,
		}
		return nxtOrder
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

func notInInternalOrders(iOrders []int, orderfloor int) bool {
	for  floor:= 0; floor < len(iOrders) ; floor++ {
		if iOrders[floor] == orderfloor {
			return false
		}
	}
	return true
}

func notInExternalOrders(eOrders [2][ElevLib.N_FLOORS]string, orderfloor ElevLib.MyOrder) bool {
	if orderfloor.Ip == eOrders[orderfloor.ButtonType][orderfloor.Floor] {
		return false
	}
	return true
}

func nextOrder(iOrder []int, eOrders [2][ElevLib.N_FLOORS]string, currentFloor int, dir int) (ElevLib.NextOrder) {
	
	var nxtOrder ElevLib.NextOrder

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
			nxtOrder = topDownSearch(eOrders, currentFloor)
			return nxtOrder

		} else if dir == -1 {
			nxtOrder =  bottomUpSearch(eOrders, currentFloor)
			return nxtOrder

		} else if dir == 0 {
			nxtOrder = search(eOrders, currentFloor)
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
		tmpNextOrder := topDownSearch(eOrders, currentFloor)
		if tmpNextOrder.Direction == 1 && tmpNextOrder.Floor < nxtOrder.Floor{
				nxtOrder.Floor = tmpNextOrder.Floor
				nxtOrder.Direction = tmpNextOrder.Direction
		}
	} else if dir == -1 {
		tmpNextOrder := bottomUpSearch(eOrders, currentFloor)
		if tmpNextOrder.Direction == -1 && tmpNextOrder.Floor > nxtOrder.Floor {
				nxtOrder.Floor = tmpNextOrder.Floor
				nxtOrder.Direction = tmpNextOrder.Direction
		}
	}
	return nxtOrder
}

func deleteOrders(internalOrders []int, externalOrders [2][ElevLib.N_FLOORS]string, order ElevLib.NextOrder) ([]int, [2][ElevLib.N_FLOORS]string){

	if order.ButtonType == ElevLib.BUTTON_COMMAND {
		if len(internalOrders) > 1{
			internalOrders = internalOrders[1:]
		} else {
			internalOrders = []int{}
		}
	}else{
		if len(internalOrders) != 0{
			if internalOrders[0] == order.Floor{
				internalOrders = internalOrders[1:]
			}
		}
		externalOrders[order.ButtonType][order.Floor] = ""
	}
	return internalOrders, externalOrders
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
