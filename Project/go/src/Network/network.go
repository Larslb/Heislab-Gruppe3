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

//


const myIPadress string 



func SendAliveMessageUDP(){
	broadcastAliveaddr,_ := net.ResolveUDPAddr("udp", "127.241.187.255:25556")
	broadcastAliveSock,_ := net.DialUDP("udp", nil, broadcastAliveaddr)
	time.Sleep(1*time.Second)
	sendMsg := []byte("I'm Alive!!")
	for {
		buf,_ := json.Marshal(sendMsg)
		_,err = broadcastAliveSock.Write(buf)
		if err != nil{
			panic(err)
		}
		time.Sleep(1*time.Second)
	}
}

func SendInfoMessageUDP(info MyInfo, IPMaster string){
	SendToMasteraddr,_ := net.ResolveUDPAddr("udp", IPMaster + "25557")
	SendToMasterSock,err := net.DialUDP("udp", nil, SendToMasteraddr)
	time.Sleep(1*time.Second)
	for {
		buf,_ := json.Marshal(info)
		_,err = SendToMasterSock.Write(buf)
		if err != nil{
			panic(err)
		}
		time.Sleep(100*time.Millisecond)
	}
}

func SendOrderMessageUDP(order chan MyOrder){
	broadcastOrderaddr,_ := net.ResolveUDPAddr("udp", "127.241.187.255:25555")
	broadcastOrderSock,err := net.DialUDP("udp", nil, broadcastOrderaddr)
	time.Sleep(1*time.Second)
	for {
		tmpOrder := <- order
		buf,_ := json.Marshal(tmpOrder)
		_,err = broadcastOrderSock.Write(buf)
		if err != nil{
			panic(err)
		}
	}
}

func ListenUDP(port string){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", port)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for {
		msglen ,_,_ := recieveSock.ReadFromUDP(buffer)
		var tempInfo Myinfo
		json.Unmarshal(buffer[:mlen], &tempInfo)
		Queue.Storeinfo(tempInfo)
		time.Sleep(1*time.Second)
	}
}

func recvInfoUDP(infochannel chan MyInfo, readyToRecv chan bool
){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", ":25557")
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for {
		msglen ,_,_ := recieveSock.ReadFromUDP(buffer)
		var tempInfo Myinfo
		json.Unmarshal(buffer[:mlen], &tempInfo)
		if <-readyToRecv{ 
			infochannel <- tempInfo
			readyToRecv <- false
		}
		time.Sleep(1*time.Second)
	}
}

func recvOrderUDP(orderchannel chan MyOrder, readyToRecvOrder chan bool){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", "127.241.187.255:25555")
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for {
		msglen ,IPMaster,_ := recieveSock.ReadFromUDP(buffer)
		var tempOrder MyOrder
		json.Unmarshal(buffer[:mlen], &tempOrder)
		if <- readyToRecv{
			orderchannel <- tempOrder
			readyToRecv <- false
		}
		time.Sleep(1*time.Second)
	}
}



func recvAliveMessageUDP(alive string){
	boolvar := true
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", "127.241.187.255:25555")
	recievesock,_ := net.ListenUDP("udp", raddr)
	for(boolvar)  {
		mlen , _,_ := recievesock.ReadFromUDP(buffer)
		alive = string(buffer[:mlen])
	}
}

func EventManager_NetworkStuff() {  //Pseudokode

	initAllConnections
	init all neccessary channels

	go all recv thread ( channels)  
	go all send threads

	for{
		select{
			case msg_info := <-infochannel: 
			//do something
				rcv <- true
			case msg_order := <-orderchannel:
			case ...
		}
	}
}



