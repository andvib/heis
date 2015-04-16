package costfunc

import(".././network"
	".././heis"
	"strconv")


func ButtonHandle(){
	var buttonEvent driver.ButtonEvent

	for ; true ; {
		buttonEvent = <- driver.ButtonChan
		println("New order: ", buttonEvent.Floor, ", ",buttonEvent.Button)
		newOrderSlave(buttonEvent)
	}	
}


func newOrderSlave(order driver.ButtonEvent){
	message := order.Button + strconv.Itoa(order.Floor)
	println(message)
	network.SendMessage(message, network.MasterConn.Conn)

}
