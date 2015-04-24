package main

import (
	"./Queue"
	"./Driver"
	"./Network"
	"./ElevLib"
	"time"
	"fmt"
)


// HVA SOM MÅ GJØRES I MORGEN
// 1. sørge for at queue får tilgang til current_floor



func main() {
	/*
	//killAllConnections()

	localIpChan := make(chan string)
	newExternalOrderChan := make(chan ElevLib.MyOrder)
	externalOrderChan := make(chan ElevLib.MyOrder)
	internalOrderChan  := make(chan ElevLib.MyOrder)
	nextFloorChan := make(chan int)
	deleteOrderOnFloorChan := make(chan int)
	newInfoChan := make(chan ElevLib.MyInfo)
	currentFloorChan := make(chan int)
	directionChan := make(chan int)
	setLightsChan := make(chan []int)


	go Queue.Queue_manager(internalOrderChan, externalOrderChan, nextFloorChan, deleteOrderOnFloorChan, newInfoChan, currentFloorChan, directionChan, setLightsChan, localIpChan)
	
	//Network.EventManager_NetworkStuff(newInfoChan, externalOrderChan, newExternalOrderChan)
	//defer close channels	

	Network.Init(localIpChan)
	go Network.SendAliveMessageUDP()
	go Network.ReadAliveMessageUDP()
	Network.PrintAddresses()
	go Driver.ReadElevPanel(internalOrderChan)
	go Driver.ReadFloorPanel(newExternalOrderChan)
	time.Sleep(1*time.Second)
	go Driver.Fsm(nextFloorChan, deleteOrderOnFloorChan, currentFloorChan, directionChan, setLightsChan)

	Network.Network(newInfoChan, externalOrderChan , newExternalOrderChan)

	*/
	
	// INIT
	current_floor := -1
	direction := 0
	//err := false
	
	localIp,_ := Network.GetLocalIP() 
	fmt.Println("localIp= ",localIp)
	
	
	//newExternalOrderChan := make(chan ElevLib.MyOrder)
	
	Driver.Io_init()
	
	// NEW CHANNELS
	sensorchan := make(chan int)
	// COMMUNICATION BETWEEN EM AND FSM
	rcvNewReqFromFSMChan := make(chan ElevLib.NewReqFSM)
	//checkDriverStatus    := make(chan int) // Only used when there have been no orders for a while
	orderHandledChan	 := make(chan int)
	setlights 			 := make(chan bool)
	setLightsOff := make(chan []int)
	currentfloorupdateFSM	:= make(chan int)
	
	
	// COMMUNICATION BETWEEN EM AND QUEUE
	sendReq2Queue 	     := make(chan ElevLib.NewReqFSM)
	receiptFromQueue     := make(chan int)
	updCrntFlrAndDir := make(chan []int)
	setLightsOn := make(chan []int)
	deleteOrderChan := make(chan []int)

	//FoR QUEUE AND DRIVER
	InternalOrderChan := make(chan ElevLib.MyOrder)
	ExternalOrderChan := make(chan ElevLib.MyOrder)


	go Driver.ReadSensors(sensorchan)
	//rcvCurrentFloorQueue := make(chan chan int)
	current_floor,_ = Driver.Elev_init(sensorchan)
	fmt.Println("Elev_init Done: current_floor = ", current_floor, " and direction = ", direction)
	
	// STARTUP PHASE, GO-ROUTINES
	//go Driver.ReadSensors(sensorchan)
	go Driver.ReadElevPanel(InternalOrderChan)
	go Driver.ReadFloorPanel(ExternalOrderChan)
	go Queue.Queue_manager(sendReq2Queue, receiptFromQueue, localIp, setLightsOn, updCrntFlrAndDir, InternalOrderChan, ExternalOrderChan, deleteOrderChan)

	
	//fmt.Println("Elev_init Done: current_floor = ", current_floor, " and direction = ", direction)
	


	go Driver.FSM(rcvNewReqFromFSMChan, orderHandledChan, setLightsOff, setlights, currentfloorupdateFSM)
	go Driver.SetLights(setLightsOn, setLightsOff)
	time.Sleep(10*time.Millisecond)

	// EVENT MANAGER
	for{
		select{
			case requestNewOrder := <-rcvNewReqFromFSMChan:
				//requestNewOrder.Current_floor = current_floor
				//requestNewOrder.Direction = direction
				//fmt.Println("MAIN: ", "order is now: currentFloor: ", requestNewOrder.Current_floor, " Direction: ", requestNewOrder.Direction )
				sendReq2Queue <- requestNewOrder
				
				receipt := <- receiptFromQueue  // We wait for Queue to tell us where the elevetor is going
				direction = receipt
				//fmt.Println("MAIN: Ready to trigger on new cases")
				
				
				
			case floor := <-orderHandledChan:
				deleteOrderChan <- []int{floor, direction}

				receipt := <- receiptFromQueue  // Trenger egentlig ikke å ta imot
				fmt.Println("Order on floor ", receipt, " in direction ", direction, " was deleted")

				setlights <- false


			//case newExtOrd := <-newExternalOrderChan: // FORELØPIG BARE FOR Å TESTE 1 HEIS
			//	newExtOrd.Ip = localIp
			//	ExternalOrderChan <- newExtOrd

			case current_floor = <-sensorchan:
				fmt.Println("CurrentFloor is: ", current_floor)
				updCrntFlrAndDir <- []int{current_floor,direction}
				select{
					case currentfloorupdateFSM <- current_floor:
					case <-time.After(1*time.Second):
				}
				
		}
		time.Sleep(time.Millisecond)
	}
}
