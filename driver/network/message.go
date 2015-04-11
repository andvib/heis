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

func sendMessage(text string, conn *net.UDPConn){
  m := new(message)
  m.from = IP
  m.message = text
  messageString := m.from + "+" + m.message
  send(messageString, conn)
}

func receiveMessage(mess string){
	m := new(message)
	text := strings.Split(mess, "+")
    m.from = text[0]
    m.message = text[1]
    whatToDo(m)
}

func printMessage(m *message){
  println("From: ", m.from)
  println("Message: ", m.message)
}