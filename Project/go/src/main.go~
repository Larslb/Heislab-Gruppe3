package main

import (
	"./Queue"
	"./Driver"
	"./Network"
	"./Notify"
	"./ElevLib"
	"time"
	"fmt"
)



func main() {
	localIp,_ := Network.GetLocalIP()
	Network.Init()

	//errorchannels
	errorDetection := make(chan bool) 

	//orderchannels
	externalOrderChan := make(chan ElevLib.MyOrder)
	internalOrderChan := make(chan ElevLib.MyOrder)


	//NETWORKCHANNELS
	newInfoChan := make(chan ElevLib.MyInfo) 
	newPanelOrderChan := make(chan ElevLib.MyOrder)
	readAndWriteAdresses := make(chan int, 1)
	masterChan := make(chan int)
	slaveChan := make(chan int)

	//Queue - Driver channels
	QM2FSM := make(chan ElevLib.QM2FSMchannels)
	setLightsOn		  := make(chan []int)
	setLightsOff      := make(chan []int)

	//network - Queue connection channels
	orderDel2Master := make(chan ElevLib.MyOrder)
	ordrDeleteFromMaster := make(chan ElevLib.MyOrder)

	//NETWORK INIT!
	go Network.SendAliveMessageUDP()
	go Network.ReadAliveMessageUDP(readAndWriteAdresses)
	readAndWriteAdresses<-1
	time.Sleep(time.Second)
	go Network.SolvMaster(readAndWriteAdresses, masterChan, slaveChan)
	go Network.Network(newInfoChan, externalOrderChan, newPanelOrderChan, masterChan, slaveChan, orderDel2Master, ordrDeleteFromMaster)
	go Driver.ReadElevPanel(internalOrderChan)
	go Driver.ReadFloorPanel(newPanelOrderChan)
	go Driver.SetLights(setLightsOn, setLightsOff)
	go Queue.Queue_Manager(QM2FSM, internalOrderChan, externalOrderChan, setLightsOn, localIp, newInfoChan, orderDel2Master, ordrDeleteFromMaster)
	go Driver.Fsm(QM2FSM, setLightsOff)

	
	for{
		select{

			case err := <-errorDetection: 
				fmt.Println(err)

			case <-time.After(200*time.Millisecond):
				time.Sleep(10*time.Second)  
				
		}
		time.Sleep(10*time.Millisecond)
	}

}
