package costfunc

import(".././network"
		".././heis/"
		"strconv"
		".././ko/"
		".././event/"
		"time")

var backup ko.Queue
var qCopy ko.Queue

var BestOrder network.OrderCost

func ButtonHandle(){
	var buttonEvent driver.ButtonEvent

	for ; true ; {
		buttonEvent = <- driver.ButtonChan
		if network.Master {	
			newOrderMaster(buttonEvent)
		}else{
			newOrderSlave(buttonEvent)
		}
	}	
}


func newOrderSlave(order driver.ButtonEvent){
	println("NewOrderSlave")
	message := "no" + order.Button + strconv.Itoa(order.Floor)
	println(message)
	network.SendMessage(message, network.Broadcast.Conn)
	println("New order message sent")
}


func newOrderMaster(order driver.ButtonEvent){
	println("NewOrderMaster")
	addToBackup(order)
	sendBackup()
	
	floor := strconv.Itoa(order.Floor)

	message := "cc" + floor + order.Button
	network.SendMessage(message,network.Broadcast.Conn)
	println("Sent calculate cost")
	timer := time.Now()
	//var bestOrder network.OrderCost
	BestOrder.Cost = 1000
	for ; (time.Since(timer) < 500*time.Millisecond) ; {}
	
	/*for o := range network.CostReceived{
		println("Something on the channel")
		if o.Cost < bestOrder.Cost {
			bestOrder.Cost = o.Cost
			bestOrder.Conn = o.Conn
		}
	}*/

	if Cost(order.Floor,order.Button) < BestOrder.Cost{
		println("Takes the order itself")
		ko.AddOrder(order.Floor,order.Button)
	}else{
		println("Sends execute order")
		newOrder := "eo" + strconv.Itoa(order.Floor) + order.Button
		network.SendMessage(newOrder,BestOrder.Conn)
	}
}


func CostChan(){
	for ;; {
		o := <-network.CostReceived
		if o.Cost < BestOrder.Cost{
			BestOrder.Cost = o.Cost
			BestOrder.Conn = o.Conn
		}
	}
}


func SlaveCalculate(){
	var temp driver.ButtonEvent
	for ;; {
		temp = <- network.CalCost
		println("Calculating cost")
		//floor, _ := strconv.Atoi(string(m.Message[2]))
		cost := Cost(temp.Floor,temp.Button)
		message := "oc" + strconv.Itoa(cost)
		network.SendMessage(message,network.Broadcast.Conn)
	}
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
			cost = 2*i  //0?
		}else if (ordersInQ[i] < orderedFloor) && (ordersInQ[i+1] > orderedFloor) && (orderedDir == "UP") && (i < cost){
			cost = 2*i
		}else if (ordersInQ[i] > orderedFloor) && (ordersInQ[i+1] < orderedFloor) && (orderedDir == "DOWN") && (i < cost){
			cost = 2*i
		}
	}

	cost = cost + event.Floor

	if (ko.EmptyQ() == 1) {
		println("Empty Q")
		cost = 1
	}

	println("Cost: ", cost)
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
	//println(message)
	network.SendMessage(message, network.Broadcast.Conn)
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
	var temp driver.ButtonEvent
	for  {
	message := <- network.OrderReceived
	temp.Floor, _ = strconv.Atoi(string(message.Message[3]))
	temp.Button = string(message.Message[2])
	println("New order: ", temp.Button, temp.Floor)
	//ko.AddOrder(floor,button)
	//Cost(floor,button)
	newOrderMaster(temp)
	}
}
