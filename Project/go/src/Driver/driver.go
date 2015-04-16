package Driver
import (
	
	//"time"
	"Queue"

)

func readElevPanel(buttonChan chan Queue.MyOrder){
	for {
		for i:=0;i<N_FLOORS;i++{
			elev_get_button_signal(BUTTON_COMMAND,i, buttonChan)
		}
	}
}

func readFloorPanel(buttonChan chan Queue.MyOrder){
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



// BURDE KOMMUNIKASJONEN MELLOM DRIVER-MODULEN OG QUEUE-MODULEN SKJE MED CHANNELS I ET NIVÅ OPP
// I STEDET FOR AT VI BRUKER PUBLIC SET OG GET FUNKSJONER FRA QUEUE???
// DET SAMME GJELDER FOR SFM NÅR VI MÅ LESE BESTILLINGER FOR Å BESTEMME DIRECTION!

// FORSLAG:
// 1. GO-ROUTINE SOM LESER INTERNE OG EKSTERNE BESTILLINGER OG RETURNERER NEXT 
//	  FLOOR GJENNOM EN CHANNEL


func ReadPanelAndSensor(intrOrd chan Queue.MyOrder, extrOrd chan Queue.MyOrder, snsrChan chan int){

	// Har som oppgave å lese knappetrykk/sensor og distribuere bestillinger/sensor ut til EventManager


	buttonFloorChan := make(chan Queue.MyOrder)
	buttonElevChan  := make(chan Queue.MyOrder) 
	sensorChan := make(chan int)

	go readElevPanel(buttonElevChan)
	go readSensors(sensorChan)
	go readFloorPanel(buttonFloorChan)
	
	for {
		select {
			
		case button := <-buttonElevChan:
				//elev_set_button_lamp(button.ButtonType,button.Floor,1)
			intrOrd <- button
				//Queue.SetInternalOrders(button.Floor)
				
		case button := <-buttonFloorChan:

			extrOrd <- button

				//master skal egentlig bestemme dette
				//elev_set_button_lamp(button.ButtonType,button.Floor,1) 
				//Queue.SetExternalOrders(button.ButtonType, button.Floor)
			
		case sensor := <-sensorChan:

			snsrChan <- sensor

				// BURDE STÅ I FSM-koden
				//elev_set_floor_indicator(sensor)
				//if Queue.CheckFloor(sensor) {
				//	elev_set_motor_dir(0)
				//	elev_set_door_open(true)
				//	time.Sleep(3*time.Second)
				//	Queue.DeleteInternalOrder(sensor)
				//send complete to master
		}
	}
}


func setLights() {
	
}