package Network
import(
	"fmt"
	"os"
	"net"
	"Queue"
	"time"
	"encoding/json"
)


// 1. Hva slags informasjon trenger vi å sende?
// 2. En melding for bestilling og en melding for enkle string-meldinger? (eks: "Jeg er Master",
//    "Mottatt"... etc)

//


const (
	N_FLOORS int = 4
	N_BUTTONS int = 3
	localHost string = "127.241.187.255"
	BRPORT string = "25556"
	COMPORT string = "25557"
	) 

var localIP string = "0"
var lowestIP string = "0"
var masterIP string = "0"
var infomap = make(map[*net.UDPAddr]Queue.ElevatorInfo)
var addresses = make(map[string]time.Time)



func Init(){
	localIP = GetLocalIP()
	adresses[localIP] = time.Now()

}

func GetLocalIP() (string){
   addr, _ := net.ResolveTCPAddr("tcp4", "google.com:80")
   conn, _ := net.DialTCP("tcp4", nil, addr)
   return strings.Split(conn.LocalAddr().String(), ":")[0]
}



func Master() {

}

func Slave() {
	
}

func SolvDeadMaster(){

	localIP = lowestIP

}
func SendAliveMessageUDP(){
	broadcastAliveaddr,_ := net.ResolveUDPAddr("udp", localHost + BRPORT)
	broadcastAliveSock,_ := net.DialUDP("udp", nil, broadcastAliveaddr)
	time.Sleep(10*time.Millisecond)
	for {
		_,err = broadcastAliveSock.Write(localIP)
		if err != nil{
			panic(err)
		}
		time.Sleep(10*time.Millisecond)
	}
}

func SendInfoMessageUDP(info chan Queue.ElevatorInfo, IPMaster chan string){
	addr,_ := net.ResolveUDPAddr("udp", <-IPMaster + COMPORT)
	conn,err := net.DialUDP("udp", nil, addr)
	time.Sleep(10*time.Millisecond)
	for {
		tempInfo := <-info
		buf,_ := json.Marshal(tempInfo)
		_,err = conn.Write(buf)
		if err != nil{
			panic(err)
		}
		time.Sleep(10*time.Millisecond)
	}
	conn.close()
}

func SendOrderMessageUDP(order chan Queue.MyOrder){
	broadcastOrderaddr,_ := net.ResolveUDPAddr("udp",  )
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

func recvInfoUDP(readyToRecv chan bool){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", ":25557")
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for {
		msglen , recvIP,_ := recieveSock.ReadFromUDP(buffer)
		var tempInfo Queue.ElevatorInfo
		json.Unmarshal(buffer[:mlen], &tempInfo)
		if <-readyToRecv{ 
			infomap[recvIP] = tempInfo
			readyToRecv <- false
		}
		time.Sleep(10*time.Millisecond)
	}
	recieveSock.close()
}

func recvOrderUDP(orderchannel chan Queue.MyOrder, readyToRecvOrder chan bool){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", localHost + BRPORT)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for {
		msglen ,masterIP,_ := recieveSock.ReadFromUDP(buffer)
		var tempOrder MyOrder
		json.Unmarshal(buffer[:mlen], &tempOrder)
		if <- readyToRecv{
			orderchannel <- tempOrder
			readyToRecv <- false
		}
		time.Sleep(10*time.Millisecond)
	}
	recieveSock.close()
}

func Listen() {
}

func ReadAliveMessageUDP(boolvar bool){
	addr,err := net.ResolveUDPAddr("udp", localHost+BRPORT)
	if err != nil {
		fmt.Println(err)
		return
	} 
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	} 
	buffer := make([]byte(1024))
	for {
		conn.ReadFromUDP(buffer)
		s := string(buffer[0:15]) //slipper nil i inlesningen
		adresses[string(s)] = time.Now()
		if s!= nil {
			boolvar=true
			for key, value := adresses{
				if ((time.Now().Sub(value) > 100*time.Millisecond) && (key != localIP)){
					delete(adresses,key)

				}
			}
		}
		else{
			boolvar=false
		}
		time.Sleep(10*time.Millisecond)
	}
	conn.close()
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



