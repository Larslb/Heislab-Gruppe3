package Driver
import (
	
	//"Nettwork"
	//"fmt"
	//"time"
)

func ElevPanelThread(buttonchannel chan int){
	for {
		for i:=0;i<N_FLOORS;i++{
			elev_get_button_signal(BUTTON_COMMAND,i, buttonchannel)
		}
	}
}
func SensorThread(sensorChan chan int){
	
}
func StopThread(stopChan chan int){

}
func FloorPanelThread(buttonchan chan int){

}



func main(){
	
	buttonFloorChan := make(chan int)
	buttonElevChan := make(chan int)
	stopChan := make(chan int)
	sensorChan := make(chan int)


	go elevPanelThread(buttonElevChan)
	go sensorThread(sensorChan)
	go stopThread(stopChan)
	go floorPanelThread(buttonFloorChan)	
	
	
	elev.init()
	bestilling <- buttonElevChan
	println(bestilling)
	
}
