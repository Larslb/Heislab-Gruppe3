package Driver

import (
	"time"
)

// SKAL ELEV_INIT KJØRE TIL NÆRMESTE ETASJE?

func ReadElevPanel(buttonChan chan Queue.MyOrder){
	for {
		for i:=0;i<N_FLOORS;i++{
			elev_get_button_signal(BUTTON_COMMAND,i, buttonChan)
		}
	}
}

func ReadFloorPanel(buttonChan chan Queue.MyOrder){
	for{
		for i:=0;i<N_BUTTONS-1;i++{
			for j:=0;j<N_FLOORS;j++{
				elev_get_button_signal(i, j, buttonChan)
			}
		}
	}
}

func readSensors(sensorChan chan int){
	tmpSensorChan := make(chan int, 1) // Nødvendig med buffer?
	for {
		elev_get_floor_sensor_signal(tmpSensorChan)
		tmpVal := <- tmpSensorChan
		if tmpVal != -1 {
			sensorChan <- tmpVal		
		}
	}
}

func Fsm(dirOrNextFloor chan int, deleteOrderOnFloor chan int) { // skal være i drivermodulen

	current_floor := -1 // initielt (må settes av ELEV_INIT)
	direction     := 0  // initielt
	next_floor    := -1 // initielt ingen bestillinger
	
	sensorChan := make(chan int)
	go readSensors(sensorChan)
	
	// ELEV_INIT HER? -> må ta inn sensorChan og sette current_floor
	
	buttonFloorChan := make(chan Queue.MyOrder) // Skal ikke brukes her
	buttonElevChan  := make(chan Queue.MyOrder) // Skal ikke brukes her
	go readElevPanel(buttonElevChan)		  // Skal ikke brukes her
	go readFloorPanel(buttonFloorChan)		  // Skal ikke brukes her
	
	
	STATE := WAIT
	
	for {
		switch STATE {
		
		// BIG PROBLEMS:
		// 1. Vi må kunne oppdatere next_floor samtidig som heisen håndterer den nåværende next_floor
		// 	-> f.eks hvis vi er i 4. etg og next_floor blir satt til 1. etg. På veien får vi en bestilling i 2.etg.
		//	   Da må vi kunne sette next_floor på nytt!
		
		//WAIT:
		//	1. Vi må sjekke om det er en ny bestilling i internalOrders eller externalOrders
			
		// UP or DOWN:
		//	1. Sette motor_direction
		//	2. kontinuerlig sjekke om det er kommet en ny next_floor!
		//		-> kan gjøre det slik at når QM får en ny bestilling, så sjekker den om denne ordren skal bli den nye next floor
		//	3. sjekke når vi har kommet til next_floor
		
			case WAIT:
				
				dirOrNextFloor <- current_floor
				next_floor = <-dirOrNextFloor
				
				if next_floor < current_floor {
				
					direction = -1
					dirOrNextFloor <- direction
					elev_set_motor_direction(direction)
					STATE = MOVING
					
				} else if next_floor > current_floor {
				
					direction = 1
					dirOrNextFloor <- direction
					elev_set_motor_direction(direction)
					STATE = MOVING
					
				} else if next_floor == current_floor {
					
					dirOrNextFloor <- direction // QM venter på direction, så vi må sende her også
					STATE = OPEN_DOOR
					
				} else if next_floor == -1 { //Ingen bestillinger i internalOrders eller externalOrder i direction
					
					direction = 0
					dirOrNextFloor <- direction
					// QM må huske på å returnere -1 dersom internal/external orders er tom for bestillinger
					//time.Sleep(???) for å ikke overbelaste QM med requests
				
				}
				
			//case UP:
				
				// elev_set_motor_direction(direction)
				// KJØR OPP SÅ LENGE VI IKKE HAR KOMMET TIL NEXT FLOOR
				// for {
				//	current_floor <- readSensors
				//	elev_set_floor_indicator(current_floor)
				//	if current_floor == next_floor{
				//		STATE = OPEN_DOOR
				//		break
				//	}
				// }
				
			//case DOWN:
				//elev_set_motor_direction(direction)
				
				// KJØR NED SÅ LENGE VI IKKE HAR KOMMET TIL NEXT FLOOR
				//for {
				//	current_floor <- readSensors
				//	elev_set_floor_indicator(current_floor)
				//	if current_floor == next_floor{
				//		STATE = OPEN_DOOR
				//		break
				//	}
				//}
			
			
			// ALTERNATIVT til UP og DOWN... 
			case MOVING:
			
				current_floor = <- sensorChan
				
				// DENNE KANALEN OG KOMMUNISERINGEN MED QM ER LITT ... mjeee..
				// 1. fsm sender sin current_floor til QM (oppdaterer tmpCurrent_floor)
				// 2. QM sender tilbake next_floor ved å sjekke internalOrders og externalOrders
				// 3. fsm mottar next_floor og sender tilbake direction til QM som oppdaterer tmpDir
				
				
				dirOrNextFloor <- current_floor // Vi spør QM om vi skal plukke opp noen i denne etasjen
				next_floor <- dirOrnextFloor	  // Vi mottar next_floor (enten er det den samme som før, eller en ny)
				dirOrNextFloor <- direction	  // Vi sender tilbake direction
				
				if current_floor == next_floor {
					STATE = DOOR_OPEN
				}
				
				
			case OPEN_DOOR:
				elev_set_motor_direction(0)
				elev_set_door_open_lamp(true)
				// Open door for 3 seconds -> then elev_set_door_open_lamp(false)
				// (Hva gjør vi hvis knapp i samme etasje trykkes inn?) -> dørene skal vel ikke lukkes?

				deleteOrderOnFloor <- current_floor
				ordersDeleted := <-deleteOrderOnFloor // continue when orders are deleted
				
				// SLETTE LYS
				if direction == 1{
					elev_set_button_lamp(BUTTON_CALL_UP, current_floor,0)  // Hva med value = false/true  i stedet for 0/1?
					elev_set_button_lamp(BUTTON_COMMAND, current_floor,0)
				} else {
					elev_set_button_lamp(BUTTON_CALL_DOWN, current_floor, 0)
					elev_set_button_lamp(BUTTON_COMMAND, current_floor,0)
				}
				
				STATE = WAIT				
		}
	}
}
