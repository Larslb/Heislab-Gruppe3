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
	localHost string = "127.241.187.255"
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



func Init(){
	localIP,localconn = GetLocalIP()
	adresses[localIP] = time.Now()
	infomap[localconn] = nil
}


/////////////////////////////////////////////////////////////////////////////
/////////////////////////Logiske funksjoner//////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
func GetLocalIP() (string,*net.UDPAddr){
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

func Master(writeToSocketMaprecvInfo chan queue.MyInfo, recvOrder chan queue.MyOrder, extOrder chan string , read chan bool, PanelOrder chan queue.MyOrder) {

	var orders := []queue.MyOrder{}
	sendorder := make(chan queue.MyOrder)
	time.Sleep(100*time.Millisecond)
	writeToSocket <- false
	read <- true
	go ReadALL(read, recvInfo, recvOrder)
	for {

		if (!boolvar) {
			fmt.Println("Going slavemode")
			//trenge å sende orders til en kanal som kan mottas når den blir slave slik at den kan sende ut til mastersocket sine tidligere ordre.
			return
		}

		select{
			case NewInfo <-recvInfo:
				//OPPDATERE INFOMAP MED INFOEN MOTTATT PÅ SOCKET
				infomap[NewInfo.IPadresse] = NewInfo
			case NewOrder <-recvOrder:
				read<-false
				//SENDE UT NY ORDRE TIL ALLE
				orders = append(orders, NewOrder)
				handledorder = orderhandler(orders)

				for _,socket := range socketmap{
					sendorder(socket, handledorder)
				}
				read<-true

			case order <- PanelOrder:

				handledorder = orderhandler(order)
				broadCastOrder(handledorder)
		}
		


		/*
		if (!boolvar) {
			fmt.Println("Going Slavemode")
			return //return stopper threaden
		}

		recInfo <- queue.MyInfo{}
		infomap[localconn] = <- recvInfo
		go RecieveOrders(recvOrder, orders)
		sendorder := make(chan queue.MyOrder)
		//Trenger en kanal for å legge til egne bestillinger fra internkø

		//LESER INN INFO FRA ALLE HEISENE OG OPPDATERER INFOMAPET
		for _,connection := range socketmap{
			info := ReadInfo(connection)
			infomap[connection.RemoteAddr()] = info
		}

		//LESER INN ORDRE FRA ALLE HEISENE OG LEGGER DET I EN EGEN LISTE
		for socket,_ := range socketmap{
			orders[socket]<-recvOrder
		}

		// må ta in en kanal
		//costfunksjon
		
		for _,socket := range socketmap {
			tempOrder <- sendorder
			json
			socket.Write([]byte())
			//sende orderen til alle, med info om hvem som skal ta seg av ordren
		}
		*/


	}


	//leser inn info fra heiser
	//leser inn ordre fra heiser.. Kanskje en kanal som står å venter på ordre fra riktig port?
	//kjører costfunksjon
	//Broadcaster orderen til alle, med rett IP 

}

func orderhandler(order queue.MyOrder) {


	for key,value := range infomap{
		if value.internalOrders == nil {
			
		}
		if order.ButtonType == elev.BUTTON_CALL_UP {
			
		}
		if  {
			
		}
	}
	
}

func ReadALL(read chan bool, chan recvInfo queue.MyInfo, chan recvOrder queue.MyInfo) {
	for _,connection := range socketmap{
		if <-read {
			buffer := make([]byte,1024)
			msglen ,_:= connection.Read(buffer)
			var temp queue.MyElev
			json.Unmarshal(buffer[:msglen], &temp)
			if tempInfo.MessageType == "INFO" {
				recvInfo <-temp.Info
			}else if tempInfo.MessageType == "ORDER" {
				recvOrder <-temp.Order
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
func Slave(chan recvInfo queue.MyOrder, extOrder chan queue.MyOrder, Panelorder chan queue.MyOrder, readyToRecv chan bool) {
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
			case NewOrder <- recievechannel:
				if (NewOrder.IPadresse == localIP) {
					extOrder <-NewOrder
				}
			case NewPanelOrder <- Panelorder:
				sendtoMastersocket()
			}
		}
		
	}

		//sende Ordre, og motta ordre
}

	//sender inn alle bestillinger den mottar fra panel til master
	//lytter på port for å motta en ordre fra master
	//setter inn bestilling i queue med hvilken IP som skal di


func TCPAccept(writeToSocket chan bool) {
	listenAddr, error := net.ResolveTCPAddr("tcp4", localIP+tcpPort)
	if error != nil {
		fmt.Println(error)
	}
	listener, error := net.ListenTCP("tcp4",listenAddr)
	if error != nil {
		fmt.Println(error)
	}
	for{
		listener.SetDeadline(time.Now().Add(time.Millisecond*100))
		remoteConn, error := listener.AcceptTCP()
		if (error == nil && <-writeToSocket){
			socketmap[strings.Split(remoteConn.RemoteAddr().String(), ":")[0]] = remoteConn
		}
		if (!boolvar) {
			return 
		}
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


func RecieveOrders(chan recvorder queue.MyOrder) {
	for {
		for socket,_ := range socketmap{
			recvorder <- ReadOrder(socket)

		}
	}
}
func SendInfoMessageUDP(info chan Queue.MyInfo, socket *net.TCPConn ){
	time.Sleep(10*time.Millisecond)
	for {
		tempInfo := <-info
		buf,_ := json.Marshal(tempInfo)
		_,err = socket.Write([]byte(buf))
		if err != nil{
			panic(err)
		}
		time.Sleep(10*time.Millisecond)
	}
	conn.close()
}


func SendOrder(socket *net.TCPConn, order queue.MyOrder) {
	buffer := json.Marshal(order)
	_,err = socket.Write(buf)
	if err != nil{
		panic(err)
	}
	
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
	broadcastOrderSock.close()
}

func recvInfoUDP(readyToRecv chan bool, recvIP string){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", recvIP + COMPORT)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for {
		msglen , _,_ := recieveSock.ReadFromUDP(buffer)
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

func ReadInfo(socket *net.TCPConn)queue.MyInfo{
	buffer := make([]byte,1024)
	msglen ,_:= socket.Read(buffer)
	var tempInfo queue.MyInfo
	json.Unmarshal(buffer[:msglen], &tempInfo)
	return tempInfo
}

func ReadOrder(socket *net.TCPConn)queue.MyOrder{
	buffer := make([]byte,1024)
	msglen ,_:= socket.Read(buffer)
	var tempOrder queue.MyOrder
	json.Unmarshal(buffer[:msglen], &tempOrder)
	return tempOrder
}


func recvOrderUDP(orderchannel chan Queue.MyOrder, readyToRecvOrder chan bool){
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", localHost + BRPORT)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for {
		msglen ,_,_ := recieveSock.ReadFromUDP(buffer)
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


func EventManager_NetworkStuff(writeToSocketmap chan bool) {  

	time.Sleep(1*time.Second)
	master = SolvMaster()
	if (!master) {
		boolvar = true		
	}
	for {
		if (master) {

			if (!boolvar) {
			fmt.Println("Im Master")
			boolvar = true
			writeToSocketmap <- true
			go AcceptTCP()
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



