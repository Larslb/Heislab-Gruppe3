package Driver
import (
	
	//"Network"
	//"fmt"
	//"time"
)

func readElevPanel(buttonChan chan myOrder){
	for {
		for i:=0;i<N_FLOORS;i++{
			elev_get_button_signal(BUTTON_COMMAND,i, buttonChan)
		}
	}
}

func readFloorPanel(){
	for{
		for i:=0;i<N_BUTTONS-1;i++{
			for j:=0;j<N_FLOORS;j++{
				elev_get_button_signal(i, j, buttonChan)
			}
		}
	}
}

func readSensors(sensorChan chan int){
	for {
		for i:=0;i<N_FLOORS;i++{
			elev_get_floor_sensor_signal(sensorChan)
		}
	}
}

func readStopSignal(stopChan chan int){
	for {
		elev_get_stop_signal(stopChan)
	}
}





func Driver(){

	buttonFloorChan := make(chan myOrder)
	buttonElevChan  := make(chan myOrder) 
	stopChan := make(chan int)
	sensorChan := make(chan int)

	go readElevPanel(buttonElevChan)
	go readSensors(sensorChan)
	go readStopSignal(stopChan)
	go readFloorPanel(buttonFloorChan)
	
	for {
		select {
			// case: motta bestillinger fra master??
			
			case button := <- buttonElevChan
				elev_set_button_lamp(button.buttonType, button.floor, button.value)
				// SETTE BESTILLING I KØ
				
			
			
			case button := buttonFloorChan
				//SEND TIL MASTER
			
			case stop := <- stopChan  //BURDE STOP CASE OVERSTYRE DE ANDRE?
				elev_set_stop_lamp(stop)
				// GJØR NOE MER?
			
			
			case sensor := <- sensorChan
				if sensor != -1{
					elev_set_floor_indicator(sensor)
				}
		}
	}	
}
