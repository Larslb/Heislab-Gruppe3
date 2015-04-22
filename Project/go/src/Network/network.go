package Network
import(
	"fmt"
	"os"
	"net"
	"Queue"
	"time"
	"encoding/json"
	"strconv"
)


// 1. Hva slags informasjon trenger vi å sende?
// 2. En melding for bestilling og en melding for enkle string-meldinger? (eks: "Jeg er Master",
//    "Mottatt"... etc)

//


const (
	N_FLOORS int = 4
	N_BUTTONS int = 3
	localHost string = "129.241.187.255"
	BRALIVE string = "25556"
	BRORDER string = "25555"
	tcpPort string = "25557"
	) 

var localIP string = "0"
var lowestIP string = "0"
var infomap = make(map[string]queue.MyInfo)
var socketmap = make(map[string]*net.TCPConn)
var addresses = make(map[string]time.Time)
var master bool = false
var boolvar bool = false



func Init(localipch chan string){
	localIP,localconn = GetLocalIP()
	adresses[localIP] = time.Now()
	infomap[localconn] = nil
}


/////////////////////////////////////////////////////////////////////////////
/////////////////////////Logiske funksjoner//////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
func GetLocalIP() (string,*net.TCPConn){
   addr, _ := net.ResolveTCPAddr("tcp4", "google.com:80")
   conn, _ := net.DialTCP("tcp4", nil, addr)
   return strings.Split(conn.LocalAddr().String(), ":")[0],conn
}


func SolvMaster() bool{
	//returner true hvis jeg er master
	//brukes til å sjekke hvem som er master basert på lavest IP
	//returnerer false hvis jeg ikke har lavest IP

	localIP = lowestIP

	for key,_ := range addresses{
		IP1,_ := strconv.Atoi(strings.SplitAfterN(key,".",-1)[3])
		IP2,_ := strconv.Atoi(strings.SplitAfterN(lowestIP,".",-1)[3])

		if (IP1 < IP2) {
			lowestIP = key
		}
	}

	if lowestIP =localIP{
		return true
	}else{
		return false
	}
}


func OrderNotInList([]orders queue.MyOrder, neworder queue.MyOrder) (bool) {
	for i := 0; i < len(orders); i++ {
		if (neworder.ButtonType == orders[i].ButtonType) && neworder.Floor == orders[i].Floor {
			return false
		}else{
			return true
		}
	}
}


///////////////////////UDP funksjoner/////////////////////////////////
//////////////////////////////////////////////////////////////////////

func SendAliveMessageUDP(){
	broadcastAliveaddr,_ := net.ResolveUDPAddr("udp", localHost + BRALIVE)
	broadcastAliveSock,_ := net.DialUDP("udp", nil, broadcastAliveaddr)
	time.Sleep(10*time.Millisecond)
	for {
		_,err = broadcastAliveSock.Write([]byte(localIP))
		if err != nil{
			return
		}
		time.Sleep(10*time.Millisecond)
	}
	broadcastAliveSock.close()
}



func ReadAliveMessageUDP(readvar bool){
	addr,err := net.ResolveUDPAddr("udp", localHost+BRALIVE)
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
			readvar=true
			for key, value := range adresses{
				if ((time.Now().Sub(value) > 100*time.Millisecond) && (key != localIP)){
					delete(adresses,key)
				}
			}
		}
		else{
			readvar=false
		}
		time.Sleep(10*time.Millisecond)
	}
	conn.close()
}


func broadCastOrder(order queue.MyOrder) {
	broadcastOrderaddr,_ := net.ResolveUDPAddr("udp", localHost + BRORDER)
	broadcastOrderSock,_ := net.DialUDP("udp", nil, broadcastOrderaddr)
	time.Sleep(10*time.Millisecond)
	for {
		tmpOrder := <- order
		buf,_ := json.Marshal(tmpOrder)
		_,err = broadcastOrderSock.Write(buf)
		if err != nil{
			panic(err)
			}
		}
	}
	broadcastOrderSock.close()
}

func RecieveOrders(readyToRecv chan bool, orderchannel queue.MyOrder) {
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", localHost + BRORDER)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for { 
		if (<-readyToRecv) {
			msglen ,_,_ := recieveSock.ReadFromUDP(buffer)
			var tempOrder MyOrder
			json.Unmarshal(buffer[:mlen], &tempOrder)
			orderchannel <- tempOrder
			readyToRecv <- false
		}
		time.Sleep(10*time.Millisecond)
	}
}

//////////////////////////TCP funksjoner/////////////////////////////////
/////////////////////////////////////////////////////////////////////////

