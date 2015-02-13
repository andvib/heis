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
var PORT = "20015"

func initialize() {
	//Retrieve local IP address
	adr, _ := net.InterfaceAddrs()
	ip := strings.Split(adr[1].String(), "/")
	IP = ip[0]
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
	
	m := new(message)
	m.from = IP
	m.message = "c"

	sendMessage(m, conn)

	return conn
}

func receive(conn *net.UDPConn) (rec string){
	received := make([]byte,1024)
	//recAddr, _ := net.ResolveUDPAddr("udp", ":20015")
	//conn, _ := net.ListenUDP("udp", recAddr)
	_, _, _ = conn.ReadFromUDP(received)
	println(string(received))
	return string(received)
}

func send(s string, conn *net.UDPConn){
	b := []byte(s)
	_, _ = conn.Write(b)
}

func sendMessage(m *message, conn *net.UDPConn){
	messageString := m.from + "+" + m.message
	send(messageString, conn)
}

func receiveMessage(conn *net.UDPConn) (melding *message){
	mess := receive(conn)
	m := new(message)
	text := strings.Split(mess, "+")
	m.from = text[0]
	m.message = text[2]

	return m
}

func printMessage(conn *net.UDPConn){
	for ; true ; {
		m := receiveMessage(conn)
		println("From: ", m.from)
		println("Message: ", m.message)
		time.Sleep(500*time.Millisecond)
	}
}

func main(){
	initialize()
	
	runtime.GOMAXPROCS(runtime.NumCPU())
	recAddr, _ := net.ResolveUDPAddr("udp", ":20015")
	conn1, _ := net.ListenUDP("udp", recAddr)
	
	go printMessage(conn1)
	
	m := new(message)
	m.from = "Data1"
	m.message = "Heisann"
	
	//conn := connectMaster()
	/*for ; true ; {
		sendMessage(m, conn)
		time.Sleep(1000*time.Millisecond)
	}*/
}
