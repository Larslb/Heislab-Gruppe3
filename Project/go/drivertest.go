package main

import (
	"Driver"
	//"Nettwork"
	"fmt"
	//"time"
)



func main(){
	
	
	buttonFloorChan := make(chan int, 1)
	buttonElevChan := make(chan int, 1)
	stopChan := make(chan int, 1)
	sensorChan := make(chan int, 1)


	go Driver.ElevPanelThread(buttonElevChan)
	go Driver.SensorThread(sensorChan)
	go Driver.StopThread(stopChan)
	go Driver.FloorPanelThread(buttonFloorChan)	
	
	
	Driver.Elev_init()

	

	bestilling := <- buttonElevChan
	fmt.Println(bestilling)
	//for {
	//	
	//	select{
		// Hvordan ta imot verdier sendt over channels i case-blokkene??
	//	case <- buttonElevChan:
	//	case 2:
	//	case 3:
	//	case 4:
	//	}
	//}
}
