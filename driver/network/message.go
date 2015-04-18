/*Message handling. Contains the message struct that is
sent across network, and the functions needed to "code"
and "decode" the message struct.*/

package network

import ("strings"
		"net"
		"time")

type Message struct{
	From string
	Message string
}


type sentMessage struct {
	message *Message
	sent time.Time
}


var sentMessages []*sentMessage


func SendMessage(text string, conn *net.UDPConn){
	m := new(Message)
	m.From = IP
	m.Message = text
	messageString := m.From + "+" + m.Message
	send(messageString, conn)
}

func receiveMessage(mess string){
	//println("Receive message")
	m := new(Message)
	text := strings.Split(mess, "+")
    m.From = text[0]
    m.Message = text[1]
	//m.message = m.message[:2]
    whatToDo(m)
}

func printMessage(m *Message){
	println("From: ", m.From)
	println("Message: ", m.Message)
}


/*func addMessage(message *message) {
	var temp sentMessage
	temp.message = message
	temp.sent = time.Now()
	sentMessages = append(sentMessages,&temp)
}


func messageAcknowledged(message *message){
	for i := 0 ; i < len(sentMessages) ; i++ {
		if (sentMessages[i].message.message == message.message) && (sentMessages[i].message.from == message.from){
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
				SendMessage(sentMessages[i].message.message, findConn(sentMessages[i].message.from))
				removeMessage(i)
			}
		}
		time.Sleep(1000*time.Millisecond)
	}
}*/
