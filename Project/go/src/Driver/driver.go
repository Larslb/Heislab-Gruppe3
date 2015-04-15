package Driver
import (
	
	//"time"
	"Queue"

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



//func floorReached () {}

func Driver(floorChan chan []int, elevChan [][]string ,sensChan []int){

	buttonFloorChan := make(chan MyOrder)
	buttonElevChan  := make(chan MyOrder) 
	sensorChan := make(chan int)

	go readElevPanel(buttonElevChan)
	go readSensors(sensorChan)
	go readFloorPanel(buttonFloorChan)
	
	for {
		select {
			// case: motta bestillinger fra master??
			//need to 
			
			case button := <-buttonElevChan:
				elev_set_button_lamp(button.ButtonType,button.Floor,1)
				// 1. Sette bestilling i intern kø (bestillinger som bare denne	 					//    heisen
	 			//    kan bruke).

				// HVA MED PRIORITERING??
				Queue.SetInternalOrders(button.Floor)

				floorChan <- Queue.GetInternalOrders()
				
			case button := <-buttonFloorChan:
				// 1. Sette bestilling i ekstern kø som må leses av Nettverk og
				//    sendes til MASTER.

				//master skal egentlig bestemme dette
				elev_set_button_lamp(button.ButtonType,button.Floor,1) 
				Queue.SetExternalOrders(button.ButtonType, button.Floor)
				
				elevChan <- Queue.GetUnprExtOrd()

			
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
