package Driver
import (
	
	//"time"
	//"Queue"

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
	tmpSensorChan := make(chan int, 1) // Nødvendig med buffer?
	for {
		elev_get_floor_sensor_signal(tmpSensorChan)
		tmpVal := <- tmpSensorChan
		if tmpVal != -1 {
			sensorChan <- tmpVal		
		}
	}
}

func readStopSignal(stopChan chan int){
	for {
		elev_get_stop_signal(stopChan)
	}
}

func floorReached





func Driver(){

	buttonFloorChan := make(chan MyOrder)
	buttonElevChan  := make(chan MyOrder) 
	stopChan := make(chan int)
	sensorChan := make(chan int)

	go readElevPanel(buttonElevChan)
	go readSensors(sensorChan)
	go readStopSignal(stopChan)
	go readFloorPanel(buttonFloorChan)
	
	for {
		select {
			// case: motta bestillinger fra master??
			//need to 
			
			case button := <-buttonElevChan:
				elev_set_button_lamp(button.ButtonType, button.Floor, button.Value)
				// 1. Sette bestilling i intern kø (bestillinger som bare denne heisen
				//    kan bruke).

				// HVA MED PRIORITERING??
				Queue.SetInternalOrders(button.Floor)
				
			
			
			case button := <-buttonFloorChan:
				// 1. Sette bestilling i ekstern kø som må leses av Nettverk og
				//    sendes til MASTER.
				Queue.SetExternalOrders(button.ButtonType, button.Floor)
			
			case stop := <-stopChan:  //BURDE STOP CASE OVERSTYRE DE ANDRE CASENE?
				elev_set_stop_lamp(stop)
				elev_set_motor_dir(STOP)
				Queue.Delete_all_internalOrders()
			
			
			case sensor := <-sensorChan:
				elev_set_floor_indicator(sensor)
				if Queue.CheckFloor(sensor) {
					elev_set_motor_dir(0)
					elev_set_door_open(true)
					time.Sleep(3*time.Second)
					Queue.DeleteInternalOrder(sensor)
					//send complete to master
				}
				
				
		}
	}
}
