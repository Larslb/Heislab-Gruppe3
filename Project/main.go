
package main


import(
	"net"
	"fmt"
	"strings"
	//"os"

)

func GetLocalIP() (string){
	addr,_ := net.ResolveTCPAddr("tcp4", "google.com:80")
	conn,_ := net.DialTCP("tcp4", nil, addr)
	return strings.Split(conn.LocalAddr().String(), ":")[0] 
}
func main() {

	var localIP string = "0"
	localIP = GetLocalIP()

	fmt.Println(localIP)
}