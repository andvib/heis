package phoenix

import("time"
		"os/exec"
		"net"
		"encoding/json")


var Master int
var LastSignal time.Time


func Phoenix(){
	Master = 0
	LastSignal = time.Now()

	slave()
	println("Exiting phoenix")
}


func spawn(){
	backup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run /home/haavardo/heis/main.go")
	backup.Start()
}


func connect() (*net.UDPConn){
	sendAddr, _ := net.ResolveUDPAddr("udp", "localhost:30011")
	conn, _ := net.DialUDP("udp", nil, sendAddr)
	return conn
}


func connReceive() (*net.UDPConn){
	recAddr, _ := net.ResolveUDPAddr("udp", ":30011")
	recConn, _ := net.ListenUDP("udp", recAddr)
	return recConn
}


func receive(conn *net.UDPConn){
	received := make([]byte, 1)
	for ; Master == 0 ; {
		_, _, _ = conn.ReadFromUDP(received)
		//_ = json.Unmarshal(received, &Counter)
		LastSignal = time.Now()
	}
	conn.Close()
}


func send(conn *net.UDPConn){
	a, _ := json.Marshal("a")
	_, _ = conn.Write(a)
}

func alive(){
	conn := connect()
	for {
		send(conn)
		time.Sleep(1000*time.Millisecond)
	}
}


func slave(){
	recConn := connReceive()
	go receive(recConn)
	for {
		if time.Since(LastSignal) > 2000*time.Millisecond{
			Master = 1
			go alive()
			spawn()
			return
		}
	}
}
