package main

import (
	"Queue"
	"Driver"
	"Network"
	"ElevLib"
	"time"
)






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
	var err bool
	//var localIp string
	
	floorSensor   := make(chan int)
	go Driver.ReadSensor(floorSensor)  // HVORDAN BRUKES DENNE I INIT??
	
	current_floor, err = Driver.Elev_init(floorSensor)
	if err {
		fmt.Println("Unable to initialize elevator!")
		//do something: return
	}
	
	localIp = ""  // Network.Init(localIpChan) 
	
	
	
	newExternalOrderChan := make(chan ElevLib.MyOrder)
	
	
	
	// NEW CHANNELS
	
	// COMMUNICATION BETWEEN EM AND FSM
	rcvNewReqFromFSMChan := make(chan ElevLib.NewReqFSM)
	sendReqStatusFSM     := make(chan int) // Only used when there have been no orders for a while
	orderHandled	   := make(chan ElevLib.MyOrder)
	
	
	// COMMUNICATION BETWEEN EM AND QUEUE
	sendReq2Queue 	   := make(chan ElevLib.NewReqFSM)
	receiptFromQueue     := make(chan ElevLib.MyOrder)
	startUpQueue 	   := make(chan ElevLib.StartQueue)
	
	
	// STARTUP PHASE, GO-ROUTINES
	
	go Queue.Queue_manager(sendReq2Queue, receiptFromQueue, startUpQueue)
	
	startUpQueue <- ElevLib.StartQueue{}
	time.Sleep(300*time.Millisecond)
	queueChannels <-startUpQueue
	
	
	go Driver.ReadElevPanel(newExternalOrderChan)
	go Driver.ReadFloorPanel(queueChannels.InternalOrderChan)
	
	checkDriverStatus := make(chan int)
	
	go Driver.FSM(rcvNewReqFromFSMChan, checkDriverStatus)
	
	// EVENT MANAGER
	for{
		select{
			case requestNewOrder := <-rcvNewReqFromFSMChan:
				sendReq2Queue <- requestNewOrder
				
				receipt := <- receiptFromQueue  // We wait for Queue to tell us where the elevetor is going
				direction = receipt.Dir
				
				
				
			case delOrder := <-orderHandledChan:
			case newExtOrd := <-newExternalOrderChan: // FORELØPIG BARE FOR Å TESTE 1 HEIS
				
				
		}
	}
	
	
}