func Master(writeToSocketMap chan bool, recvInfo chan queue.MyInfo, recvOrder chan queue.MyOrder, extOrder chan string , read chan bool, PanelOrder chan queue.MyOrder) {

	var orders := []queue.MyOrder{}
	sendorder := make(chan queue.MyOrder)
	time.Sleep(100*time.Millisecond)
	writeToSocketmap <- false
	go ReadALL(writeToSocketmap, recvInfo, recvOrder)
	for {

		if (!boolvar) {
			fmt.Println("Going slavemode")
			//trenge å sende orders til en kanal som kan mottas når den blir slave slik at den kan sende ut til mastersocket sine tidligere ordre.
			return
		}

		select{
			case NewInfo := <-recvInfo:
				//OPPDATERE INFOMAP MED INFOEN MOTTATT PÅ SOCKET
				infomap[NewInfo.IPadresse] = NewInfo
			case NewOrder := <-recvOrder:
				//SENDE UT NY ORDRE TIL ALLE
				orders = append(orders, NewOrder)
				handledorder = orderhandler(orders)

				broadCastOrder(handledorder)

			case Ownorder := <- PanelOrder:

				handledorder = orderhandler(Ownorder)

				if handledorder.IPadresse == localIP {
					extOrder <- handledorder
				}
				broadCastOrder(handledorder)
		}
	}
}

func orderhandler(order queue.MyOrder)(string) {

	var besteheis := make(map[string]int){}

	for key,value := range infomap{
		if  value.CurrentFloor == order.Floor{
			besteheis[
			return besteheis
		}

		if  {
			
		}
		/*
		else if value.internalOrders == nil {
			
		}
		else if order.ButtonType == elev.BUTTON_CALL_UP {

		}
		else if  order.ButtonType == elev.BUTTON_CALL_DOWN{

		}
		if value.internalOrders[0] == order.Floor && value.dir == {
			
		}
		)
		for i := 0; i < len(value.internalOrders); i++ {
			if value.internalOrders[i] == order.Floor && value.dir 
				
			}
		}*/
	}
}

func ReadALL(writing chan bool, chan recvInfo queue.MyInfo, chan recvOrder queue.MyInfo) {
	for  {
		<-writing
		for _,connection := range socketmap{
			buffer := make([]byte,1024)
			msglen ,_:= connection.Read(buffer)
			var temp Queue.MyElev
			json.Unmarshal(buffer[:msglen], &temp)
			if temp.MessageType == "INFO" {
				recvInfo <-temp.Info
				writing<-1
			}else if temp.MessageType == "ORDER" {
				recvOrder <-temp.Order
				writing<-1
			}else{
				continue
			}
		}
	}
}


func ReadOrders(chan recvOrder queue.MyOrder){
	for _,connection := range socketmap{
		buffer := make([]byte,1024)
		msglen ,_:= connection.Read(buffer)
		var tempOrder
	}
}
func Slave(chan sendInfo queue.MyInfo, extOrder chan queue.MyOrder, Panelorder chan queue.MyOrder, readyToRecv chan bool) {
	var masterSocket *net.TCPConn 
	var connected bool = false
	for(connected==false){
		masterSocket,connected = ConnectToIP(lowestIP)
	}
	var recievechannel queue.MyOrder
	readyToRecv<-true
	go RecieveOrders(readyToRecv, recievechannel)
	for {
		if (boolvar) {
			fmt.Println("Going from slave to master")
		}

		for {

			select{
			case NewOrder := <- recievechannel:
				if (NewOrder.IPadresse == localIP) {
					extOrder <-NewOrder
				}
			case NewPanelOrder := <- Panelorder:
				//newOrder needs to be sent to mastersocket
			case InfoUpdate := <- sendInfo:
				//send Info to mastersocket
			}
		}
		
	}

		//sende Ordre, og motta ordre
}

	//sender inn alle bestillinger den mottar fra panel til master
	//lytter på port for å motta en ordre fra master
	//setter inn bestilling i queue med hvilken IP som skal di


func TCPAccept(writeToSocket chan int) {
	listenAddr, error := net.ResolveTCPAddr("tcp4", localIP+tcpPort)
	if error != nil {
		fmt.Println(error)
	}
	listener, error := net.ListenTCP("tcp4",listenAddr)
	if error != nil {
		fmt.Println(error)
	}
	for{
		<-writeToSocket
		listener.SetDeadline(time.Now().Add(time.Millisecond*100))
		remoteConn, error := listener.AcceptTCP()
		if (error == nil){
			socketmap[strings.Split(remoteConn.RemoteAddr().String(), ":")[0]] = remoteConn
		}
		if (!boolvar) {
			return 
		}
		writeToSocket<-1
	}
}

func ConnectToIP(IP string)(*net.TCPConn, bool){
	remoteAddr,error := net.ResolveTCPAddr("tcp4", IP + tcpPort)
	if error != nil{
		fmt.Println(error)
		panic(error)
	}
	conn, error := net.DialTCP("tcp4", nil, remoteAddr) 
    if(error==nil){
    	return conn, true
    }else{
    	return conn, false
    }

}

///////////////////////////////diverse funksjoner/////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////


func EventManager_NetworkStuff() {  

	time.Sleep(1*time.Second)
	master = SolvMaster()
	writeToSocketmap := make(chan int,1)
	if (!master) {
		boolvar = true		
	}
	for {
		if (master) {

			if (!boolvar) {
			fmt.Println("Im Master")
			boolvar = true
			go AcceptTCP(writeToSocketmap)
			go master(writeToSocketmap)
			}

		master = SolvMaster()
		}else{
			if (boolvar) {
				fmt.Println("Im a Slave biatch")
				boolvar = false
				go slave()
			}
		master = SolvMaster()
		}
	}


}



