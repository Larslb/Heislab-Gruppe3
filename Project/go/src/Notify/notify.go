package notify
import(
	"fmt"
	"os"
	"os/signal"
	".././Driver"
	".././Queue"
	"syscall"
)


func cleanup() {
    fmt.Println("cleanup!!")
    Driver.Elev_set_speed(0)
    fo, _ := os.Create("Backup.txt")
    fo.Write(Queue.GetInternalOrders())
   	fo.Close()
}



func Notify() {
	c := make(chan os.signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func () {
		sig := <-c
		cleanup()
		fmt.Println(sig)
		done<-true
	}
	<-done
	fmt.Println(exiting)
}