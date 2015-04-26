package Notify
import(
	"fmt"
	"os"
	"os/signal"
	".././Driver"
	".././Queue"
	"syscall"
)


func cleanup(backup ElevLib.MyInfo) {
    fmt.Println("cleanup!!")
    Driver.Elev_set_speed(0)
    fo, _ := os.Create("Backup.txt")
    fo.Write(backup)
   	fo.Close()
}

func Notify(notify chan chan ElevLib.MyInfo) {
	c := make(chan os.signal, 1)
	done := make(chan bool, 1)

	getbackup := make(chan ElevLib.MyInfo)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM )

	go func () {
		sigs := <-c	
		notify <- getbackup

		backup := <- getbackup

		cleanup(backup)
		fmt.Println(sig)
		done<-true
	}()

	<-done
	fmt.Println("Interrupt: exiting")
}