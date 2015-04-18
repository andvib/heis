package costfunc

import(".././network"
		".././heis/"
		"strconv"
		".././ko/"
		".././event/")

var backup ko.Queue
var qCopy ko.Queue

func ButtonHandle(){
	var buttonEvent driver.ButtonEvent

	for ; true ; {
		buttonEvent = <- driver.ButtonChan
		//println("New order: ", buttonEvent.Floor, ", ",buttonEvent.Button)
		println("NewOrderCOSTFUNC")		
		newOrderSlave(buttonEvent)
	}	
}


func newOrderSlave(order driver.ButtonEvent){
	println("NewOrderSlave")
	message := "no" + order.Button + strconv.Itoa(order.Floor)
	println(message)
	//network.SendMessage(message, network.MasterConn.Conn)
	network.SendMessage(message, network.Broadcast.Conn)
}


func newOrderMaster(order driver.ButtonEvent){
	addToBackup(order)
	sendBackup()
}


func Cost (orderedFloor int, orderedDir string) (int){
	println("Cost()")
	qCopy = ko.Q
	var ordersInQ []int
	moreOrders := true
	cost := 100

	for ; moreOrders ; {
		temp := NextOrdered()
		if (temp == -1) {
			moreOrders = false
		}else{
			ordersInQ = append(ordersInQ,temp)
			qCopy.UP[temp] = 0
			qCopy.DOWN[temp] = 0
			qCopy.CMD[temp] = 0
		}
	}
	
	/*for i := 0; i < len(ordersInQ) ; i++ {
		println(ordersInQ[i])
	}*/

	for i := 0 ; i < len(ordersInQ) - 1 ; i++ {
		if (ordersInQ[i] == orderedFloor) && (i < cost){
			cost = i
		}else if (ordersInQ[i] < orderedFloor) && (ordersInQ[i] > orderedFloor) && (orderedDir == "UP") && (i < cost){
			cost = i
		}else if (ordersInQ[i] > orderedFloor) && (ordersInQ[i] < orderedFloor) && (orderedDir == "DOWN") && (i < cost){
			cost = i
		}
	}

	return cost	
}


func NextOrdered() (int){
	println("NextOrdered()")
	switch event.Dir {
	case "UP" :
		for i := event.Floor ; i < 4 ; i++ {
			if (qCopy.UP[i] == 1) || (qCopy.CMD[i] == 1) {
				return i
			}
		}
		
		for j := 3 ; j > event.Floor ; j-- {
			if (qCopy.DOWN[j] == 1) {
				return j
			}
		}

	case "DOWN" :
		for i := event.Floor ; i > -1 ; i-- {
			if (qCopy.DOWN[i] == 1) || (qCopy.CMD[i] == 1) {
				return i
			}
		}
	
		for j := 0 ; j < event.Floor ; j++ {
			if (qCopy.UP[j] == 1){
				return j
			}
		}
	}

	for i := 0 ; i < 4 ; i++{
		if (qCopy.UP[i] == 1) || (qCopy.DOWN[i] == 1) || (qCopy.CMD[i] == 1){
			return i
		}
	}

	return -1
}
		


func addToBackup(order driver.ButtonEvent){
	switch order.Button {
	case "U":
		backup.UP[order.Floor] = 1
	
	case "D":
		backup.DOWN[order.Floor] = 1

	case "C":
		backup.CMD[order.Floor] = 1
	}
}


func sendBackup(){
	var message string
	for i := 0 ; i < 4 ; i++ {
		message = message + strconv.Itoa(backup.UP[i])
	}
	for i := 0 ; i < 4 ; i++ {
		message = message + strconv.Itoa(backup.DOWN[i])
	}
	for i := 0 ; i < 4 ; i++ {
		message = message + strconv.Itoa(backup.CMD[i])
	}
	println(message)
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

func ReceiveOrder(){
	for  {
	message := <- network.OrderReceived
	floor, _ := strconv.Atoi(string(message.Message[3]))
	button, _ := strconv.Atoi(string(message.Message[2]))
	println("New order: ", button, floor)
	}
}
