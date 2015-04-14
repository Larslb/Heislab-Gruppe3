package Network
import(
	"fmt"
	"Queue"
	"time"
	"encoding/json"
)


// 1. Hva slags informasjon trenger vi å sende?
// 2. En melding for bestilling og en melding for enkle string-meldinger? (eks: "Jeg er Master",
//    "Mottatt"... etc)


const myIPadress string{""}

type NetworkOrder struct{
	Message string
	IPaddr int
	Order MyOrder 
}



func BroadcastOrder(){

}

func ListenUDP(){

}

func sendMessageUDP( msg NetworkOrder ){ // HVA SLAGS TYPE MELDING SKAL VI TA INN?

	baddr,err := net.ResolveUDPAddr("udp", "129.241.187.255:20004")
	sendSock, err := net.DialUDP("udp", nil ,baddr) // connection
	
	// INNKAPSULERING AV MELDING ----> JSON

	send_msg := []byte()
	time.Sleep(1*time.Second)
	_,err = sendSock.Write(send_msg)
	

	// ERROR HANDLING
	//fmt.Println(err)
	//if err != nil{
	//	panic(err)
	//}
}

func recvMessageUDP(){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", ":20004")
	receivesock,_ := net.ListenUDP("udp", raddr)
	for  {
		mlen ,_,_ := receivesock.ReadFromUDP(buffer)
		fmt.Println(string(buffer[:mlen]))
	}

	// ERROR HANDLING
}


