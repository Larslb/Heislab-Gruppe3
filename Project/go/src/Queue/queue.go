package Queue

import(
	"time"
	"fmt"
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
	tmp = append(tmp, orders[i:]...)
	return tmp
}

func setExternalOrder(eOrders [2][ElevLib.N_FLOORS]string, order ElevLib.MyOrder) ([2][ElevLib.N_FLOORS]string) {
	eOrders[order.ButtonType][order.Floor] = order.Ip
	return eOrders
}

func nextOrder(iOrder []int, eOrders [2][ElevLib.N_FLOORS]string, currentFloor int, dir int)int{
	if len(iOrder)==0 {
		if dir == 1 {
			for floor := currentFloor; floor < ElevLib.N_FLOORS ; floor++ {
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
			for floor := currentFloor; floor < ElevLib.N_FLOORS ; floor++ {
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

	tmpNextOrder := iOrder[0]
	
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

func Queue_manager(intrOrdChan chan ElevLib.MyOrder, extrOrdChan chan ElevLib.MyOrder, nextFloorChan chan int, deleteOrdFloorChan chan int, sendInfoChan chan ElevLib.MyInfo, currentFloorChan chan int, directionChan chan int, setLightsChan chan []int, localIpChan chan string){

	localIp = <- localIpChan

	dir := 0
	current_floor := -1

	internalOrders := []int{}
	externalOrders := [2][ElevLib.N_FLOORS]string{}
	
	for {
		select{
		case order := <- intrOrdChan:
			time.Sleep(30*time.Millisecond)
			fmt.Sprintf("QUEUE: ","Internal order received on floor: %v", order.Floor)
			internalOrders = setInternalOrder(internalOrders, order.Floor ,current_floor, dir)

			setLightsChan <- []int{ElevLib.BUTTON_COMMAND, order.Floor, 1}
			fmt.Print("QUEUE: ", "Sending new info on internal orders to MASTER")
			sendInfoChan <- ElevLib.MyInfo{
				Ip: localIp,
				Dir: dir,
				CurrentFloor: current_floor,
				InternalOrders: internalOrders,
				}
			
		case order := <- extrOrdChan:
			time.Sleep(30*time.Millisecond)
			fmt.Sprintf("QUEUE: ","External order received on floor: %v", order.Floor)
			externalOrders = setExternalOrder(externalOrders, order)
			setLightsChan <- []int{order.ButtonType, order.Floor, 1}
		
		
		case tmpCurrent_floor := <-currentFloorChan:
			time.Sleep(30*time.Millisecond)
			fmt.Println("QUEUE: ", "FSM requires next_order")
			next_floor := nextOrder(internalOrders, externalOrders, tmpCurrent_floor, dir)

			fmt.Println("QUEUE: ", "sending next_order to FSM")
			nextFloorChan <- next_floor
			
			fmt.Println("QUEUE: ", "waiting for new direction")
			if tmpDir := <-directionChan; tmpDir != dir || tmpCurrent_floor != current_floor {
				fmt.Println("QUEUE: ", "Sending new info on current_floor and direction to MASTER")
				dir = tmpDir
				current_floor = tmpCurrent_floor
				/*info := ElevLib.MyInfo{
					Ip: localIp,
					Dir: dir,
					CurrentFloor: current_floor,
					InternalOrders: internalOrders,
					}
				sendInfoChan <- info
				*/
			}
			fmt.Println("QUEUE: ", "dir =", dir, "current_floor = ", current_floor)
			
		case deleteOrder := <- deleteOrdFloorChan:
		
			internalOrders = internalOrders[1:]
			if dir == 1 {
				externalOrders[ElevLib.BUTTON_CALL_UP][deleteOrder] = ""
				setLightsChan <- []int{ElevLib.BUTTON_CALL_UP, deleteOrder, 0}
				
				
			} else if dir == -1 {
				externalOrders[ElevLib.BUTTON_CALL_DOWN][deleteOrder] = ""
				setLightsChan <- []int{ElevLib.BUTTON_CALL_DOWN, deleteOrder, 0}
			}
			
			// else ? ERROR HANDLING
			
			time.Sleep(10*time.Millisecond)
			setLightsChan <- []int{ElevLib.BUTTON_COMMAND, deleteOrder, 0}
			
			
			
			sendInfoChan <- ElevLib.MyInfo{
				Ip: localIp,
				Dir: dir,
				CurrentFloor: current_floor,
				InternalOrders: internalOrders,
				}
		}
	}
}
