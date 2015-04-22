package main

import (

	"fmt"
	"time"

)

func main() {
	
	boolVar := false

	if !boolVar {
		fmt.Println("Sleeping")
		time.Sleep(100*time.Millisecond)
		fmt.Println("Test ok!")
	}
}
