
package main

import (
    . "fmt"    
    "runtime"
    //"time"
)

var i int

func thread1(c1 chan int, cdone chan int) {
	
    	for j := 0; j < 1000000; j++ {
		<-c1
		i=i-1	
		c1<-1
	}
	cdone <-1
}

func thread2(c1 chan int, cdone chan int) {
	for j := 0; j < 1000001; j++ {
		<-c1
		i=i+1
		c1<- 1
	}
	cdone <-1 
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	c1 := make(chan int,1)
	cdone := make(chan int)
    	
	go thread1(c1, cdone)                     
	go thread2(c1, cdone)
	c1 <- 1
    				
    //time.Sleep(100*time.Millisecond)
	<-cdone
	<-cdone
	Println(i)
}

