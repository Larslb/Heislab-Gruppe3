package main

import (
	"Queue"
	"Driver"
	"Network"
)



func main() {
	


	localIpChan := make(chan string)
	newExternalOrderChan := make(chan Elevlib.MyOrder)
	externalOrderChan := make(chan string)
	internalOrderChan  := make(chan int)
	nextFloorChan := make(chan int)
	deleteOrderOnFloorChan := make(chan int)
	newInfoChan := make(chan Elevlib.MyInfo)
	currentFloorAndDirChan := make(chan int)
	setLightsChan := make(chan int)


	go Queue.Queue_manager(internalOrderChan, externalOrderChan, nextFloorChan, deleteOrderOnFloorChan, newInfoChan, currentFloorAndDirChan, setLightsChan, localIpChan)
	
	
	//Network.EventManager_NetworkStuff(newInfoChan, externalOrderChan, newExternalOrderChan)
	//defer close channels	

	Network.Init(localIpChan)
	go Network.SendAliveMessageUDP()
	go Network.ReadAliveMessageUDP()

	go Driver.ReadElevPanel(internalOrderChan)
	go Driver.ReadFloorPanel(newExternalOrderChan)
	go Driver.Fsm(nextFloorChan, deleteOrderOnFloorChan, currentFloorAndDirChan, setLightsChan)


	master := Network.SolvMaster()
	if (!master) {
		boolvar = true		
	}
	for {
		if (master) {

			if (!boolvar) {
			fmt.Println("Im Master")
			boolvar = true
			writeToSocketmap := make(chan int,1)
			go Network.AcceptTCP(writeToSocketmap)
			go Network.Master(writeToSocketmap,newInfoChan, externalOrderChan, newExternalOrderChan)
			}

		master = Network.SolvMaster()
		}else{
			if (boolvar) {
				fmt.Println("Im a Slave biatch")
				boolvar = false
				go Network.Slave(newInfoChan, externalOrderChan, newExternalOrderChan)
			}
		master = Network.SolvMaster()
		}
	}

	// Vi må forsikre oss om at thread Master og Slave avsluttes når man går f.eks. fra Slave til Master
}
