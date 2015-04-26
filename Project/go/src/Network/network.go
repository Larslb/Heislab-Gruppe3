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
var deadadresses []string 
var master bool = false
var slave bool = false



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
			//sfmt.Print(key)
			s1 := strings.SplitAfterN(key,".",-1)
			s2 := strings.SplitAfterN(lowestIP,".",-1)
			IP1,_ := strconv.Atoi(s1[3])
			IP2,_ := strconv.Atoi(s2[3])

			if (IP1 < IP2) && IP1 > 0 && IP2 > 0{
				lowestIP = key
			}
		}
		if lowestIP == localIP && !master {
			masterchan<-1
		}else if lowestIP != localIP && !slave {
			slavechan <-1
		}

	//fmt.Println("MASTERIP :",lowestIP)
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
		//fmt.Print(s)
		addresses[string(s)] = time.Now()
		if s!= "" {
			for key, value := range addresses{
				if time.Now().Sub(value) > time.Second && key != localIP{
					if key == lowestIP {
						lowestIP = localIP
					}
					delete(addresses,key)
					deadadresses = append(deadadresses, key)

				}
			}
			
		}
		write<-1
		time.Sleep(10*time.Millisecond)
	}
	conn.Close()
}

func PrintAddresses() {
	for key,_ := range addresses {
		fmt.Println(key)
	}
}

func PrintDeadAddresses() {
	for key := 0; key <len(deadadresses); key++ {
		fmt.Println(deadadresses[key])
	}
}

func printInfo() {
	for _,value := range infomap {
		fmt.Println(value.Ip, value.Dir, value.CurrentFloor, value.InternalOrders)
	}
}

func broadCastOrder(order ElevLib.MyOrder) {
	broadcastOrderaddr,_ := net.ResolveUDPAddr("udp", localHost + BRORDER)
	broadcastOrderSock,_ := net.DialUDP("udp", nil, broadcastOrderaddr)
	time.Sleep(10*time.Millisecond)
	for i:=0;i<10;i++ {
		fmt.Println("BROADCASTINGORDER!!!!")
		buf,_ := json.Marshal(order)
		_,err := broadcastOrderSock.Write(buf)
		if err != nil{
			panic(err)
		}
		time.Sleep(10*time.Millisecond)
	}
	broadcastOrderSock.Close()
}


func RecieveOrders(orderchannel chan ElevLib.MyOrder, stopRecieving chan int) {
	buffer := make([]byte,1024) 
	raddr,_ := net.ResolveUDPAddr("udp", localHost + BRORDER)
	recieveSock,_ := net.ListenUDP("udp", raddr)
	for { 
		select{
		default:
			msglen ,_,_ := recieveSock.ReadFromUDP(buffer)
			var tempOrder ElevLib.MyOrder
			json.Unmarshal(buffer[:msglen], &tempOrder)
			orderchannel <- tempOrder
			time.Sleep(10*time.Millisecond)
		case <-stopRecieving:
			return
		}
	}
}

//////////////////////////TCP funksjoner/////////////////////////////////
/////////////////////////////////////////////////////////////////////////

