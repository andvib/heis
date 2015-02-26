package main

import("net"
	   "strings"
	   "time"
	   "runtime")

type message struct{
	from string
	message string
}

var connected = []*net.UDPConn{}
var Master = false
var IP string
var PORT = "30011"
var LastSignal time.Time

func initialize(){
	//Retrieve local IP address
	adr, _ := net.InterfaceAddrs()
	ip := strings.Split(adr[1].String(), "/")
	IP = ip[0]

	//Connect as listener
	recAddr, _ := net.ResolveUDPAddr("udp",":" + PORT)
	recConn, _ := net.ListenUDP("udp", recAddr)

	//Start listening-thread
	go receive(recConn)

	//Check for master on network
	LastSignal = time.Now()
	
	var temp time.Time
	temp = time.Now()
	
	for ; time.Since(temp) < 350*time.Millisecond && Master == false; {
		timeout()
	}
	//recConn.Close()
	conn := connect("")

	if Master == true{
		go alive(conn)
		//recConn.Close()
	}else{
		go slave()
		go receive(recConn)
	}
}

func connect(ip string) (connection *net.UDPConn){
	var address string

	if ip == ""{
		address = "129.241.187.255:" + PORT
	}else{
		address = ip + ":" + PORT
	}

	sendAddr, err := net.ResolveUDPAddr("udp", address)
	conn, err2 := net.DialUDP("udp", nil, sendAddr)

	if (err != nil || err2 != nil){
		println("Kunne ikke koble til!")
		return nil
	}
	
	/*m := new(message)
	m.from = IP
	m.message = "c"

	sendMessage(m, conn)*/

	return conn
}

func receive(conn *net.UDPConn){
	//Receive message from network
	received := make([]byte,1024)	
	for ; Master == false ; {
		_, _, _ = conn.ReadFromUDP(received)
		println("Motatt: ", string(received))
		LastSignal = time.Now()
	}
	conn.Close()
	println("Stenger mottak")
	//Send received message to message-handling
	//go receiveMessage(string(received))
}

func send(s string, conn *net.UDPConn){
	b := []byte(s)
	_, _ = conn.Write(b)
}

func sendMessage(m *message, conn *net.UDPConn){
	messageString := m.from + "+" + m.message
	send(messageString, conn)
}

func receiveMessage(mess string){
	m := new(message)
	text := strings.Split(mess, "+")
	
	m.from = text[0]
	m.message = text[2]
	
	whatToDo(m)
}

func printMessage(m *message){
	println("From: ", m.from)
	println("Message: ", m.message)
}

func alive(conn *net.UDPConn){
	for ; Master == true ; {
		send("a", conn)
		time.Sleep(100*time.Millisecond)
		println("Alive")
	}
}

func slave() {
	for ; Master == false ; {
		//timeout()
	}
}	

func whatToDo(m *message){
	if m.message == "a"{
		LastSignal = time.Now()
	}else if m.message == "b"{
		println("Ny bestilling")
	}else if m.message == "c"{
		connect(m.from)
	}
}

func timeout(){
	if time.Since(LastSignal) > 200*time.Millisecond {
		//No master on the network
		Master = true
		println("Tar over som master")
		//go alive()
	}
}

func main(){
	initialize()
	
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	for ; true ; {
		if Master == true{
		}
	}

	
	/*m := new(message)
	m.from = "Data1"
	m.message = "Heisann"*/
	
	//conn := connectMaster()
	/*for ; true ; {
		sendMessage(m, conn)
		time.Sleep(1000*time.Millisecond)
	}*/
}
