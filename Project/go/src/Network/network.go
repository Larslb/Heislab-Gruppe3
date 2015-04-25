package Network
import(
	"fmt"
	"net"
	"./../ElevLib"
	"time"
	"encoding/json"
	"strconv"
	"strings"
	
)



// 1. Hva slags informasjon trenger vi å sende?
// 2. En melding for bestilling og en melding for enkle string-meldinger? (eks: "Jeg er Master",
//    "Mottatt"... etc)

//


const (
	N_FLOORS int = 4
	N_BUTTONS int = 3
	localHost string = "129.241.187.255"
	BRALIVE string = ":25556"
	BRORDER string = ":25555"
	tcpPort string = ":25557"
	) 

var localIP string = "0"
var localconn *net.TCPConn 
var lowestIP string = "0"
var infomap = make(map[string]ElevLib.MyInfo)
var socketmap = make(map[string]*net.TCPConn)
var addresses = make(map[string]time.Time)
var master bool = false
var boolvar bool = false



func Init(){
	localIP,localconn = GetLocalIP()
	addresses[localIP] = time.Now()
	socketmap[localIP] = localconn
}


/////////////////////////////////////////////////////////////////////////////
/////////////////////////Logiske funksjoner//////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
func GetLocalIP() (string,*net.TCPConn){
   addr, _ := net.ResolveTCPAddr("tcp4", "google.com:80")
   conn, _ := net.DialTCP("tcp4", nil, addr)
   return strings.Split(conn.LocalAddr().String(), ":")[0],conn
}


func SolvMaster(read chan int, masterchan chan int , slavechan chan int) {
	//returner true hvis jeg er master
	//brukes til å sjekke hvem som er master basert på lavest IP
	//returnerer false hvis jeg ikke har lavest IP
	lowestIP = localIP

	for {
		<-read 

		for key,_ := range addresses{
			s1 := strings.SplitAfterN(key,".",-1)
			s2 := strings.SplitAfterN(lowestIP,".",-1)
			IP1,_ := strconv.Atoi(s1[3])
			IP2,_ := strconv.Atoi(s2[3])

			if (IP1 < IP2) && IP1 > 0 && IP2 > 0{
				lowestIP = key
			}
		}
<<<<<<< HEAD
			
		if lowestIP == localIP && !master {
			masterchan<-1
		}
		else if !slave {
			slavechan <-1
		}
		
=======
	}
	fmt.Println(lowestIP)
	if lowestIP == localIP{
		return true
	}else{
		return false
	}
>>>>>>> cebda29012a2f236b513ad88fe18cf54e78ecccd

		read<-1
		time.Sleep(10*time.Millisecond)
	}
}

