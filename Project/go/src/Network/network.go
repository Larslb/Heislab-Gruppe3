package Network
import(
	"fmt"
	"os"
	"net"
	"Queue"
	"Elevlib"
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
var infomap = make(map[string]Elevlib.MyInfo)
var socketmap = make(map[string]*net.TCPConn)
var addresses = make(map[string]time.Time)
var master bool = false
var boolvar bool = false



func Init(localIpChan chan string){
	localIP,localconn = GetLocalIP()
	localIpChan <- localIP
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

/*

func OrderNotInList([]orders Elevlib.MyOrder, neworder Elevlib.MyOrder) (bool) {
	for i := 0; i < len(orders); i++ {
		if (neworder.ButtonType == orders[i].ButtonType) && neworder.Floor == orders[i].Floor {
			return false
		}else{
			return true
		}
	}
}
*/

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



func ReadAliveMessageUDP(){
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
			for key, value := range adresses{
				if ((time.Now().Sub(value) > 100*time.Millisecond) && (key != localIP)){
					delete(adresses,key)
				}
			}
		}
		time.Sleep(10*time.Millisecond)
	}
	conn.close()
}


func broadCastOrder(order Elevlib.MyOrder) {
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

func RecieveOrders(orderchannel Elevlib.MyOrder) {
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", localHost + BRORDER)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for { 
		msglen ,_,_ := recieveSock.ReadFromUDP(buffer)
		var tempOrder Elevlib.MyOrder
		json.Unmarshal(buffer[:mlen], &tempOrder)
		orderchannel <- tempOrder
		time.Sleep(10*time.Millisecond)
	}
}

//////////////////////////TCP funksjoner/////////////////////////////////
/////////////////////////////////////////////////////////////////////////

func Master(writeToSocketMap chan bool, sendInfo chan Elevlib.MyInfo, extOrder chan string , PanelOrder chan Elevlib.MyOrder) {

	//var orders := []queue.MyOrder{}
	recvInfo := make(chan Elevlib.MyInfo)
	recvOrder := make(chan )
	sendorder := make(chan queue.MyOrder)

	time.Sleep(10*time.Millisecond)
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
				handledorder = orderhandler(order)

				broadCastOrder(handledorder)

			case Ownorder := <- PanelOrder:

				handledorder = orderhandler(Ownorder)

				if handledorder.IPadresse == localIP {
					extOrder <- handledorder
				}
				broadCastOrder(handledorder)
			case UpdateInfo := <- sendInfo:
				infomap[localIP] = UpdateInfo

			//LAGE GO ROUTINE FOR BOOLVARCHECK
		}
	}
}

func orderhandler(order Elevlib.MyOrder)(Elevlib.MyOrder) {

	//var besteheis Elevlib.MyInfo

	order.FromIp = localIP
	return order
	/*
	for key,value := range infomap {
		if value.CurrentFloor == order.Floor {
			besteheis = infomap[key]			
		}
		else{
			if abs(float(value.CurrentFloor) - float(order.Floor)) > 1 {
				besteheis = infomap[key]
			}
			else{
				if abs(value.CurrentFloor -order.Floor) > 2 {
					besteheis = infomap[key]
				}
				else{}
			}
		}

	}
	//for key,value := range infomap{
	//	besteheis = key
	//	return besteheis
		
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

func ReadALL(writing chan bool, recvInfo chan queue.MyInfo, recvOrder chan queue.MyInfo) {
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

func writetoSocket(socket *net.TCPConn, object Elevlib.MyElev)(bool){
	if object.MessageType == "INFO" {
		buffer,_ := json.Marshal(object.Info)
		_,err:= socket.Write(buffer)
		if err != nil {
			fmt.Println("error", err)
			return false
		//errorhandle
		}
		return true
	}
	else if object.MessageType == "ORDER" {
		buffer,_ := json.Marshal(object.Order)
		_,err:= socket.Write(buffer)
		if err != nil {
			fmt.Println("error", err)
			return false
			//errorhandle
		}
		return true
	}
	else{
		return false
	}
}



func Slave(sendInfo chan Elevlib.MyInfo, extOrder chan Elevlib.MyOrder, Panelorder chan Elevlib.MyOrder) {
	var masterSocket *net.TCPConn 
	var connected bool = false
	for(connected==false){
		masterSocket,connected = ConnectToIP(lowestIP)
	}
	var recievechannel Elevlib.MyOrder
	var sendObject Elevlib.MyElev

	go RecieveOrders(recievechannel)
	for {
		if (boolvar) {
			fmt.Println("Going from slave to master")
			return
		}

		for {

			select{
			case NewOrder := <- recievechannel:
				if (NewOrder.IPadresse == localIP) {
					extOrder <- NewOrder
				}
			case NewPanelOrder := <- Panelorder:
				sendObject.MessageType = "ORDER"
				sendObject.Order = NewPanelOrder
				sendObject.Info = nil

				sentorder := writetoSocket(masterSocket, sendObject)
				for !sentorder {
					sentorder = writetoSocket(masterSocket, sendObject)
				}

			case InfoUpdate := <- sendInfo:
				sendObject.MessageType = "INFO"
				sendObject.Order = nil
				sendObject.Info = InfoUpdate

				sentinfo := writetoSocket(masterSocket, sendObject)
				for !sentinfo {
					sentinfo = writetoSocket(masterSocket, sendObject)
				}

				//LAGE CASE FOR IKKE BOOLVAR!!!
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
