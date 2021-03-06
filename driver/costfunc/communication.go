package elevlog

import(.".././network"
		"strconv"
		".././queue/"
		".././heis/")

func ReceiveMessage(){
	//Finds out what to do with message
	var received *Message
	var order driver.ButtonEvent

	for {
		received = <- NewMessage
		
		command := received.Message[:2]

		if (command == "oc"){
			//Ordercost from one of the slaves
			var temp OrderCost
			temp.Cost, _ = strconv.Atoi(string(received.Message[2]))
			temp.Conn = FindConn(received.From)
			CheckCost(temp)

		}else if (command == "no"){
			//New order from one of the slaves
			if (Master){
				order.Floor, _ = strconv.Atoi(string(received.Message[3]))
				order.Button = string(received.Message[2])
				newOrderMaster(order)
				newMess := "ac" + received.Message
				SendMessage(newMess,FindConn(received.From),false)
			}

		}else if (command == "cc"){
			//Master tells slaves to calulate cost of a order
			order.Floor, _ = strconv.Atoi(string(received.Message[2]))
			order.Button = string(received.Message[3])
			lightsOn(order)	
			slaveCalculate(order)

		}else if (command == "eo"){
			//Master tells one of the slaves to execute a certain order
			floor, _ := strconv.Atoi(string(received.Message[2]))
			queue.AddOrder(floor, string(received.Message[3]))

		}else if (command == "ba"){
			//Master sends the backup to the slaves
			receiveBackup(received.Message[2:])

		}else if (command == "nm"){
			//New master on the work, distributes orders in backup
			distributeBackup()

		}else if (command == "rm"){
			//One of the slaves tells all the other units that a order has been completed
			if (Master){
				RemoveBUOrder(received)
			}
			lightsOff(int(received.Message[2]))
		}
	}
}


func sendBackup(){
	message := "ba"
	for i := 0 ; i < 4 ; i++ {
		message = message + strconv.Itoa(backup.UP[i])
	}
	for i := 0 ; i < 4 ; i++ {
		message = message + strconv.Itoa(backup.DOWN[i])
	}
	for i := 0 ; i < 4 ; i++ {
		message = message + strconv.Itoa(backup.CMD[i])
	}

	SendMessage(message, Broadcast.Conn,false)
}


func receiveBackup(message string){
	for i := 0 ; i < 4 ; i++ {
		letter1, _ := strconv.Atoi(string(message[i]))
		letter2, _ := strconv.Atoi(string(message[i+4]))
		letter3, _ := strconv.Atoi(string(message[i+8]))

		backup.UP[i] = letter1
		backup.DOWN[i] = letter2
		backup.CMD[i] = letter3
	}
}


func distributeBackup(){
	//New master on the network or slave-timeout. Distributes all orders in backup
	var temp driver.ButtonEvent
	
	for i := 0 ; i < 4 ; i++{
		if (backup.UP[i] == 1){
			temp.Floor = i
			temp.Button = "U"
			newOrderMaster(temp)
		}
	
		if (backup.DOWN[i] == 1) {
			temp.Floor = i
			temp.Button = "D"
			newOrderMaster(temp)
		}
	}
}
