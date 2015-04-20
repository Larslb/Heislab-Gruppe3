package Queue




// 1. Hva gjør vi med bestillinger når en heis dør?
//  --> interne kan slettes og externe må ha backup?


// Variabler og typer

const (
	N_BUTTONS int = 3
	N_FLOORS int = 4
	
	IDLE int = 0
	UP int = 1
	DOWN int = 2
	OPEN_DOOR = 3
)


type MyInfo struct {
	Ip string
	Dir int
	InternalOrders []int
}

type MyOrder struct {
	Ip string
	ButtonType int
	Floor int
	//Value int
}



var ElevatorInfo map[string]map[int]bool  //string-key kan være IPadresse
// allOrders = map[IPadresse]map[floor]struct{ direction(up/down), order(true/false) }




// Ikke ferdige funksjoner

func setInternalOrder(iOrders []int, floor, dir int) ([]int) {
	if dir == 1 {
		// If floor - current position < 0, then append at back
		// if floor - current position = 0, open door
		for i := 0; i < len(iOrders); i++ {
			if floor < iOrders[i]{
				return insert(iOrders, floor, i)
			}
		
		}
	} else if dir == -1 {
		// If current position - floor < 0, then append at back
		// if current position - floor = 0, then open door
		for i := 0; i < len(iOrders); i++{
			if floor > iOrders[i] {
				return insert(iOrders, floor, i)
			}
		}
	}
	return append(iOrders, floor)
}

// Ferdige funksjoner

func insert (orders []int ,floor, i int) ([]int) {  // FUNKER BRA!
	tmp := make([]int, len(orders[:i]), len(orders)+1)
	copy(tmp, orders[:i])
	tmp = append(tmp, floor)
	tmp = append(tmp, internalOrders[i:]...)
	return tmpSlice
}

func setExternalOrder(eOrders [2][N_FLOORS]string, order MyOrder) ([2][N_FLOORS]string) {
	eOrders[order.ButtonType][order.Floor] = order.Ip
	return eOrders	
}


func queue_manager(intrOrd chan int, extrOrd chan myOrder, dirOrNF chan int, deleteOrdFloor chan int, reqInfo chan myInfo, msg chan string){

	tmpDir := 0 // brukes av setInternalOrders, men er lokalt lagret i FSM

	internalOrders := []int{}
	externalOrders := [2][N_FLOORS]string{}   // [0][...] er opp bestillinger og [1][...] er ned bestillinger
	
	for {
		select{
		case order := <- intrOrd:
			// Må passe på at tmpDir er oppdatert
			internalOrders = setInternalOrder(internalOrders, order, tmpDir)
			
		case order := <- extrOrd:
			externalOrders = setExternalOrder(externalOrders, order)
			
		case tmpDir = <- dirOrNF: // Bedre navnsetting på channel?
		
			// Brukes på følgende måte: FSM ber om neste bestilling ved å sende en oppdatert retning dir som lagres lokalt i QM. 
			// QM responderer på samme kanal ved å sende nextOrder til FSM.
			dirOrNF <- nextOrder(internalOrders, externalOrders, dir)
			
			
		case deleteOrder := <- deleteOrdFloor:
			// BRUKES NÅR FSM HAR HÅNDTERT EN BESTILLING OG ER KLAR TIL Å TA NESTE BESTILLING
			// OBS! Må passe på at tmpDir er oppdatert!
			
			internalOrders = internalOrders[1:]
			
			if tmpDir == 1 {
				externalOrders[BUTTON_CALL_UP][deleteOrder] = ""
			} else if tmpDir == -1 {
				externalOrders[BUTTON_CALL_DOWN][deleteOrder] = ""
			}
			
			// else ?
			
		case <-reqInfo:
			// Brukes til: Når "Network" har behov for å vite info om heisens interne ordre og kjøreretning.
			// Dette gjøres ved at "Network" sender en request (tom MyInfo), og QM responderer ved å sende infoen til heisen tilbake
			
			reqInfo <- MyInfo{
					Ip: myIP,
					Dir: direction,
					InternalOrders: internalOrders,
					}
		}
	}
}

func Fsm(chan dirOrNextFloor int, chan deleteOrderOnFloor) { // skal være i drivermodulen

	current_floor := -1 // initielt
	direction     := 0  // initielt
	
	next_floor    := -1 // initielt ingen bestillinger
	
	STATE := IDLE
	
	for {
		switch state {
			
			// HVILKE TILSTANDER TRENGER VI??
			
			case IDLE:
				// 1. request next floor
				dirOrNextFloor <- direction
				next_floor = <-dirOrNextFloor
				
				if current_floor
				if next_floor == -1 { //Ingen bestillinger
					//time.Sleep(???) for å ikke overbelaste QM med requests
				}
				
				
				
				
				
			case UP:
			case DOWN:
			case OPEN_DOOR:
		}
	
	}
}
