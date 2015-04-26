package main

import (

	"fmt"
	//"./Driver"
	//"./Network"
	//"./ElevLib"
	//"time"
)


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

func insert(orders []int ,floor, i int) []int {
	tmp := make([]int, len(orders[:i]), len(orders)+1)
	copy(tmp, orders[:i])
	tmp = append(tmp, floor)
	tmp = append(tmp, orders[i:]...)
	return tmp
}


func main() {

	direction := 1
	currentFloor := 0
	internalOrders := []int{3,1,4,2,0,10}

	internalOrders = sortInDirection(internalOrders, currentFloor, direction)

	fmt.Println(internalOrders)

	/*
	Network.Init()
	fmt.Println("INTI!")
	newInfoChan := make(chan ElevLib.MyInfo)
	externalOrderChan := make(chan ElevLib.MyOrder) 
	newExternalOrderChan := make(chan ElevLib.MyOrder)
	readAndWriteAdresses := make(chan int, 1)
	masterChan := make(chan int,1)
	slaveChan := make(chan int,1)



	
	
	go Network.SendAliveMessageUDP()
	go Network.ReadAliveMessageUDP(readAndWriteAdresses)
	readAndWriteAdresses<-1
	Network.PrintAddresses()
	time.Sleep(time.Millisecond)
	go Network.SolvMaster(readAndWriteAdresses, masterChan, slaveChan)

	time.Sleep(time.Second)

	go Network.Network3(newInfoChan, externalOrderChan, newExternalOrderChan, masterChan, slaveChan)

	time.Sleep(time.Second)

	//Driver.Elev_init()
	//go ReeeadSensors()
	//time.Sleep(100*time.Second)

	fmt.Println("info sent")
	/*
	rcv := make(chan int)

	go send(rcv)

	t := time.Now()
	t2 := time.After()
	for {
		select{
				case <- rcv:
					fmt.Println("Received value")
				case <-time.After(2*time.Second):
					fmt.Println("Time out")
			}
	}
	
	time.Sleep(100*time.Second)
	fmt.Println("Main: Terminating")

	*/
}
