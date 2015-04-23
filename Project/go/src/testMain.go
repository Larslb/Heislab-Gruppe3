package main

import (

	"fmt"
	//"time"

)

func rcv(c1 chan chan chan chan int) {

	c2 := <- c1
	


	c2 <- 1
}

func main() {
	
	c3 := make(chan int)
	c3 := make(chan int)
	c3 := make(chan int)
	c3 := make(chan int)
	boolVar := make(chan chan chan chan int)


	go rcv(boolVar)
	boolVar <- c3

	b := <- c3

	fmt.Println(b)
}
