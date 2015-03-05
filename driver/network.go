package main

import("net"
	   "strings"
	   "time"
	   "runtime")

type message struct{
	from string
	message string
}

var connected = make(map[string]*net.UDPConn)
var Master = false
var IP string
var PORT = "30011"
var LastSignal time.Time
var MasterConn *net.UDPConn

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

	conn := connect("")

	//Asks other units to connect
	send("n:"+IP, connected["bc"])
	
	if Master == true{
		go alive(conn)
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

	//Creates connection 
	sendAddr, err := net.ResolveUDPAddr("udp", address)
	conn, err2 := net.DialUDP("udp", nil, sendAddr)

	if (err != nil || err2 != nil){
		println("Kunne ikke koble til!")
		return nil
	}

	//Adds the new connection to the map
	if ip == ""{
		connected["bc"] = conn
	}else{
		connected[ip] = conn
	}
	
	//Sends a message to the new connection that they should connect as well
	text := "c"

	sendMessage(text, conn)

	return conn
}


func receive(conn *net.UDPConn){
	//Receive message from network
	received := make([]byte,1024)	
	for ; true ; {
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
	}
}


func slave() {
	for ; Master == false ; {
		//timeout()
	}
}


func whatToDo(m *message){
	//Checks what to do with the new message
	if m.message == "a"{
		LastSignal = time.Now()

		if connected["Master"] == nil{
			connect(m.from)
			MasterConn = connected[m.from]
		}

	}else if m.message == "b"{
		println("Ny bestilling")

	}else if m.message == "c"{
		shouldConnect(m.from)

	}else if m.message[0] == 110{
		ip := strings.Split(m.message, ":")
		shouldConnect(ip[1])
	}
}


func shouldConnect(addr string){
	// Iterate throught connected-map to check if the IP is allready connected

	for key, _ := range connected{
		if key == addr{
			//Returns from function if the connection already exists
			return;
		}
	}

	connect(addr)
}

	
func timeout(){
	//Checks for alive signal from master
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
}
