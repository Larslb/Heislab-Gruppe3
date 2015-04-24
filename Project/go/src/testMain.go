package main

import (

	"fmt"
	"./../Driver"
	"ElevLib"
	"time"

)

func main() {
	buttonChan := make(chan ElevLib.MyOrder)
	buttonChan2 := make(chan ElevLib.MyOrder)
	go Driver.ReadElevPanel(buttonChan)
	go Driver.ReadFloorPanel(buttonChan2)
	for {
		select{
		case buttonpress:= <-buttonChan:
			fmt.Println(buttonpress.ButtonType, buttonpress.Floor)
		case buttonpressed := <-buttonChan2:
			fmt.Println(buttonpressed.ButtonType, buttonpressed.Floor)

		}
		
	}
	time.Sleep(100*time.Second)
}
