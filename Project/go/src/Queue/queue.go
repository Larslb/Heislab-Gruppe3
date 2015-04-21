package Queue

const (
	N_BUTTONS int = 3
	N_FLOORS int = 4
)

type MyInfo struct {
	Ip string
	Dir int
	CurrentFloor int
	InternalOrders []int
}

type MyOrder struct {
	FromIp string
	// ToIp string
	ButtonType int
	Floor int
}

type MyElev struct {
	MessageType string
	Order MyOrder
	Info MyInfo
}

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

func insert (orders []int ,floor, i int) ([]int) {
	tmp := make([]int, len(orders[:i]), len(orders)+1)
	copy(tmp, orders[:i])
	tmp = append(tmp, floor)
	tmp = append(tmp, internalOrders[i:]...)
	return tmpSlice
}

func setExternalOrder(eOrders [2][N_FLOORS]string, order MyOrder) ([2][N_FLOORS]string) {  // HVORDAN ER DET MED BUTTONTYPE NÅR BUTTON_CALL_UP OG BUTTON_CALL_DOWN IKKE LIGGER I QUEUE???
	eOrders[order.ButtonType][order.Floor] = order.Ip
	return eOrders
}

func nextOrder(iOrder int, eOrders [2][N_FLOORS]string, currentFloor int, dir int) int {
	
	if iOrder == [] {	
		if dir == 1 {
			for floor := currentFloor; floor < N_FLOORS ; floor++ {
				if eOrders[0][floor] == Network.LocalIp { // HVORDAN GJØR VI DET HER?? HAR VI LAGRET MyIp?
					return floor
				} else {
					return -1 // NO ORDERS PENDING
				}
			}
		} else if dir == -1 {
			for floor := currentFloor; floor > -1 ; floor-- {
				if eOrders[1][floor] == Network.LocalIp { // HVORDAN GJØR VI DET HER?? HAR VI LAGRET MyIp?
					return floor
				} else {
					return -1 // NO ORDERS PENDING
				}
			}
		} else if dir == 0 {
			for floor := currentFloor; floor < N_FLOORS ; floor++ {
				if eOrders[0][floor] == Network.LocalIp { // HVORDAN GJØR VI DET HER?? HAR VI LAGRET MyIp?
					return floor
				} else {
					return -1 // NO ORDERS PENDING
				}
			for floor := currentFloor; floor > -1 ; floor-- {
				if eOrders[1][floor] == Network.LocalIp { // HVORDAN GJØR VI DET HER?? HAR VI LAGRET MyIp?
					return floor
				} else {
					return -1 // NO ORDERS PENDING
				}
			}
			
		} // else -> ERROR
	}


	tmpNextOrder = iOrder[0]
	
	if dir == 1{ 
		if eOrders[0][currentFloor] == Network.LocalIp {
			return currentFloor
		}
	} else if dir == -1 {
		if eOrders[0][currentFloor] == Network.LocalIp {
			return currentFloor
		}
	}
	
	return tmpNextOrder
}

func queue_manager(intrOrd chan int, extrOrd chan myOrder, dirOrNF chan int, deleteOrdFloor chan int, sendInfo chan myInfo, currentFloor chan int){

	dir := 0
	current_floor := -1

	internalOrders := []int{}
	externalOrders := [2][N_FLOORS]string{}   // [0][...] er opp bestillinger og [1][...] er ned bestillinger
	
	for {
		select{
		case order := <- intrOrd:
			internalOrders = setInternalOrder(internalOrders, order, dir)
			sendInfo <- MyInfo{
				Ip: Network.LocalIp,
				Dir: dir,
				CurrentFloor: current_floor,
				InternalOrders: internalOrders,
				}
			
		case order := <- extrOrd:
			externalOrders = setExternalOrder(externalOrders, order)
		
		
		case tmpCurrent_floor := <-currentFloor:
			
			if tmpDir := <-dirOrNF; tmpDir != dir || tmpCurrent_floor != current_floor {
				dir = tmpDir
				current_floor = tmpCurrent_floor
				sendInfo <- MyInfo{
					Ip: Network.LocalIp,
					Dir: dir,
					CurrentFloor: current_floor,
					InternalOrders: internalOrders,
					}
			}

			dirOrNF <- nextOrder(internalOrders, externalOrders, current_floor, dir)
			
			
		case deleteOrder := <- deleteOrdFloor:
		
			internalOrders = internalOrders[1:]
			if tmpDir == 1 {
				externalOrders[BUTTON_CALL_UP][deleteOrder] = ""
			} else if tmpDir == -1 {
				externalOrders[BUTTON_CALL_DOWN][deleteOrder] = ""
			}
			
			// else ?
			
			sendInfo <- MyInfo{
				Ip: Network.LocalIp,
				Dir: dir,
				CurrentFloor: current_floor,
				InternalOrders: internalOrders,
				}
		}
	}
}
