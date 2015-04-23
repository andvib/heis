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
	message *Message
	sent time.Time
}


var sentMessages []*sentMessage


func SendMessage(text string, conn *net.UDPConn, ack bool){
	m := new(Message)
	m.From = IP
	m.Message = text
	messageString := m.From + "+" + m.Message
	
	if(ack) {
		m.To = MasterConn.IP
		addMessage(m, conn)
	}

	send(messageString, conn)
}

func receiveMessage(mess string){
	//println("Receive message")
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
	temp.message = message
	temp.sent = time.Now()
	temp.To = conn
	sentMessages = append(sentMessages,&temp)
}


func messageAcknowledged(message *Message){
	for i := 0 ; i < len(sentMessages) ; i++ {
		if (sentMessages[i].message.Message[2:] == message.Message) && (sentMessages[i].message.From == message.From){
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
			if (time.Since(sentMessages[i].sent) > 2000*time.Millisecond){
				SendMessage(sentMessages[i].message.Message, FindConn(sentMessages[i].message.From),true)
				removeMessage(i)
			}
		}
		time.Sleep(1000*time.Millisecond)
	}
}
