package Queue

import(
	"time"
	"fmt"
	".././ElevLib"
)


var localIp string


	
var InternalOrderChan chan ElevLib.MyOrder
var ExternalOrderChan chan ElevLib.MyOrder
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

func recievingthread(ready2recieve chan bool, ready chan bool) {
	for {
		select {
		case <-ready2recieve:
			ready <- true
		}
	}
	
}

func nextOrder(iOrder []int, eOrders [2][ElevLib.N_FLOORS]string, currentFloor int, dir int)(int,int){
	if len(iOrder)==0 {
		fmt.Println("dir ==: ",dir )
		if dir == 1 {
			for floor := currentFloor; floor < ElevLib.N_FLOORS ; floor++ {
				if eOrders[0][floor] == localIp { 
					return dir, floor
				} else {
					return dir, -1
				}
			}
		} else if dir == -1 {
			for floor := currentFloor; floor > -1 ; floor-- {
				if eOrders[1][floor] == localIp { 
					return dir, floor
				} else {
					return dir, -1

				}
			}
		} else if dir == 0 {
			for floor := currentFloor; floor < ElevLib.N_FLOORS ; floor++ {
				if eOrders[0][floor] == localIp { 
					fmt.Println(eOrders)
					return 1, floor
				} else {
					return 0, -1
				}
			}
			for floor := currentFloor; floor > -1 ; floor-- {
				if eOrders[1][floor] == localIp { 
					return -1, floor
				} else {
					return 0 , -1
				}
			}
		} 
	}

	tmpNextOrder := iOrder[0]
	if dir == 1{ 
		for floor := currentFloor; floor < ElevLib.N_FLOORS ; floor++ {
			if eOrders[0][floor] == localIp {
				return dir, floor
			}
		}
	} else if dir == -1 {
		for floor := currentFloor; floor > -1 ; floor-- {
			if eOrders[0][floor] == localIp {
				return dir, floor
			}
		}
	}else{
		if currentFloor < tmpNextOrder {
			return 1, tmpNextOrder
		}else{
			return -1, tmpNextOrder
		}
	}
	return dir, tmpNextOrder
}

// NEW SHIT

func Queue_manager(rcvFromEMChan chan ElevLib.NewReqFSM, sendReceipt2EM chan int, localIpsent string, setLightsOn chan []int, updCrntFlrAndDir chan []int, InternalOrderChan chan ElevLib.MyOrder, ExternalOrderChan chan ElevLib.MyOrder, deleteOrderChan chan []int) {
	
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
	time.Sleep(10*time.Millisecond)
	for{
		select{


			case order := <- InternalOrderChan:

				
				select{
					case newInternalOrder2check <- order:
					case <-time.After(10*time.Millisecond):
						fmt.Println("newInternalOrder2check TIMEOUT")
				}



				fmt.Println(internalOrders)

				time.Sleep(time.Millisecond)
				fmt.Sprintf("QUEUE: ","Internal order received on floor: %v", order.Floor)

				internalOrders = setInternalOrder(internalOrders, order.Floor , currentFloor , direction)

				setLightsOn <- []int{ElevLib.BUTTON_COMMAND, order.Floor, 1}
				fmt.Println(" ")
				/*fmt.Print("QUEUE: ", "Sending new info on internal orders to MASTER")
				sendInfoChan <- ElevLib.MyInfo{
					Ip: localIp,
					Dir: dir,
					CurrentFloor: current_floor,
					InternalOrders: internalOrders,
					}

				*/

			case order := <- ExternalOrderChan:

				select{
					case newExternalOrder2check <- order:
					case <-time.After(10*time.Millisecond):
						fmt.Println("newExternalOrder2check TIMEOUT")
				}

				
				fmt.Sprintf("QUEUE: ","External order received on floor: %v", order.Floor)
				externalOrders = setExternalOrder(externalOrders, order)
				setLightsOn <- []int{order.ButtonType, order.Floor, 1}
				fmt.Println(" ")
			


			case reqNewOrderFSM := <-rcvFromEMChan:
				fmt.Println("Queue: currentFloor = ", currentFloor)
				//direction = reqNewOrderFSM.Direction
				dir, nextFloor := nextOrder(internalOrders, externalOrders, currentFloor, direction) // OPPDATER nextOrder to return dir

				fmt.Println("Queue: Detected new request. Sending FSM to floor ", nextFloor, " in direction ", dir, " Current floor is: ", currentFloor)
				fmt.Println("INternalORders: ",internalOrders)
				if nextFloor != -1 {
					go checkForUpdOrders(reqNewOrderFSM.UpdateOrderChan,reqNewOrderFSM.KillThread, newInternalOrder2check, newExternalOrder2check, nextFloor, dir, localIp)

					// sending order to FSM
					reqNewOrderFSM.OrderChan <- [2]int{dir, nextFloor}
					sendReceipt2EM <- dir
					//fmt.Println("Queue: Ready to trigger on new cases")

				} else {
					reqNewOrderFSM.OrderChan <- [2]int{dir, nextFloor}

					sendReceipt2EM <- dir
					//fmt.Println("Queue: Ready to trigger on new cases")
				}
				fmt.Println(" ")
				

			case delOrder := <-deleteOrderChan:
				fmt.Println("QUEUE: ","DELETING ORDERS")
				internalOrders = internalOrders[1:]
				if delOrder[1] == 1 {
					externalOrders[ElevLib.BUTTON_CALL_UP][delOrder[0]] = ""
					
				
				} else if delOrder[1] == -1 {
					externalOrders[ElevLib.BUTTON_CALL_DOWN][delOrder[0]] = ""
					
				}

				sendReceipt2EM <- delOrder[0]
				fmt.Println(" ")
			case c := <-updCrntFlrAndDir:
				currentFloor = c[0]
				direction = c[1]
				fmt.Println("CurrentFloor Update to: ", currentFloor)
				fmt.Println(" ")
		}
		time.Sleep(time.Millisecond)
	}
	
}

func checkForUpdOrders(updateOrderChan chan int, killThread chan bool, iOrder chan ElevLib.MyOrder, eOrder chan ElevLib.MyOrder, nextFloor int, dir int, localIp string){

	currentNextFloor := nextFloor

	fmt.Println("checkForUpdOrders start")
	for {
		select{
			case <- killThread: // FSM tell checkForUpdOrders to terminate
				fmt.Println("checkForUpdOrders closing")
				fmt.Println(" ")
				return


			case order := <- iOrder:
				if order.Floor < currentNextFloor  && dir == 1 {
					currentNextFloor = order.Floor
					updateOrderChan <- currentNextFloor
				}else if order.Floor > currentNextFloor && dir == -1 {
					currentNextFloor = order.Floor
					updateOrderChan <- currentNextFloor
				}

			case order := <- eOrder:

				if order.Ip == localIp{
					if order.Floor < currentNextFloor && dir == 1 && order.ButtonType == ElevLib.BUTTON_CALL_UP {
						currentNextFloor = order.Floor
						updateOrderChan <- currentNextFloor
					}else if order.Floor > currentNextFloor && dir == -1 && order.ButtonType == ElevLib.BUTTON_CALL_DOWN{
						currentNextFloor = order.Floor
						updateOrderChan <- currentNextFloor
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