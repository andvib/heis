/*Message handling. Contains the message struct that is
sent across network, and the functions needed to "code"
and "decode" the message struct.*/

package network

import ("strings"
		"net")

type message struct{
	from string
	message string
}

func SendMessage(text string, conn *net.UDPConn){
	m := new(message)
	m.from = IP
	m.message = text
	messageString := m.from + "+" + m.message
	send(messageString, conn)
}

func receiveMessage(mess string){
	//println("Receive message")
	m := new(message)
	text := strings.Split(mess, "+")
    m.from = text[0]
    m.message = text[1]
	m.message = m.message[:2]
    whatToDo(m)
}

func printMessage(m *message){
	println("From: ", m.from)
	println("Message: ", m.message)
}
