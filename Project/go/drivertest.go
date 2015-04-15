package main

import (
	"Driver"
	//"Nettwork"
	"fmt"
	//"time"
)



func main(){
	
	if Driver.Elev_init() {
		
		floorChan := make(chan string)
		elevChan  := make(chan string) 
		sensChan := make(chan string)
		
		go Driver.Driver(floorChan, elevChan, sensChan)

		for {
			select{
				case msg := <-floorChan:
					fmt.Println(msg)

				case msg := <-elevChan:
					fmt.Println(msg)

				case msg := <-sensChan:
					fmt.Println(msg)
			}
		}
	}
	
	fmt.Println("Error: Could not initiate driver")
}
