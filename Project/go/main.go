package main


import (
	"Driver"
	"Queue"
	"Network"
	
	//"fmt"
	//"time"
)

func main(){ // TIL NÅ ER DETTE BARE EVENT-MANAGER DELEN
	
	// ReadDriver (output)channels
	intrOrd = make(chan Queue.MyOrder) 		// BRUKES SOM CASE TRIGGER, input til queue
	newExtrOrd = make(chan Queue.MyOrder) 	// BRUKES SOM CASE TRIGGER, input til network
	elevSensor  = make(chan int)			// BRUKES SOM CASE TRIGGER, sette currentfloor
	

	// FSM (input)channels
	nextOrd = make(chan int)				// input: requestOrder sendes gjennom her
	currentFloor = make(chan int)			// input: sensor sendes gjennom her


	// Network channels
	extrOrdIn  = make(chan Queue.MyOrder) 	 // input:  newExtrOrd sendes gjennom her
	estrOrdOut = make(chan Queue.MyOrder) 	 // output: brukes som input til Queue handleExtrOrd
	// Hva med synkronisering av myInfo og internalOrders (ligger i queue) og Dir (ligger i driver, fsm) ?

	// Queue channels
	handleIntrOrd = make(chan Queue.MyOrder)
	handleExtrOrd = make(chan Queue.MyOrder) // input: queueOrder sendes inn her
	requestOrder  = make(chan int) 			 // output: brukes som input til FSM nextOrd

	// (kan disse ligge inne i Queue_Manager?)
	queueRead  = make(chan bool) // Hvordan gjør vi synkroniseringen mellom read og write?
	queueWrite = make(chan bool)
	


	go Driver.ReadPanelAndSensor(intrOrd, newExtrOrd, elevSensor)
	go Driver.FSM(nextOrd, currentFloor)
	go Network.Network_Manager(extrOrd)
	go Queue.Queue_Manager(handleIntrOrd, handleExtrOrd, requestOrder)

	for {
		select{ // DISTRIBUERING AV OPPGAVER

		case order := <-intrOrd:
			handleIntrOrd <- order 	// QUEUE LEGGER TIL INTERN ORDRE I KØ

		case order := <-newExtrOrd:
			extrOrdIn <- order 		// NETWORK (MASTER) MÅ HÅNDTERE NY EKSTERN BESTILLING

		case floor := <-elevSensor:
			currentFloor <- floor 	// FSM SETTER CURRENT FLOOR

		case order := <-extrOrderOut:
			handleExtrOrd <- order 	// QUEUE HÅNDTERER EKSTERN BESTILLING SOM MASTER HAR TILDELT

		}
	}
}

func Queue_Manager(intrOrder chan MyOrder)