package Queue



// 1. Hvordan representere køen? Map?
// 2. Hva med retningsbit?
// 3. Hvor bør variabler ligge? Queue, Driver?
// 4. Hva gjør vi med bestillinger når en heis dør?
//  --> interne kan slettes og externe må ha backup?



// Variabler

const (
	N_BUTTONS int = 3
	N_FLOORS int = 4
)


type MyOrder struct{
	ButtonType int
	Floor int
	Value int
	Direction int
}

type externalOrder struct {
	buttonType int
	floor	   int
}



// OVERSIKT OVER ALLE BESTILLINGER (MASTER)
// 1. Den må kanskje oppdateres ettersom om en heis kobler seg på nettverket
var allOrders map[string]map[int]bool  //string-key kan være IPadresse




// Hvordan initialisere internalOrders?
// Kan internalOrders være en referanse til allOrders? ->  &allorders[myIPadress]...
// Hva med å kalle variabelen for myOrders?
var internalOrders map[int]bool




var externalOrders []externalOrder







// Ikke ferdige funksjoner

func initOrders() {
	
}



// Ferdige funksjoner
func DeleteInternalOrder(floor int) {
	internalOrder[floor] = false
}

func SetInternalOrders(floor int){
	internalOrders[floor] = true
}


func SetExternalOrders(button, fl int) {
	extOrd := externalOrder{
		buttonType:	button
		floor:		fl
	}
	externalOrders = append(externalOrders, extOrd)
}

func CheckFloor(floor int) (bool){
	return internalOrders[floor]
}

func Delete_all_internalOrders(){
	for i,_:= range internalOrders{
		DeleteInternalOrder(i)
	}
}