/*

func OrderNotInList([]orders ElevLib.MyOrder, neworder ElevLib.MyOrder) (bool) {
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
	broadcastAliveaddr,_ := net.ResolveUDPAddr("udp", localHost+BRALIVE)
	broadcastAliveSock,_ := net.DialUDP("udp", nil, broadcastAliveaddr)
	time.Sleep(10*time.Millisecond)
	for {
		_,err := broadcastAliveSock.Write([]byte(localIP))
		if err != nil{
			fmt.Println("ERRORR!", "SendAliveMessageUDP closing")
			return
		}
		time.Sleep(10*time.Millisecond)
	}
	broadcastAliveSock.Close()
}



func ReadAliveMessageUDP(write chan int){
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
	buffer := make([]byte,1024)
	for {
		<-write
		conn.ReadFromUDP(buffer)
		s := string(buffer[0:15]) //slipper nil i inlesningen
		fmt.Print(s)
		addresses[string(s)] = time.Now()
		if s!= "" {
			for key, value := range addresses{
				if time.Now().Sub(value) > 100*time.Millisecond && key != localIP{
					delete(addresses,key)
				}
			}
			
		}write<-1
		time.Sleep(10*time.Millisecond)
	}
	conn.Close()
}

func PrintAddresses() {
	for key,_ := range addresses {
		fmt.Println(key)
	}
}

func broadCastOrder(order ElevLib.MyOrder) {
	broadcastOrderaddr,_ := net.ResolveUDPAddr("udp", localHost + BRORDER)
	broadcastOrderSock,_ := net.DialUDP("udp", nil, broadcastOrderaddr)
	time.Sleep(10*time.Millisecond)
	for i:=0;i<10;i++ {
		buf,_ := json.Marshal(order)
		_,err := broadcastOrderSock.Write(buf)
		if err != nil{
			panic(err)
		}
	}
	broadcastOrderSock.Close()
}

func RecieveOrders(orderchannel chan ElevLib.MyOrder) {
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", localHost + BRORDER)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for { 
		msglen ,_,_ := recieveSock.ReadFromUDP(buffer)
		var tempOrder ElevLib.MyOrder
		json.Unmarshal(buffer[:msglen], &tempOrder)
		orderchannel <- tempOrder
		time.Sleep(10*time.Millisecond)
	}
}

//////////////////////////TCP funksjoner/////////////////////////////////
/////////////////////////////////////////////////////////////////////////

<<<<<<< HEAD
func Master(sendInfo chan ElevLib.MyInfo, extOrder chan ElevLib.MyOrder , PanelOrder chan ElevLib.MyOrder, slaveChan chan int, closing chan int, stopTCP chan int, stopRead chan int) {
=======
func Master(sendInfo chan ElevLib.MyInfo, extOrder chan ElevLib.MyOrder , PanelOrder chan ElevLib.MyOrder, recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder) {
>>>>>>> cebda29012a2f236b513ad88fe18cf54e78ecccd

	//var orders := []queue.MyOrder{}
	//recvInfo := make(chan ElevLib.MyInfo)
	//recvOrder := make(chan ElevLib.MyOrder)
	
	masterchange := make(chan bool)

	time.Sleep(10*time.Millisecond)
	go masterToSlaveMode(masterchange)
	//go ReadALL(writeToSocketMap, recvInfo, recvOrder)
	fmt.Println("MASTER:", "Going on")
	for {


		//PrintAddresses()
		fmt.Println("")
		time.Sleep(1*time.Second)
		/*select{/*
			case NewInfo := <-recvInfo:
				//OPPDATERE INFOMAP MED INFOEN MOTTATT PÅ SOCKET
				infomap[NewInfo.Ip] = NewInfo
			case NewOrder := <-recvOrder:
				handledorder := orderhandler(NewOrder)

				broadCastOrder(handledorder)

			case Ownorder := <- PanelOrder:

				handledorder := orderhandler(Ownorder)
				fmt.Println("NETWORK: ","new panel Order recieved: ")
				if handledorder.Ip == localIP {
					extOrder <- handledorder
				}
				broadCastOrder(handledorder)
			case UpdateInfo := <- sendInfo:
				infomap[localIP] = UpdateInfo
			case <- slaveChan:
				fmt.Println("Going slavemode")
				stopTCP<-1
				stopRead<-1
				closing<-1
				return
			default:
				time.Sleep(1*time.Second)
				fmt.Println("MASTER")
				PrintAddresses()
		}*/
	}
}




func orderhandler(order ElevLib.MyOrder)(ElevLib.MyOrder) {

	//var besteheis ElevLib.MyInfo
	//order.Ip = localIP
	//return order
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

	}*/
	for _,value := range infomap{
		order.Ip = value.Ip
		return order
	}
	order.Ip = localIP
	return order
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

