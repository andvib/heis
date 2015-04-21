package costfunc

import(.".././network"
		".././heis/"
		"strconv"
		".././ko/"
		".././event/"
		"time"
		"net")

var backup ko.Queue
var qCopy ko.Queue
var BestOrder OrderCost

type OrderCost struct{
	Cost int
	Conn *net.UDPConn
}


func ReceiveMessage(){
	var received *Message

	for ;; {
		received = <- NewMessage
		
		order := received.Message[:2]

		if (order == "oc"){
			var temp OrderCost
			temp.Cost, _ = strconv.Atoi(string(received.Message[2]))
			temp.Conn = FindConn(received.From)
			CheckCost(temp)
		}else if (order == "no"){
			var temp driver.ButtonEvent
			temp.Floor, _ = strconv.Atoi(string(received.Message[3]))
			temp.Button = string(received.Message[2])
			newOrderMaster(temp)
		}else if (order == "cc"){
			var temp driver.ButtonEvent
			temp.Floor, _ = strconv.Atoi(string(received.Message[2]))
			temp.Button = string(received.Message[3])
			lightsOn(temp)	
			slaveCalculate(temp)
		}else if (order == "eo"){
			floor, _ := strconv.Atoi(string(received.Message[2]))
			ko.AddOrder(floor, string(received.Message[3]))
		}else if (order == "ba"){
			receiveBackup(received.Message[2:])
		}else if (order == "nm"){
			newMaster()
		}else if (order == "rm"){
			if (Master){
				RemoveBUOrder(received)
			}
			lightsOff(int(received.Message[2]))
		}
	}
}


func ButtonHandle(){
	var buttonEvent driver.ButtonEvent

	for ; true ; {
		buttonEvent = <- driver.ButtonChan
		if (buttonEvent.Button == "C"){
			ko.AddOrder(buttonEvent.Floor,buttonEvent.Button)
		}else if Master {	
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
	SendMessage(message, Broadcast.Conn)
	println("New order message sent")
}


func newOrderMaster(order driver.ButtonEvent){
	println("NewOrderMaster")
	addToBackup(order)
	sendBackup()
	
	floor := strconv.Itoa(order.Floor)

	lightsOn(order)

	message := "cc" + floor + order.Button
	SendMessage(message,Broadcast.Conn)
	println("Sent calculate cost")
	timer := time.Now()
	//var bestOrder network.OrderCost
	BestOrder.Cost = 1000
	for ; (time.Since(timer) < 200*time.Millisecond) ; {}

	if Cost(order.Floor,order.Button) <= BestOrder.Cost{
		println("Takes the order itself")
		ko.AddOrder(order.Floor,order.Button)
	}else{
		println("Sends execute order")
		newOrder := "eo" + strconv.Itoa(order.Floor) + order.Button
		SendMessage(newOrder,BestOrder.Conn)
	}
}


func CheckCost(c OrderCost){
	if c.Cost < BestOrder.Cost{
		BestOrder.Cost = c.Cost
		BestOrder.Conn = c.Conn
	}
}


func slaveCalculate(order driver.ButtonEvent){
	cost := Cost(order.Floor,order.Button)
	message := "oc" + strconv.Itoa(cost)
	SendMessage(message,Broadcast.Conn)
}


func Cost (orderedFloor int, orderedDir string) (int){
	println(orderedDir)
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
	
	for i := 0; i < len(ordersInQ) ; i++ {
		println(ordersInQ[i])
	}


	if (ko.EmptyQ() == 1) {
		println("Empty Q")
		cost = 1
	}

	for i := 0 ; i < len(ordersInQ) ; i++ {
		if (ordersInQ[i] == orderedFloor){
			println("Etasje allerede bestilt!")
			cost = 0
		}
	}
	
	if (len(ordersInQ) != 0){
		if (orderedFloor < ordersInQ[0]) && (orderedDir == "U") && (event.Floor < orderedFloor){
			println("Forst i koen, opp")
			cost = 0
		}else if (orderedFloor > ordersInQ[0]) && (orderedDir == "D") && (event.Floor > orderedFloor){
			println("Forst i koen, ned")
			cost = 0
		}
	}

	for i := 0 ; i < len(ordersInQ) -1 ; i++{
		if (ordersInQ[i] < orderedFloor) && (ordersInQ[i+1] > orderedFloor) && (orderedDir == "U") && (2*(i+1) < cost){
			cost = 2*(i+1)
		}else if (ordersInQ[i] > orderedFloor) && (ordersInQ[i+1] < orderedFloor) && (orderedDir == "D") && (2*(i+1) < cost){
			cost = 2*(i+1)
		}
	}

	if (event.Floor == orderedFloor){
		cost = 0
	}

	if ((len(ordersInQ) + 1)*2 < cost){
		println("Last in Q")
		cost = (len(ordersInQ)+1)*2
	}
	

	println("Cost: ", cost)
	return cost	
}


func NextOrdered() (int){
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
	SendMessage(message, Broadcast.Conn)
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


func newMaster(){
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
		}
	}
}


func RemoveBUOrder(m *Message){
	floor,_ := strconv.Atoi(string(m.Message[2]))
	backup.UP[floor] = 0
	backup.DOWN[floor] = 0
	backup.CMD[floor] = 0


	/*for i := 3 ; i > -1 ; i-- {
		println(backup.UP[i], "\t", backup.DOWN[i], "\t", backup.CMD[i])
	}
	println("")*/
}


func lightsOff(m int){
	floor,_ := strconv.Atoi(string(m))
	println(floor)
	driver.ELEV_set_button_lamp(0,floor,0)
	driver.ELEV_set_button_lamp(1,floor,0)
}


func lightsOn(m driver.ButtonEvent){
	switch m.Button{
	case "U":
		driver.ELEV_set_button_lamp(0,m.Floor,1)
	
	case "D":
		driver.ELEV_set_button_lamp(1,m.Floor,1)
	}
}
