package Queue



// 1. Hvordan representere køen? Map?
// 2. Hva med retningsbit?
// 3. Hvor bør variabler ligge? Queue, Driver?
// 4. Hva gjør vi med bestillinger når en heis dør?
//  --> interne kan slettes og externe må ha backup?

// 5. Vi må ha en måte å behandle bestillinger vi ikke kjører forbi
//  --> f.eks: Vi kan ikke bare kjøre kjøre fra 1. etasje til 4. etasje for å forsikre oss om at
//     	       vi får med alle bestillinger



// Variabler

const (
	N_BUTTONS int = 3
	N_FLOORS int = 4
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





// OVERSIKT OVER ALLE BESTILLINGER (MASTER)
// 1. Den må kanskje oppdateres ettersom om en heis kobler seg på nettverket

var ElevatorInfo map[string]map[int]bool  //string-key kan være IPadresse
// allOrders = map[IPadresse]map[floor]struct{ direction(up/down), order(true/false) }




var internalOrders []int
// Shared variable??



var externalOrders [2][N_FLOORS]string 

var unprocessedExtOrders []MyOrder  // Bare midlertidig for å teste Driver





// Ikke ferdige funksjoner

func DeleteInternalOrder(floor int) {
	internalOrder[floor] = false
}

func SetInternalOrders(floor int){
	for i := 0; i < len(internalOrders);i++{
		if floor < internalOrders[i] && dir == 1 {
			internalOrders = insert(internalOrders, floor, i)
		}
		if floor > internalOrders[i] && dir == -1 {
			internalOrders = insert(internalOrders, floor, i+1)
		}
	}
}

func insert (orders []int ,floor, i int) ([]int) {
	// Kanskje vi må passe på størrelsen til orders slik at vi vet at i finnes i orders??
	tmpSlice = orders[:i]
	tmpSlice = append(tmpSlice, floor)
	return append(tmpSlice, orders[i:]...)
}

func SetExternalOrders(button, fl int) {
	extOrd := MyOrder{
		buttonType:	button
		floor:		fl
	}

	unprocessedExtOrders = append(externalOrders, extOrd)
}

func CheckFloor(floor int) (bool){ // Kan brukes til å plukke opp bestillinger på veien
	return internalOrders[floor]
}

func Delete_all_internalOrders(){
	for i,_:= range internalOrders{
		DeleteInternalOrder(i)
	}
}


// Midlertidige testfunksjoner

func GetUnprExtOrd() []MyOrder {
	return unprocessedExtOrders
}

func GetInternalOrders() []int {
	return internalOrders
}



// Ferdige funksjoner
func initexternalOrders() {
	for i:=0;i<2;i++{
		for j:=0;j<N_FLOORS;j++{
			externalOrders[i][j] == ""		
		}		
	}
}




