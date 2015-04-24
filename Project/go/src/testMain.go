package main

import (

	"fmt"
	"Driver"
	//"time"

)

func readbuttons(buttonChan chan int){
	for {
			for i:=0;i<4;i++{
			elev_get_button_signal(ElevLib.2,i, buttonChan)
		}
	}

}

func main() {
	buttonChan := make(chan int)
	go readbuttons(buttonChan)
	fmt.Println(<-buttonChan)
}
