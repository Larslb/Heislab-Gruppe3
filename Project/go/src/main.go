package main

import (
	"fmt"
	//"time"
)


const (
	BUTTON_DOWN    = 0
	BUTTON_UP      = 1
	BUTTON_COMMAND = 2
	
	N_FLOORS = 4
)

type myOrder struct{
	// fromIP string (f.eks: Master. Da vet alle slavene at master har gitt ordren til toIP slave)
	// toIP string
	
	IP string
	buttonType int
	floor int
}

type myInfo struct {
	Ip string
	Dir int
	InternalOrders []int
}


func queue_manager(intrOrd chan int, extrOrd chan myOrder, dirOrNF chan int, deleteOrdFloor chan int, reqInfo chan myInfo, msg chan string){
	
	
	
	// lastOrder := -1
	myIP := "0.0.1"
	direction := 1
	internalOrders := []int{}
	externalOrders := [2][N_FLOORS]string{}
	
	for {
		select{
		case order := <- intrOrd:
			internalOrders = setInternalOrder(internalOrders, order, direction)
			msg <- fmt.Sprintf("QM: Added floor %v to internalOrders\n\n", order)
			
		case order := <- extrOrd:
			externalOrders = setExternalOrder(externalOrders, order)
			msg <- fmt.Sprintf("QM: Added an order to externalOrders\n\n")
			
			
		//case dir := <- dirOrNF:
		//	dirOrNF <- nextOrder(internalOrders, externalOrders, dir)
		//	msg <- fmt.Sprintf("QM: returned nextOrder to FSM\n\n")
			
			
		case deleteOrder := <- deleteOrdFloor:
			//internalOrders, externalOrders = deleteOrderFloor(internalOrders, externalOrders, deleteOrder)
			msg <- fmt.Sprintf("QM: deleted order on floor %v\n\n", deleteOrder)
			
		case <-reqInfo:
			reqInfo <- myInfo{
					Ip: myIP,
					Dir: direction,
					InternalOrders: internalOrders,
					}
			msg <- fmt.Sprintf("QM: returned requested info to NM")
		}
	}
}

func sendOrders(intrOrd chan int, extrnOrd chan myOrder, reqInfo chan myInfo) {
	
	j := []int{1,4,2,3}
	
	myOrders := []myOrder{}
	
	myOrder1 := myOrder{
		IP: "0.0.1",
		buttonType: BUTTON_DOWN,
		floor: 3,
	}
	myOrder2 := myOrder{
		IP: "0.0.2",
		buttonType: BUTTON_UP,
		floor: 2,
	}
	
	myOrders = append(myOrders, myOrder1)
	myOrders = append(myOrders, myOrder2)
	
	
	info := make(map[int]myInfo)
	
	intCount := 0
	extCount := 0
	infoCount := 0
	for i := 0; i < 8; i ++ {
		if i < 4 {
			intrOrd <- j[intCount]
			intCount++
		} else if i >= 4 && i < 6 {
			extrnOrd <- myOrders[extCount]
			extCount++
		} else if i >= 6 && i < 8 {
			reqInfo <- myInfo{}
			info[infoCount] = <- reqInfo
			infoCount++
		}
		
	}
	
	//fmt.Println(info)
}

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

func insert (orders []int ,floor, i int) ([]int) {
	// Kanskje vi må passe på størrelsen til orders slik at vi vet at i finnes i orders??
	tmpSlice := orders[:i]
	tmpSlice = append(tmpSlice, floor)
	return append(tmpSlice, orders[i:]...)
}

func setExternalOrder(eOrders [2][N_FLOORS]string, order myOrder) ([2][N_FLOORS]string) {
	eOrders[order.buttonType][order.floor] = order.IP
	return eOrders
	
}


func main() {
	
	intrOrd     := make(chan int)
	//newExtrnOrd := make(chan myOrder)
	extrnOrd    := make(chan myOrder)
	
	//sensor := make(chan int)
	
	dirOrNextFloor := make(chan int)
	reqInfo	   := make(chan myInfo)
	deleteOrder    := make(chan int)
	
	msg := make(chan string)
	
	go queue_manager(intrOrd, extrnOrd, dirOrNextFloor, deleteOrder, reqInfo, msg)
	go sendOrders(intrOrd, extrnOrd, reqInfo)
	
	
	for{
		fmt.Println(<-msg)		
	}
}