func Master(sendInfo chan ElevLib.MyInfo, extOrder chan ElevLib.MyOrder , PanelOrder chan ElevLib.MyOrder, slaveChan chan int, closing chan int, stopTCP chan int, stopRead chan int, recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder) {
	//var orders := []queue.MyOrder{}
	//recvInfo := make(chan ElevLib.MyInfo)
	//recvOrder := make(chan ElevLib.MyOrder)
	
	//masterchange := make(chan bool)

	time.Sleep(10*time.Millisecond)
	//go masterToSlaveMode(masterchange)
	//go ReadALL(writeToSocketMap, recvInfo, recvOrder)
	fmt.Println("MASTER:", "Going on")
	fmt.Println("")
	time.Sleep(1*time.Second)
	for {
		//PrintAddresses()
		select{
			case NewInfo := <-recvInfo:
				//OPPDATERE INFOMAP MED INFOEN MOTTATT PÅ SOCKET
				infomap[NewInfo.Ip] = NewInfo
				fmt.Println("NETWORK:   INFO RECIEVED!!")
			case NewOrder := <-recvOrder:
				handledorder := orderhandler(NewOrder)

				extOrder <- handledorder
				broadCastOrder(handledorder)

			case Ownorder := <- PanelOrder:

				handledorder := orderhandler(Ownorder)
				fmt.Println("NETWORK: ","new panel Order recieved: ")
				
				extOrder <- handledorder
				
				
				broadCastOrder(handledorder)
			case UpdateInfo := <- sendInfo:
				infomap[localIP] = UpdateInfo
			case <- slaveChan:
				fmt.Println("Going slavemode")
				stopTCP<-1
				stopRead<-1
				closing<-1
				return
				

		}
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

func readfromsocket( conn *net.TCPConn,  recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder ) bool {
	buffer := make([]byte,1024)
	conn.SetReadDeadline(time.Now().Add(80*time.Millisecond))
	
	msglen ,err:= conn.Read(buffer)
	if err != nil {
		time.Sleep(10*time.Millisecond)
		return false
	}
	//fmt.Println("READALL using socketmap")
	var temp ElevLib.MyElev
	json.Unmarshal(buffer[:msglen], &temp)


	fmt.Println(" ")
	fmt.Println("-------------------------")
	fmt.Println("RECIEVED  TEMP: ", temp.MessageType, temp.Info, temp.Order)
	fmt.Println("-------------------------")
	fmt.Println(" ")
	if temp.MessageType == "INFO" {
		fmt.Println("INFO recieved")
		recvInfo <-temp.Info
		return true
	}else if temp.MessageType == "ORDER" {
		fmt.Println("ORDER recieved")
		recvOrder <-temp.Order
		return true
	}
	return false
	
}
func ReadALL(writing chan int, recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder, stopRead chan int) {
	for  {
		select{
		
		case <-writing:
			fmt.Println("ReadALL using socketmap")
			for _,connection := range socketmap{
				readfromsocket(connection, recvInfo, recvOrder)
			}
			writing<-1
			time.Sleep(10*time.Millisecond)
		case <- stopRead:
			return
		}
	}
}

func writetoSocket(socket *net.TCPConn, object ElevLib.MyElev )(bool){
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



func Slave(sendInfo chan ElevLib.MyInfo, extOrder chan ElevLib.MyOrder, Panelorder chan ElevLib.MyOrder, masterchan chan int, closing chan int, stopRecieving chan int) {
	var masterSocket *net.TCPConn 
	var connected bool = false
	for(connected==false){
		masterSocket,connected = ConnectToIP(lowestIP)
	}
	recievechannel := make(chan ElevLib.MyOrder)
	//var sendObject ElevLib.MyElev
	//slavechange := make(chan bool)

	//go slaveToMasterMode(slavechange)
	go RecieveOrders(recievechannel, stopRecieving)
	fmt.Println("GOING IN FOR SELECT LOOP")
	for {

		for {

			select{
			case NewOrder := <- recievechannel:

				fmt.Println(NewOrder.Ip, NewOrder.ButtonType, NewOrder.Floor)
				extOrder <- NewOrder
			case NewPanelOrder := <- Panelorder:

				sendObject := ElevLib.MyElev {
					MessageType: "ORDER",
					Order: NewPanelOrder,
					Info: ElevLib.MyInfo{},
				}

				sentorder := writetoSocket(masterSocket, sendObject)
				for !sentorder {
					sentorder = writetoSocket(masterSocket, sendObject)
				}

			case InfoUpdate := <- sendInfo:
				fmt.Println("Sending InfoUpdate to master")
				sendObject := ElevLib.MyElev{
					MessageType: "INFO",
					Order: ElevLib.MyOrder{},
					Info: InfoUpdate,
				}

				fmt.Println("Sending: ", sendObject.MessageType, sendObject.Order, sendObject.Info)
				PrintAddresses()

				sentinfo := writetoSocket(masterSocket, sendObject)

				for !sentinfo {
					sentinfo = writetoSocket(masterSocket, sendObject)
				}
				fmt.Println("info sent")
			case <-masterchan:
				fmt.Println("Going from slave To Master!")
				stopRecieving<-1
				closing<-1
				return

			}
		}
		
	}

}



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
		select {

			case <-writeToSocket:
				//fmt.Println("Writing to sockets!")
				listener.SetDeadline(time.Now().Add(time.Millisecond*100))
				remoteConn, error := listener.AcceptTCP()
				if (error == nil){
					socketmap[strings.Split(remoteConn.RemoteAddr().String(), ":")[0]] = remoteConn
					fmt.Println("added in socketmap: ", strings.Split(remoteConn.RemoteAddr().String(), ":")[0])
				}
				writeToSocket<-1
				time.Sleep(time.Millisecond)
			case <-stopTCP:
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

func Network3(newInfoChan chan ElevLib.MyInfo, externalOrderChan chan ElevLib.MyOrder, newExternalOrderChan chan ElevLib.MyOrder, masterChan chan int, slaveChan chan int) {
	
	writeToSocketMap := make(chan int,1)
	recvInfo := make(chan ElevLib.MyInfo)
	recvOrder := make(chan ElevLib.MyOrder)
	closingMaster := make(chan int, 1)
	closingSlave := make(chan int, 1)
	stopTCP := make(chan int, 1)
	stopRead := make(chan int, 1)


	for {

		select {
			case <- masterChan:

				master = true
				fmt.Println("IM A MASTER")
				go TCPAccept(writeToSocketMap, stopTCP)
				time.Sleep(time.Millisecond)
				writeToSocketMap<-1
				go ReadALL(writeToSocketMap, recvInfo, recvOrder, stopRead)
				go Master(newInfoChan, externalOrderChan, newExternalOrderChan, slaveChan, closingMaster, stopTCP, stopRead, recvInfo, recvOrder)
				<- closingMaster
				master = false

			case <- slaveChan:

				slave = true
				fmt.Println("IM SLAVE")
				go Slave(newInfoChan, externalOrderChan, newExternalOrderChan, masterChan, closingSlave, stopRead)
				time.Sleep(time.Millisecond)
				<- closingSlave
				slave = false


		time.Sleep(time.Millisecond)
		}

	}


}

//////////////////////////////////////////////////////
//													//
// DET SOM HØRER TIL COST FUNCTION STÅR UNDER HER   //
//													//
//////////////////////////////////////////////////////
func costfunction( info map[string]ElevLib.MyInfo, order ElevLib.MyOrder) string {
	
	elevsInDirection := inDirection(info, order)

	nearestElev := shortestRoute(info, order, elevsInDirection)

	return nearestElev	

	/*
	if len(nearestElevs) == 1 { return elevsInDirection[0] }   

	bestElevs := fewestOrders(info, nearestElevs)	KAN OPPGRADERES TIL Å RETURNERE DEN MED FÆRREST BESTILLINGER

	return bestElevs[0]*/
}

func inDirection(info map[string]ElevLib.MyInfo, order ElevLib.MyOrder ) []string {
	elevs := []string{}
	var orderDirectionRelativeElev int

	for key, val := range info {
		if order.Floor < val.CurrentFloor {
			orderDirectionRelativeElev = -1
		} else if order.Floor > val.CurrentFloor {
			orderDirectionRelativeElev = 1
		} else {
			orderDirectionRelativeElev = 0
		}

		if val.Dir == orderDirectionRelativeElev || val.Dir == 0  {
			elevs = append(elevs, key)
		}
	}

	return elevs

}

func shortestRoute(info map[string]ElevLib.MyInfo, order ElevLib.MyOrder, elevlist []string ) string {

	if length := len(elevlist); length > 2 {
		m := length/2
		list1 := shortestRoute(info, order, elevlist[0:m])
		list2 := shortestRoute(info, order, elevlist[m:length])

		dist1 := abs(info[list1].CurrentFloor-order.Floor)
		dist2 :=  abs(info[list2].CurrentFloor-order.Floor)

		if dist1 < dist2 {
			return list1
		} else { 
			return list2
		}

	} else if len(elevlist) == 1 {
		return elevlist[0]
	} else {
		dist1 := abs(info[elevlist[0]].CurrentFloor-order.Floor)
		dist2 :=  abs(info[elevlist[1]].CurrentFloor-order.Floor)

		if dist1 < dist2 {
			return elevlist[0]
		} else { 
			return elevlist[1]
		}
	}

	return elevlist[0]  // HVIS IKKE VI RETURNERER PÅ SLUTTEN KLAGER KOMPILATOREN
}

func abs(number int) int {
	if number < 0 {
		return -number
	}
	return number
}
