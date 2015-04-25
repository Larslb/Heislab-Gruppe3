package main

import (

	"fmt"
	//"./Driver"
	"./Network"
	"./ElevLib"
	"time"

)

/*func ReeeadSensors(){  // ENDRET TIL EXPORT FUNC
	
	//current_floor := -1
	
	for {
		tmpVal := Driver.Elev_get_floor_sensor_signal()
		fmt.Println("Sensor read: " ,tmpVal)	
		time.Sleep(time.Millisecond)
	}
}
*/


func main() {
	
	Network.Init()
	fmt.Println("INTI!")
	newInfoChan := make(chan ElevLib.MyInfo)
	externalOrderChan := make(chan ElevLib.MyOrder) 
	newExternalOrderChan := make(chan ElevLib.MyOrder)


	go Network.SendAliveMessageUDP()
	go Network.ReadAliveMessageUDP()

	time.Sleep(time.Second)

	go Network.Network(newInfoChan, externalOrderChan, newExternalOrderChan)


	//Driver.Elev_init()
	//go ReeeadSensors()
	//time.Sleep(100*time.Second)

	
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
	*/
	time.Sleep(100*time.Second)
	fmt.Println("Main: Terminating")
}
