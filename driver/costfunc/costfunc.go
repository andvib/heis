package elevlog

import(.".././network"
		".././heis/"
		"strconv"
		".././queue/"
		".././event/"
		"time"
		"net")

var backup queue.Queue
var qCopy queue.Queue
var BestOrder OrderCost

type OrderCost struct{
	Cost int
	Conn *net.UDPConn
}


func ButtonHandle(){
	//Decides what to be done when a button is pushed
	var buttonEvent driver.ButtonEvent

	for ; true ; {
		buttonEvent = <- driver.ButtonChan
		if (buttonEvent.Button == "C"){
			//Internal order, added to local q
			queue.AddOrder(buttonEvent.Floor,buttonEvent.Button)

		}else if Master {	
			newOrderMaster(buttonEvent)

		}else{
			newOrderSlave(buttonEvent)
		}
	}	
}


func newOrderSlave(order driver.ButtonEvent){
	//Sends message to master containing the new order
	message := "no" + order.Button + strconv.Itoa(order.Floor)
	println(message)
	SendMessage(message, Broadcast.Conn, true)
}


func newOrderMaster(order driver.ButtonEvent){
	//New order received by master
	addToBackup(order)
	sendBackup()
	
	floor := strconv.Itoa(order.Floor)

	lightsOn(order)

	//Sends calculate cost message to slaves
	message := "cc" + floor + order.Button
	SendMessage(message,Broadcast.Conn,false)
	
	//Waits to receive costs from slaves
	timer := time.Now()
	BestOrder.Cost = 1000
	for ; (time.Since(timer) < 1000*time.Millisecond) ; {}


	if Cost(order.Floor,order.Button) <= BestOrder.Cost{
		println("Takes the order itself")
		queue.AddOrder(order.Floor,order.Button)
	}else{
		println("Sends execute order")
		newOrder := "eo" + strconv.Itoa(order.Floor) + order.Button
		SendMessage(newOrder,BestOrder.Conn,false)
	}
}


func CheckCost(c OrderCost){
	//Checks if cost received by slave is better than current cost
	if c.Cost < BestOrder.Cost{
		BestOrder.Cost = c.Cost
		BestOrder.Conn = c.Conn
		println("Better cost: ", c.Cost)
	}
}


func slaveCalculate(order driver.ButtonEvent){
	//Slave calculates its own cost and sends result to master
	cost := Cost(order.Floor,order.Button)
	message := "oc" + strconv.Itoa(cost)
	SendMessage(message,Broadcast.Conn,false)
}


func Cost (orderedFloor int, orderedDir string) (int){
	//Calculates cost for executing order
	qCopy = queue.Q
	var ordersInQ []int
	moreOrders := true
	cost := 100

	//Creates an array containing the order that the local orders will be executed
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

	//FINDS THE CURRENT COST OF THE ORDER
	//Checks if the queue is empty
	if (queue.EmptyQ() == 1) {
		cost = 1
	}
	
	//Checks if the floor is already ordered
	for i := 0 ; i < len(ordersInQ) ; i++ {
		if (ordersInQ[i] == orderedFloor){
			cost = i+1
		}
	}
	
	//Checks if the order can be put ahead of existing orders	
	if (len(ordersInQ) != 0){
		if (orderedFloor < ordersInQ[0]) && (orderedDir == "U") && (event.Floor < orderedFloor){
			cost = 2
		}else if (orderedFloor > ordersInQ[0]) && (orderedDir == "D") && (event.Floor > orderedFloor){
			cost = 2
		}
	}

	//Checks if the order can be put between to existing orders in the given direction
	for i := 0 ; i < len(ordersInQ) -1 ; i++{
		if (ordersInQ[i] < orderedFloor) && (ordersInQ[i+1] > orderedFloor) && (orderedDir == "U") && (3*(i+1) < cost){
			cost = 3*(i+1)
		}else if (ordersInQ[i] > orderedFloor) && (ordersInQ[i+1] < orderedFloor) && (orderedDir == "D") && (3*(i+1) < cost){
			cost = 3*(i+1)
		}
	}

	//Checks if the elevator is currently IDLE at ordered floor
	if (event.Floor == orderedFloor) && (event.State != "MOVING"){
		cost = 0
	}

	//If the order can't be put between existing orders, how much is the cost of putting it at the end
	if ((len(ordersInQ) + 1)*2 < cost){
		cost = (len(ordersInQ)+1)*3
	}
	

	println("Cost: ", cost)
	return cost	
}


func NextOrdered() (int){
	//Function used to find the next order in the local queue. Same logic as in the queue-module
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


func RemoveBUOrder(m *Message){
	floor,_ := strconv.Atoi(string(m.Message[2]))
	backup.UP[floor] = 0
	backup.DOWN[floor] = 0
	backup.CMD[floor] = 0
}


func lightsOff(m int){
	floor,_ := strconv.Atoi(string(m))
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
