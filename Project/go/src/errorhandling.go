

// I MAIN

errorChan := make(chan string)


// Channels to communicate with errorFunctions

go Notify.Notify(notify)

for {
	select {

		case err  := <- myError:
			if err == "no connection to network" {  // denne typen error må trigge i Network
				connectionChan <- false 			// der udp/tct-connection leses

				// Når QM trigger på connection chan, må QM sette connection2network = false


			} else if  backupChan := <-notify {

				
				// MÅ LEGGE TIL EN CASE I QM SOM DEN KAN TRIGGE PÅ BACKUP
				backup := 


				// TRENGER JEG os.Exit(1) ?? 

			}			
	}
}



/*
1. Nettverkskabel dratt ut: Trigger på timeout.
		-> Handling 1: Fortsette å betjene de ordrene som den har mottat
		-> Handling 2: ikke ta imot externe bestillinger mens den fortsetter å håndtere gjenværende ordre
				-> må signalisere til QM at ikke er noen connections slik at eksterne bestillinger
				   ikke blir håndtert  (f.eks at networkConnection = false når connection lost)
		-> Handling 3: Ikke ta nye ordre når lastOrderHandled og !networkConnection

2. cntrl C trykket inn:
	-> Handling 1: Be QM om backup
	-> Handling 2: Skrive til backupfil
	-> Handling 3: kalle på Notify
	-> Handling 4: Avslutt program

*/



