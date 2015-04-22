package main

import (
	"Queue"
	"Driver"
	"Network"
	"ElevLib"
	"time"
)



func main() {
	
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
}
