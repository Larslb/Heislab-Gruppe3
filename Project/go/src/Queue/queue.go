package Queue

import(
	"time"
	"ElevLib"
)


var localIp string

func setInternalOrder(iOrders []int, floor, currentFloor, dir int) ([]int) {
	if dir == 1 {
		if (floor - currentFloor < 0) {
			return append(iOrders,floor)
		} else if (floor - currentFloor == 0) {
			return insert(iOrders, floor, 0)
		}

		for i := 0; i < len(iOrders); i++ {
			if floor < iOrders[i]{
				return insert(iOrders, floor, i)
			}
		
		}
	} else if dir == -1 {
		if (floor - currentFloor > 0) {
			return append(iOrders,floor)
		} else if (floor - currentFloor == 0) {
			return insert(iOrders, floor, 0)
		}

		for i := 0; i < len(iOrders); i++{
			if floor > iOrders[i] {
				return insert(iOrders, floor, i)
			}
		}
	}
	return append(iOrders, floor)
}

func insert(orders []int ,floor, i int) ([]int) {
	tmp := make([]int, len(orders[:i]), len(orders)+1)
	copy(tmp, orders[:i])
	tmp = append(tmp, floor)
	tmp = append(tmp, internalOrders[i:]...)
	return tmpSlice
}

func setExternalOrder(eOrders [2][N_FLOORS]string, order ElevLib.MyOrder) ([2][N_FLOORS]string) {
	eOrders[order.ButtonType][order.Floor] = order.Ip
	return eOrders
}

func nextOrder(iOrder []int, eOrders [2][N_FLOORS]string, currentFloor int, dir int) int {
	
	if iOrder == [] {	
		if dir == 1 {
			for floor := currentFloor; floor < N_FLOORS ; floor++ {
				if eOrders[0][floor] == localIp { 
					return floor
				} else {
					return -1
				}
			}
		} else if dir == -1 {
			for floor := currentFloor; floor > -1 ; floor-- {
				if eOrders[1][floor] == localIp { 
					return floor
				} else {
					return -1
				}
			}
		} else if dir == 0 {
			for floor := currentFloor; floor < N_FLOORS ; floor++ {
				if eOrders[0][floor] == localIp { 
					return floor
				} else {
					return -1
				}
			}
			for floor := currentFloor; floor > -1 ; floor-- {
				if eOrders[1][floor] == localIp { 
					return floor
				} else {
					return -1
				}
			}
			
		} else {fmt.Println("ERROR: nextOrder bæsj ")
			return -1
			}
	}


	tmpNextOrder = iOrder[0]
	
	if dir == 1{ 
		if eOrders[0][currentFloor] == localIp {
			return currentFloor
		}
	} else if dir == -1 {
		if eOrders[0][currentFloor] == localIp {
			return currentFloor
		}
	}
	
	return tmpNextOrder

}

func Queue_manager(intrOrdChan chan int, extrOrdChan chan ElevLib.MyOrder, nextFloorChan chan int, deleteOrdFloorChan chan int, sendInfoChan chan ElevLib.MyInfo, currentFloorAndDirChan chan int, setLightsChan chan []int, localIpChan chan string){

	localIp = <- localIpChan

	dir := 0
	current_floor := -1

	internalOrders := []int{}
	externalOrders := [2][N_FLOORS]string{}
	
	for {
		select{
		case order := <- intrOrdChan:
			internalOrders = setInternalOrder(internalOrders, order, dir)

			setLightsChan <- []int{ElevLib.BUTTON_COMMAND, order, 1}

			sendInfoChan <- ElevLib.MyInfo{
				Ip: Network.LocalIp,
				Dir: dir,
				CurrentFloor: current_floor,
				InternalOrders: internalOrders,
				}
			
		case order := <- extrOrdChan:
			externalOrders = setExternalOrder(externalOrders, order)
			setLightsChan <- []int{order.ButtonType, order.Floor, 1}
		
		
		case tmpCurrent_floor := <-currentFloorAndDirChan:

			nextFloorChan <- nextOrder(internalOrders, externalOrders, tmpCurrent_floor, dir)
			
			if tmpDir := <-currentFloorAndDirChan; tmpDir != dir || tmpCurrent_floor != current_floor {
				dir = tmpDir
				current_floor = tmpCurrent_floor
				sendInfoChan <- ElevLib.MyInfo{
					Ip: Network.LocalIp,
					Dir: dir,
					CurrentFloor: current_floor,
					InternalOrders: internalOrders,
					}
			}
			
		case deleteOrder := <- deleteOrdFloorChan:
		
			internalOrders = internalOrders[1:]
			if tmpDir == 1 {
				externalOrders[ElevLib.BUTTON_CALL_UP][deleteOrder] = ""
				setLightsChan <- []int{ElevLib.BUTTON_CALL_UP, deleteOrder, 0}
				
				
			} else if tmpDir == -1 {
				externalOrders[ElevLib.BUTTON_CALL_DOWN][deleteOrder] = ""
				setLightsChan <- []int{ElevLib.BUTTON_CALL_DOWN, deleteOrder, 0}
			}
			
			// else ? ERROR HANDLING
			
			time.Sleep(10*time.Millisecond)
			setLightsChan <- []int{ElevLib.BUTTON_COMMAND, deleteOrder, 0}
			
			
			
			sendInfoChan <- ElevLib.MyInfo{
				Ip: Network.LocalIp,
				Dir: dir,
				CurrentFloor: current_floor,
				InternalOrders: internalOrders,
				}
		}
	}
}
