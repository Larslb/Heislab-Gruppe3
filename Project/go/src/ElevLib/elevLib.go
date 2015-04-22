package ElevLib

const (
	N_FLOORS int = 4
	BUTTON_CALL_UP int = 0
	BUTTON_CALL_DOWN int = 1
	BUTTON_COMMAND int = 2
	
	WAIT int = 0
	MOVING int = 1
	OPEN_DOOR int = 2
)

type MyInfo struct {
	Ip string
	Dir int
	CurrentFloor int
	InternalOrders []int
}

type MyOrder struct {
	FromIp string
	// ToIp string
	ButtonType int
	Floor int
}

type MyElev struct {
	MessageType string
	Order MyOrder
	Info MyInfo
}
