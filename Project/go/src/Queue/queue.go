package Queue


// PLAN 21.04.15

// 1. Skrive ferdig fsm
// 2. skrive ferdig setInternalOrders
// 3. Lage findNextFloor



const (
	N_BUTTONS int = 3
	N_FLOORS int = 4
)


type MyInfo struct {
	Ip string
	Dir int
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

// Ikke ferdige funksjoner

func setInternalOrder(iOrders []int, floor, dir int) ([]int) {
	if dir == 1 {
		// If floor - current position < 0, then append at back
		// if floor - current position = 0, open door
		for i := 0; i < len(iOrders); i++ {
			if floor < iOrders[i]{
				return insert(iOrders, floor, i)
			}
		
		}
	} else if dir == -1 {
		// If current position - floor < 0, then append at back
		// if current position - floor = 0, then open door
		for i := 0; i < len(iOrders); i++{
			if floor > iOrders[i] {
				return insert(iOrders, floor, i)
			}
		}
	}
	return append(iOrders, floor)
}

func queue_manager(intrOrd chan int, extrOrd chan myOrder, dirOrNF chan int, deleteOrdFloor chan int, reqInfo chan myInfo, currentFloor chan int){

	tmpDir := 0
	tmpCurrent_floor := -1

	internalOrders := []int{}
	externalOrders := [2][N_FLOORS]string{}   // [0][...] er opp bestillinger og [1][...] er ned bestillinger
	
	for {
		select{
		case order := <- intrOrd:
			internalOrders = setInternalOrder(internalOrders, order, tmpDir)
			// SEND OPPDATERT INFO
			
		case order := <- extrOrd:
			externalOrders = setExternalOrder(externalOrders, order)
		
		// DOBBELTSJEKK AT DENNE KANALEN ALDRI BLIR OVERFYLT
		case tmpCurrent_floor = <-currentFloor: // Bedre navnsetting på channel?
			dirOrNF <- nextOrder(internalOrders, externalOrders, dir)
			if dir := <-dirOrNF; dir != tmpDir{ // Hvis retningen har endret seg
				tmpDir = dir
				// SEND OPPDATERT INFO
				// reqInfo <- MyInfo{
				//		Ip: myIP,
				//		Dir: tmpDir,
				//		InternalOrders: internalOrders,
				//		}
			}
			
			
		case deleteOrder := <- deleteOrdFloor:
		
			internalOrders = internalOrders[1:]
			if tmpDir == 1 {
				externalOrders[BUTTON_CALL_UP][deleteOrder] = ""
			} else if tmpDir == -1 {
				externalOrders[BUTTON_CALL_DOWN][deleteOrder] = ""
			}
			
			// else ?
			
			// SEND OPPDATERT INFO
			// reqInfo <- MyInfo{
			//		Ip: myIP,
			//		Dir: tmpDir,
			//		InternalOrders: internalOrders,
			//		}
		}
	}
}

// Ferdige funksjoner

func insert (orders []int ,floor, i int) ([]int) {  // FUNKER BRA!
	tmp := make([]int, len(orders[:i]), len(orders)+1)
	copy(tmp, orders[:i])
	tmp = append(tmp, floor)
	tmp = append(tmp, internalOrders[i:]...)
	return tmpSlice
}

func setExternalOrder(eOrders [2][N_FLOORS]string, order MyOrder) ([2][N_FLOORS]string) {
	eOrders[order.ButtonType][order.Floor] = order.Ip
	return eOrders	
}

