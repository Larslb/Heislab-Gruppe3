package main

import (

	"fmt"
	"./Driver"
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
func send(sendChan chan int) {

	time.Sleep(3*time.Second)
	sendChan <- 1
}

func main() {
	sensorchan := make(chan int)
	go Driver.ReadSensors(sensorchan)
	Driver.Elev_init(sensorchan)
	
	fmt.Println("TESTMAIN GO")
	buttonChan := make(chan ElevLib.MyOrder)
	buttonChan2 := make(chan ElevLib.MyOrder)
	go Driver.ReadElevPanel(buttonChan)
	go Driver.ReadFloorPanel(buttonChan2)
	breakBool := false
	for {
		select{
		case buttonpress:= <-buttonChan:
			fmt.Println(buttonpress.ButtonType, buttonpress.Floor)
			breakBool = true
		case buttonpressed := <-buttonChan2:
			fmt.Println(buttonpressed.ButtonType, buttonpressed.Floor)

		}
		if breakBool {
			break
		}
		
	}
	fmt.Println("OUT!")
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
	fmt.Println("Main: Terminating")
}
