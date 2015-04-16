package Driver

import(
	"Queue"
	"Network"
	"os"
)

const(
	UP
	DOWN
	IDLE
	DOOR
)

func DriverStateMachine() {
	// Hva vi kan få bruk for:
	// Current floor
	// Next floor
	// Dir 

	



	for  {
		switch STATE {
		case IDLE:
			// 1. lese internalOrders og/eller external Orders for å finne next floor?
			//    --> kan være et problem å håndtere external orders når vi kjører forbi etasjene
			//    1.1 Lage en getNextFloor funksjon i Queue
			// 2. Sette nextFloor og sammenlikne med current floor for velge direction 
			// (og state)

		case UP:
			// 1. Sette hastighet opp
			// 2. Sjekke når vi har kommet til etasjen
			// 3. Sette STATE = stop/door hvis floor reached
		case DOWN:
			// 1. Sette hastighet ned
			// 2. Sjekke når vi har kommet til etasjen
			// 3. Sette STATE = stop/door hvis floor reached
		case DOOR:
			// 1. Stoppe heisen
			// 2. Åpne og lukke døren
			// 3. Sette STATE = IDLE



		}
	}
}


func Driveup(){
	elev

}