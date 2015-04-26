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

	fmt.Println("NFE")

	// DETTE ER NYTT (Se lenger ned for det gamle)

	localIp,_ := Network.GetLocalIP()

	fmt.Println("MAIN: localIp= ",localIp)
		Network.Init()
	errorDetection := make(chan bool) // Brukes ikke til noe enda. Må finne en passende type.  string??
	
	//sensorchan 		  := make(chan int)
	externalOrderChan := make(chan ElevLib.MyOrder)
	internalOrderChan := make(chan ElevLib.MyOrder)


	//NETWORKCHANNELS
	newInfoChan := make(chan ElevLib.MyInfo) 
	newPanelOrderChan := make(chan ElevLib.MyOrder)
	readAndWriteAdresses := make(chan int, 1)
	masterChan := make(chan int)
	slaveChan := make(chan int)

	qM2FSM := make(chan ElevLib.QM2FSMchannels)
	//oh2fsmChans       := make(chan ElevLib.OrderHandler2FSMchannels)
	//send2fsm		  := make(chan ElevLib.OrderHandler2FSMchannels)

	setLightsOn		  := make(chan []int)
	setLightsOff      := make(chan []int)
	//currentFloorChan  := make(chan int)
	//currentFloorUdateFSM := make(chan int)


	//NETWORK INIT!
	go Network.SendAliveMessageUDP()
	go Network.ReadAliveMessageUDP(readAndWriteAdresses)
	readAndWriteAdresses<-1
	Network.PrintAddresses()
	time.Sleep(time.Second)
	go Network.SolvMaster(readAndWriteAdresses, masterChan, slaveChan)

	time.Sleep(time.Second)

	go Network.Network3(newInfoChan, externalOrderChan, newPanelOrderChan, masterChan, slaveChan)

	time.Sleep(time.Second)
	
	go Driver.ReadElevPanel(internalOrderChan)
	go Driver.ReadFloorPanel(newPanelOrderChan)
	//go Driver.ReadSensors(sensorchan)
	go Driver.SetLights(setLightsOn, setLightsOff)

	/*current_floor, err := Driver.Elev_init(sensorchan)
	if err {
		fmt.Println("MAIN: Could not initiate elevator!")
	}*/

	go Queue.Queue_Manager(qM2FSM, internalOrderChan, externalOrderChan, setLightsOn, localIp, newInfoChan)
	go Driver.Fsm(qM2FSM, setLightsOff)

	//fmt.Println("MAIN: Ready to synchronize Queue and Fsm")
	//fmt.Println(" ")

	
	for{
		select{

			case err := <-errorDetection:
				fmt.Println(err)

			// Trenger jeg denne lenger??

			//case initiateOrderHandling := <-oh2fsmChans:   // BRUKES IKKE
			//	fmt.Println("MAIN: Received channels from orderHandler. Sending them to Fsm.")
			//	fmt.Println(" ")
			//	send2fsm <- initiateOrderHandling
			
			/*case current_floor := <-sensorchan:
				fmt.Println("MAIN: currentFloor is ", current_floor)
				
				// GÅR DET FOR TREIGT Å OPPDATERE FSM GJENNOM ORDERHANDLER??
				currentFloorChan <- current_floor
				fmt.Println("MAIN: Current Floor updated throughout")
			*/
			case <-time.After(200*time.Millisecond):
				time.Sleep(10*time.Second)   // zzzzzzZZZZZZZZZZZZZZZZZzzzzz
				
		}
		time.Sleep(10*time.Millisecond) // zzzzzzzZZZZZZZZZZZZzzzzz
	}




	// DETTE ER GAMMELT



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
	
	/*

	// INIT
	current_floor := -1
	direction := 0
	//err := false
	
	localIp,_ := Network.GetLocalIP() 
	fmt.Println("localIp= ",localIp)
	
	
	//newExternalOrderChan := make(chan ElevLib.MyOrder)
	
	Driver.Io_init()
	Network.Init()

	
	// NEW CHANNELS
	sensorchan := make(chan int)
	// COMMUNICATION BETWEEN EM AND FSM
	rcvNewReqFromFSMChan := make(chan ElevLib.NewReqFSM)
	//checkDriverStatus    := make(chan int) // Only used when there have been no orders for a while
	orderHandledChan	 := make(chan ElevLib.NextOrder)
	setlights 			 := make(chan bool)
	setLightsOff := make(chan []int)
	currentfloorupdateFSM	:= make(chan int)
	
	
	// COMMUNICATION BETWEEN EM AND QUEUE
	sendReq2Queue 	     := make(chan ElevLib.NewReqFSM)
	receiptFromQueue     := make(chan int)
	updateCurrentFloor := make(chan int)
	setLightsOn := make(chan []int)
	deleteOrderChan := make(chan ElevLib.NextOrder)


	//COMMUNICATION BETWEEN Queue AND NETWORK
	//newExternalOrderChan := make(chan ElevLib.MyOrder)
	//newInfoChan := make(chan ElevLib.MyInfo)
	//externalOrderChan := make(chan ElevLib.MyOrder)


	//FoR QUEUE AND DRIVER
	InternalOrderChan := make(chan ElevLib.MyOrder)
	ExternalOrderChan := make(chan ElevLib.MyOrder)


	go Driver.ReadSensors(sensorchan)
	//rcvCurrentFloorQueue := make(chan chan int)
	current_floor,_ = Driver.Elev_init(sensorchan)
	fmt.Println("Elev_init Done: current_floor = ", current_floor, " and direction = ", direction)
	
	// STARTUP PHASE, GO-ROUTINES
	//go Driver.ReadSensors(sensorchan)
	//go Network.SendAliveMessageUDP()
	//go Network.ReadAliveMessageUDP()
	go Driver.ReadElevPanel(InternalOrderChan)
	go Driver.ReadFloorPanel(ExternalOrderChan)
	go Queue.Queue_manager(sendReq2Queue, receiptFromQueue, localIp, setLightsOn, updateCurrentFloor, InternalOrderChan, ExternalOrderChan, deleteOrderChan)

	
	//fmt.Println("Elev_init Done: current_floor = ", current_floor, " and direction = ", direction)
	


	go Driver.FSM(rcvNewReqFromFSMChan, orderHandledChan, setLightsOff, setlights, currentfloorupdateFSM)
	go Driver.SetLights(setLightsOn, setLightsOff)
	time.Sleep(10*time.Millisecond)

	//go Network.Network(newInfoChan, externalOrderChan ,newExternalOrderChan)
	fmt.Println("MAIN:", "GOING IN FOOR LOOP")
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
				fmt.Println("MAIN: Direction : ", direction)
				
				
				
			case floor := <-orderHandledChan:
				deleteOrderChan <- floor

				//receipt := <- receiptFromQueue  // Trenger egentlig ikke å ta imot
				//fmt.Println("Order on floor ", receipt, " in direction ", direction, " was deleted")

				setlights <- false

			//
			//case newExtOrd := <-newExternalOrderChan: // FORELØPIG BARE FOR Å TESTE 1 HEIS
			//	newExtOrd.Ip = localIp
			//	ExternalOrderChan <- newExtOrd
			
			case current_floor = <-sensorchan:
				fmt.Println("CurrentFloor is: ", current_floor)
				updateCurrentFloor <- current_floor
				select{
					case currentfloorupdateFSM <- current_floor:
					case <-time.After(1*time.Second):
				}
				
		}
		time.Sleep(time.Millisecond)
	}
	*/
}
