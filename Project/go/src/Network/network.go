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
	BRPORT string = "25556"
	tcpPort string = "25557"
	) 

var localIP string = "0"
var lowestIP string = "0"
var infomap = make(map[*net.TCPAddr]queue.MyInfo)
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
   conn.close()
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



///////////////////////UDP funksjoner/////////////////////////////////
//////////////////////////////////////////////////////////////////////

func SendAliveMessageUDP(){
	broadcastAliveaddr,_ := net.ResolveUDPAddr("udp", localHost + BRPORT)
	broadcastAliveSock,_ := net.DialUDP("udp", nil, broadcastAliveaddr)
	time.Sleep(10*time.Millisecond)
	for {
		_,err = broadcastAliveSock.Write([]byte(localIP))
		if err != nil{
			panic(err)
		}
		time.Sleep(10*time.Millisecond)
	}
	broadcastAliveSock.close()
}



func ReadAliveMessageUDP(readvar bool){
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
			readvar=true
			for key, value := adresses{
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





//////////////////////////TCP funksjoner/////////////////////////////////
/////////////////////////////////////////////////////////////////////////

func Master() {

	for {
		if (!boolvar) {
			fmt.Println("Going Slavemode")
			return //return stopper threaden
		}

		//Trenger en kanal for å legge til egne bestillinger fra internkø

		//LESER INN INFO FRA ALLE HEISENE
		for _,connection := range socketmap{
			info := ReadInfo(connection)
			infomap[connection.LocalAddr()] = info
		}

		go ReadOrders() // må ta in en kanal

		
		for _,socket := range socketmap {
			//sende orderen til alle, med info om hvem som skal ta seg av ordren
		}
		


	}


	//leser inn info fra heiser
	//leser inn ordre fra heiser.. Kanskje en kanal som står å venter på ordre fra riktig port?
	//kjører costfunksjon
	//Broadcaster orderen til alle, med rett IP 

}

func Slave() {
	

	//sender inn alle bestillinger den mottar fra panel til master
	//lytter på port for å motta en ordre fra master
	//setter inn bestilling i queue med hvilken IP som skal dit
}


func TCPAccept() {
	listenAddr, error := net.ResolveTCPAddr("tcp", localIP+tcpPort)
	if error != nil {
		fmt.Println(error)
	}
	listener, error := net.ListenTCP("tcp",listenAddr)
	if error != nil {
		fmt.Println(error)
	}
	for{
		listener.SetDeadline(time.Now().Add(time.Millisecond*100))
		remoteConn, error := listener.AcceptTCP()
		if error == nil {
			socketmap[]/////UFERDIG!!
		}
	}
}
}
func ConnectToIP(IP string)(*net.TCPConn, bool){
	remoteAddr,error := net.ResolveTCPAddr("tcp", IP + tcpPort)
	if error != nil{
		fmt.Println(error)
		panic(error)
	}
	conn, error := net.DialTCP("tcp", nil, remoteAddr) 
    if(error==nil){
    	return conn, true
    }else{
    	return conn, false
    }

}

///////////////////////////////diverse funksjoner/////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

func SendInfoMessageUDP(info chan Queue.MyInfo){
	addr,_ := net.ResolveUDPAddr("udp", lowestIP + COMPORT)
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

func Listen() {
}


func EventManager_NetworkStuff() {  //Pseudokode

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
			go master()
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



