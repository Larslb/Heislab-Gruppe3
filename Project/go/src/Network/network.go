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

		if lowestIP == localIP && !master {
			masterchan<-1
		}else if lowestIP != localIP && !slave {
			slavechan <-1
		}

	read<-1
	time.Sleep(10*time.Millisecond)	
	}
}


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



func Master(sendInfo chan ElevLib.MyInfo, extOrder chan ElevLib.MyOrder , PanelOrder chan ElevLib.MyOrder, slaveChan chan int, closing chan int, stopTCP chan int, stopRead chan int, recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder, orderdeletion chan ElevLib.MyOrder, orderDelFromMaster chan ElevLib.MyOrder) {
	
	time.Sleep(10*time.Millisecond)
	for {
		
		select{
			case NewInfo := <-recvInfo:
				infomap[NewInfo.Ip] = NewInfo
				
			case NewOrder := <-recvOrder:
				NewOrder.Ip = costfunction(NewOrder)
				if NewOrder.Set {
					extOrder <- NewOrder
				}else{
					orderDelFromMaster <- NewOrder
				}
				broadCastOrder(NewOrder)

			case Ownorder := <- PanelOrder:
				Ownorder.Set = true
				Ownorder.Ip = costfunction(Ownorder)
				extOrder <- Ownorder
				broadCastOrder(Ownorder)

			case UpdateInfo := <- sendInfo:
				infomap[localIP] = UpdateInfo
				
			case OrderDeleted := <- orderdeletion:
				broadCastOrder(OrderDeleted)

			case <- slaveChan:
				fmt.Println("Going To Slavemode")
				stopTCP<-1
				stopRead<-1
				closing<-1
				return
		}
		time.Sleep(10*time.Millisecond)
	}
}


func Slave(sendInfo chan ElevLib.MyInfo, extOrder chan ElevLib.MyOrder, Panelorder chan ElevLib.MyOrder, masterchan chan int, closing chan int, stopRecieving chan int, orderdeletion chan ElevLib.MyOrder , orderDelFromMaster chan ElevLib.MyOrder) {
	var masterSocket *net.TCPConn 
	var connected bool = false
	for(connected==false){
		masterSocket,connected = ConnectToIP(lowestIP)
	}
	recievechannel := make(chan ElevLib.MyOrder)
	go RecieveOrders(recievechannel, stopRecieving)
	
	for {
		select{
			case NewOrder := <- recievechannel:

				fmt.Println(NewOrder.Ip, NewOrder.ButtonType, NewOrder.Floor)
				if NewOrder.Set {
					extOrder <- NewOrder
				}else {
					orderDelFromMaster <- NewOrder
				}

			case NewPanelOrder := <- Panelorder:


				NewPanelOrder.Set = true
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
				//PrintAddresses()

				sentinfo := writetoSocket(masterSocket, sendObject)

				for !sentinfo {
					sentinfo = writetoSocket(masterSocket, sendObject)
				}
				fmt.Println("info sent")
			case deleteOrder := <-orderdeletion:
				sendObject := ElevLib.MyElev {
					MessageType: "ORDER",
					Order: deleteOrder,
					Info: ElevLib.MyInfo{},
				}

				sentorder := writetoSocket(masterSocket, sendObject)
				for !sentorder {
					sentorder = writetoSocket(masterSocket, sendObject)
				}


			case <-masterchan:
				fmt.Println("Going from slave To Master!")
				stopRecieving<-1
				closing<-1
				return

		}
		time.Sleep(10*time.Millisecond)
	}
		
}






func readfromsocket( conn *net.TCPConn,  recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder ) bool {
	buffer := make([]byte,1024)
	conn.SetReadDeadline(time.Now().Add(80*time.Millisecond))
	
	msglen ,err:= conn.Read(buffer)
	if err != nil {
		time.Sleep(10*time.Millisecond)
		return false
	}
	var temp ElevLib.MyElev
	json.Unmarshal(buffer[:msglen], &temp)

	if temp.MessageType == "INFO" {
		recvInfo <-temp.Info
		return true
	}else if temp.MessageType == "ORDER" {
		recvOrder <-temp.Order
		return true
	}
	return false
	
}
func ReadALL(writing chan int, recvInfo chan ElevLib.MyInfo, recvOrder chan ElevLib.MyOrder, stopRead chan int) {
	for  {
		select{
		case <-writing:
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
		buffer,_ := json.Marshal(object)
		_,err:= socket.Write(buffer)
		if err != nil {
			fmt.Println("error", err)
			return false
		}
		return true
	}else if object.MessageType == "ORDER" {
		buffer,_ := json.Marshal(object)
		_,err:= socket.Write(buffer)
		if err != nil {
			fmt.Println("error", err)
			return false
		}
		return true
	}else{
		return false
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
				listener.SetDeadline(time.Now().Add(time.Millisecond*100))
				remoteConn, error := listener.AcceptTCP()
				if (error == nil){
					socketmap[strings.Split(remoteConn.RemoteAddr().String(), ":")[0]] = remoteConn
				}
				writeToSocket<-1
				time.Sleep(time.Millisecond)
			case <-stopTCP:
				return
		}
		time.Sleep(10*time.Millisecond)
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


func Network(newInfoChan chan ElevLib.MyInfo, externalOrderChan chan ElevLib.MyOrder, newExternalOrderChan chan ElevLib.MyOrder, masterChan chan int, slaveChan chan int, orderdeletion chan ElevLib.MyOrder ,ordrDeleteFromMaster chan ElevLib.MyOrder) {
	
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
				go Master(newInfoChan, externalOrderChan, newExternalOrderChan, slaveChan, closingMaster, stopTCP, stopRead, recvInfo, recvOrder, orderdeletion,ordrDeleteFromMaster)
				<- closingMaster
				master = false

			case <- slaveChan:

				slave = true
				fmt.Println("IM SLAVE")
				go Slave(newInfoChan, externalOrderChan, newExternalOrderChan, masterChan, closingSlave, stopRead, orderdeletion, ordrDeleteFromMaster)
				time.Sleep(time.Millisecond)
				<- closingSlave
				time.Sleep(time.Millisecond)
				slave = false


		time.Sleep(time.Millisecond)
		}

	}


}

func costfunction( order ElevLib.MyOrder) string {
	
	elevsInDirection := inDirection(order)

	nearestElev := shortestRoute(order, elevsInDirection)

	return nearestElev	


}

func inDirection(order ElevLib.MyOrder ) []string {
	elevs := []string{}
	var orderDirectionRelativeElev int

	for key, val := range infomap {
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

func shortestRoute(order ElevLib.MyOrder, elevlist []string ) string {

	if length := len(elevlist); length > 2 {
		m := length/2
		list1 := shortestRoute(order, elevlist[0:m])
		list2 := shortestRoute(order, elevlist[m:length])

		dist1 := abs(infomap[list1].CurrentFloor-order.Floor)
		dist2 :=  abs(infomap[list2].CurrentFloor-order.Floor)

		if dist1 < dist2 {
			return list1
		} else { 
			return list2
		}

	} else if len(elevlist) == 1 {
		return elevlist[0]
	} else {
		dist1 := abs(infomap[elevlist[0]].CurrentFloor-order.Floor)
		dist2 :=  abs(infomap[elevlist[1]].CurrentFloor-order.Floor)

		if dist1 < dist2 {
			return elevlist[0]
		} else { 
			return elevlist[1]
		}
	}

	return elevlist[0] 
}

func abs(number int) int {
	if number < 0 {
		return -number
	}
	return number
}
