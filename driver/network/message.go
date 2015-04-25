/*Message handling. Contains the message struct that is
sent across network, and the functions needed to "code"
and "decode" the message struct.*/

package network

import ("strings"
		"net"
		"time")

type Message struct{
	From string
	To string
	Message string
}


type sentMessage struct {
	To *net.UDPConn
	messageSent *Message
	timeSent time.Time
}


var sentMessages []*sentMessage


func SendMessage(text string, conn *net.UDPConn, acknowledge bool){
	m := new(Message)
	m.From = IP
	m.Message = text
	messageString := m.From + "+" + m.Message
	
	if(acknowledge) {
		m.To = MasterConn.IP
		addMessage(m, conn)
	}

	send(messageString, conn)
}


func receiveMessage(mess string){
	m := new(Message)
	text := strings.Split(mess, "+")
    m.From = text[0]
    m.Message = text[1]
	
	if m.From != IP{
		whatToDo(m)
	}
}

func printMessage(m *Message){
	println("")
	println("From: ", m.From)
	println("Message: ", m.Message)
	println("")
}


func addMessage(message *Message, conn *net.UDPConn) {
	var temp sentMessage
	temp.messageSent = message
	temp.timeSent = time.Now()
	temp.To = conn

	sentMessages = append(sentMessages,&temp)
}


func messageAcknowledged(message *Message){
	for i := 0 ; i < len(sentMessages) ; i++ {
		if (sentMessages[i].messageSent.Message == message.Message[2:6]){
			removeMessage(i)
		}
	}
}


func removeMessage(i int){
	sentMessages = append(sentMessages[:i], sentMessages[i+1:]...)
}


func updateMessages(){
	for ; true ; {
		for i := 0 ; i < len(sentMessages) ; i++ {
			if (time.Since(sentMessages[i].timeSent) > 2000*time.Millisecond){
				SendMessage(sentMessages[i].messageSent.Message, sentMessages[i].To,true)
				removeMessage(i)
			}
		}
		print(len(sentMessages))
		time.Sleep(1000*time.Millisecond)
	}
}
