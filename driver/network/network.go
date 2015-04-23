/*"Main"-file of the network package. Contains global variables
and init function.*/

package network

import("net"
	   "strings"
	   "time"
       "strconv")


type Connection struct{
	IP string
	Conn *net.UDPConn
	LastSignal time.Time
}

var NewMessage = make(chan *Message)

var Connected []Connection
var Master = false
var IP string
var PORT = "30020"
var MasterConn Connection
var Broadcast Connection


func NETWORK_init(){
    println("NETWORK_init()")
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
	MasterConn.LastSignal = time.Now()
	var temp time.Time
	temp = time.Now()

	go timeout()

	for ; (time.Since(temp) < 350*time.Millisecond) && (Master == false) ; {
	}

    println("Master: ", Master)

	conn := connect("")

	//Asks other units to connect
	//SendMessage("nw", Broadcast.Conn)

	go alive(conn)
	go updateMessages()
}


func alive(conn *net.UDPConn){
	for ; true ; {
		if Master {
    		SendMessage("am", conn,false)
    		time.Sleep(100*time.Millisecond)
		}else{
			SendMessage("as", conn,false)
			time.Sleep(200*time.Millisecond)
		}
	}
}


func whatToDo(m *Message){
	//Checks what to do with the new message
	order := m.Message[:2]

	if (m.From == IP){
		return
	}

	if order == "am"{
		//Alive-signal from master
		MasterConn.LastSignal = time.Now()
		//println(Master)
	    if /*(MasterConn.IP != m.From)*/(Master){
			println("Other master")
			Master = false
            MasterConn.IP = m.From
		    MasterConn.LastSignal = time.Now()
            for i := 0 ; i < len(Connected) ; i++ {
                if m.From == Connected[i].IP {
                    MasterConn.Conn = Connected[i].Conn
					RemoveConn(i)
                }
            }
			//WhosMaster()
		}

    }else if order == "as" {
		//println("Alive slave")
		found := false
        for i := 0 ; i < len(Connected) ; i++ {
            if m.From == Connected[i].IP {
                Connected[i].LastSignal = time.Now()
				found = true
            }
        }

		if (!found) && (m.From != MasterConn.IP) {
			connect(m.From)
		}

	}else if order == "ac" {
		messageAcknowledged(m)

    /*}else if order == "nw" {
		
        println("New connection from: ", m.From)
        conn := connect(m.From)
        SendMessage("co", conn)
    }else if order == "co" {	
        println("Connect to: ", m.From)
        connect(m.From)*/
    }else{
		NewMessage <- m
	}
}


func timeout(){
	for ; true ; {
		//println("Timeout")
		if (time.Since(MasterConn.LastSignal) > 900*time.Millisecond) && (Master == false){
			//No master on the network
            //println("Master timeout")
			go WhosMaster()
			time.Sleep(500*time.Millisecond)
		}

		for i := 0 ; i < len(Connected) ; i++ {
			//println("Check slave timeout")
			if (time.Since(Connected[i].LastSignal) > 1200*time.Millisecond) {
                println("Slave timeout: ", Connected[i].IP)
				Connected[i].Conn.Close()
				RemoveConn(i)
				if (Master) {
					var temp Message
					temp.From = IP
					temp.Message = "nm"
					NewMessage <- &temp
				}
			}
		}
	}
}


func AppendConn(conn *net.UDPConn, ip string){
    var temp Connection
    temp.IP = ip
    temp.Conn = conn
    temp.LastSignal = time.Now()
    Connected = append(Connected, temp)
}


func RemoveConn(index int) {
	Connected = append(Connected[:index], Connected[index+1:]...)
}


func FindConn(ip string) (*net.UDPConn){
	for i := 0 ; i < len(Connected) ; i++ {
		if Connected[i].IP == ip {
			return Connected[i].Conn
		}
	}
	return nil
}


func WhosMaster() {
	println("Master timeout : WhosMaster()")
    me := true
    number := len(IP)
    own, _ := strconv.Atoi(string(IP[number-2:number]))

	for i := 0 ; i < len(Connected) ; i++ {
		println("Connected[",i,"] : ", Connected[i].IP)
	}

    for i := 0 ; i < len(Connected) ; i++ {
        other, _ := strconv.Atoi(string(Connected[i].IP[number-2:number]))

		if other < own{
			me = false
			println("I am not the new master")
		}
    }
    
    if (me == true) {
        println("I am the new master")
        Master = true
		
		var temp Message
		temp.From = IP
		temp.Message = "nm"
		NewMessage <- &temp
    }

    MasterConn.IP = ""
    //MasterConn.Conn.Close()
	MasterConn.Conn = nil
}