func ReadALL(writing chan int, recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder, stopRead chan int) {
	for  {
<<<<<<< HEAD

		select{
		
		case <-writing: 
			for _,connection := range socketmap{
				buffer := make([]byte,1024)
				msglen ,_:= connection.Read(buffer)
				var temp ElevLib.MyElev
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
		case <- stopRead:
			return
=======
		fmt.Println("READALL:", "Cannot read from socketmap" )
		<-writing
		fmt.Println("READALL:", "REading from socketmap" )
		for _,connection := range socketmap{
			buffer := make([]byte,1024)
			msglen ,_:= connection.Read(buffer)
			var temp ElevLib.MyElev
			json.Unmarshal(buffer[:msglen], &temp)
			if temp.MessageType == "INFO" {
				recvInfo <-temp.Info
				writing<-1
				time.Sleep(time.Millisecond)
			}else if temp.MessageType == "ORDER" {
				recvOrder <-temp.Order
				writing<-1
				time.Sleep(time.Millisecond)
			}else{
				continue
			}
			time.Sleep(time.Millisecond)
>>>>>>> cebda29012a2f236b513ad88fe18cf54e78ecccd
		}
	}
}


func masterToSlaveMode( masterchange chan bool ){
	for {
		if !boolvar {
			masterchange<-true
		}
	}
}


func slaveToMasterMode(slavechange chan bool ){
	for {
		if boolvar {
			slavechange<-true
		}
	}
}
/*

func ReadOrders(chan recvOrder queue.MyOrder){
	for _,connection := range socketmap{
		buffer := make([]byte,1024)
		msglen ,_:= connection.Read(buffer)
		var tempOrder
	}
}
*/
func writetoSocket(socket *net.TCPConn, object ElevLib.MyElev)(bool){
	if object.MessageType == "INFO" {
		buffer,_ := json.Marshal(object.Info)
		_,err:= socket.Write(buffer)
		if err != nil {
			fmt.Println("error", err)
			return false
		//errorhandle
		}
		return true
	}else if object.MessageType == "ORDER" {
		buffer,_ := json.Marshal(object.Order)
		_,err:= socket.Write(buffer)
		if err != nil {
			fmt.Println("error", err)
			return false
			//errorhandle
		}
		return true
	}else{
		return false
	}
}



func Slave(sendInfo chan ElevLib.MyInfo, extOrder chan ElevLib.MyOrder, Panelorder chan ElevLib.MyOrder, masterchan chan int, closing chan int) {
	var masterSocket *net.TCPConn 
	var connected bool = false
	for(connected==false){
		masterSocket,connected = ConnectToIP(lowestIP)
	}
	recievechannel := make(chan ElevLib.MyOrder)
	var sendObject ElevLib.MyElev
	slavechange := make(chan bool)

	go slaveToMasterMode(slavechange)
	go RecieveOrders(recievechannel)
	for {
		if (boolvar) {
			fmt.Println("Going from slave to master")
			return
		}

		for {

			select{
			case NewOrder := <- recievechannel:
				if (NewOrder.Ip == localIP) {
					extOrder <- NewOrder
				}
			case NewPanelOrder := <- Panelorder:
				sendObject.MessageType = "ORDER"
				sendObject.Order = NewPanelOrder
				sendObject.Info = ElevLib.MyInfo{}

				sentorder := writetoSocket(masterSocket, sendObject)
				for !sentorder {
					sentorder = writetoSocket(masterSocket, sendObject)
				}

			case InfoUpdate := <- sendInfo:
				sendObject.MessageType = "INFO"
				sendObject.Order = ElevLib.MyOrder{}
				sendObject.Info = InfoUpdate

				sentinfo := writetoSocket(masterSocket, sendObject)
				for !sentinfo {
					sentinfo = writetoSocket(masterSocket, sendObject)
				}
			case <-masterchan:
				fmt.Println("Going from slave To Master!")
				closing<-1
				return
				//LAGE CASE FOR IKKE BOOLVAR!!!
			}
		}
		
	}

		//sende Ordre, og motta ordre
}

	//sender inn alle bestillinger den mottar fra panel til master
	//lytter på port for å motta en ordre fra master
	//setter inn bestilling i queue med hvilken IP som skal di


func TCPAccept(writeToSocket chan int, stopTCP chan int) {
	listenAddr, error := net.ResolveTCPAddr("tcp4", localIP+tcpPort)
	if error != nil {
		fmt.Println(error)
	}
	listener, error := net.ListenTCP("tcp4",listenAddr)
	if error != nil {
		fmt.Println(error)
	}
	for{
<<<<<<< HEAD
		select {

			case <-writeToSocket:
				listener.SetDeadline(time.Now().Add(time.Millisecond*100))
				remoteConn, error := listener.AcceptTCP()
				if (error == nil){
					socketmap[strings.Split(remoteConn.RemoteAddr().String(), ":")[0]] = remoteConn
				}
				time.Sleep(time.Millisecond)
				writeToSocket<-1

			case <-stopTCP:
				return
		}
=======
		<-writeToSocket
		fmt.Println("TCPAccept: ", "writing to socketmap")
		listener.SetDeadline(time.Now().Add(time.Millisecond*100))
		remoteConn, error := listener.AcceptTCP()
		if (error == nil){
			socketmap[strings.Split(remoteConn.RemoteAddr().String(), ":")[0]] = remoteConn
		}
		if (!boolvar) {
			return 
		}
		writeToSocket<-1
		time.Sleep(time.Millisecond)
>>>>>>> cebda29012a2f236b513ad88fe18cf54e78ecccd
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

func Network(newInfoChan chan ElevLib.MyInfo, externalOrderChan chan ElevLib.MyOrder, newExternalOrderChan chan ElevLib.MyOrder) {
	writeToSocketmap := make(chan int,1)
	//recvInfo := make(chan ElevLib.MyInfo)
	//recvOrder := make(chan ElevLib.MyOrder)



	writeToSocketmap <- 1
	master := SolvMaster()
	if (!master) {
		boolvar = true		
	}
	for {
		fmt.Println("else")

		if (master) {
			if (!boolvar) {
				fmt.Println("Im Master")
				boolvar = true
				//go TCPAccept(writeToSocketmap)
				time.Sleep(time.Millisecond)
				//go ReadALL(writeToSocketmap, recvInfo, recvOrder)
				go Master(newInfoChan, externalOrderChan, newExternalOrderChan, recvInfo, recvOrder)
			}
		master = SolvMaster()
		}else{
			fmt.Println("else")
			if (boolvar) {
				fmt.Println("Im a Slave biatch")
				boolvar = false
				go Slave(newInfoChan, externalOrderChan, newExternalOrderChan)
			}
		master = SolvMaster()
		}
		time.Sleep(10*time.Millisecond)
	}
	fmt.Println("IIm dead")
}


func Network2(newInfoChan chan ElevLib.MyInfo, externalOrderChan chan ElevLib.MyOrder, newExternalOrderChan chan ElevLib.MyOrder) {

	writeToSocketMap := make(chan int,1)
	readAndWriteAdresses := make(chan int, 1)
	masterChan := make(chan int,bool)
	changeChan := make(chan int,1)

	go SendAliveMessageUDP()
	go ReadAliveMessageUDP(readAndWriteAdresses)
	go SolvMaster(readAndWriteAdresses, masterChan)
	if check := <-masterChan; check == true{
		fmt.Println("Im Master")
		go TCPAccept(writeToSocketmap)
		go ReadALL(writeToSocketMap)
		go Master(newInfoChan,externalOrderChan, newExternalOrderChan, Masterchan , changeChan)
		<-changechan
	}
	else{
		fmt.Println("Im a Slave")
		go Slave(newInfoChan, externalOrderChan, newExternalOrderChan)
	}
	time.Sleep(100*time.Millisecond)
	
}


func Network3(newInfoChan chan ElevLib.MyInfo, externalOrderChan chan ElevLib.MyOrder, newExternalOrderChan chan ElevLib.MyOrder) {
	
	writeToSocketMap := make(chan int,1)
	readAndWriteAdresses := make(chan int, 1)
	masterChan := make(chan int,1)
	slaveChan := make(chan int,1)
	closing := make(chan int, 1)
	stopTCP := make(chan int, 1)
	stopRead := make(chan int, 1)


	go SendAliveMessageUDP()
	go ReadAliveMessageUDP(readAndWriteAdresses)
	time.Sleep(time.Millisecond)
	go SolvMaster(readAndWriteAdresses, masterchan, slaveChan)

	for {

		select {
			case <- masterChan:

				master = true
				fmt.Println("IM A MASTER")
				go ReadALL(writeToSocketMap)
				go TCPAccept(writeToSocketMap)
				go Master(newInfoChan, externalOrderChan, newExternalOrderChan, slaveChan, closing, stopTCP, stopRead)
				<- closing
				master = false

			case <- slaveChan:

				fmt.Println("IM SLAVE")
				go Slave(newInfoChan, externalOrderChan, newExternalOrderChan, masterchan, closing)
				<- closing


		time.Sleep(time.Millisecond)
		}

	}


}
