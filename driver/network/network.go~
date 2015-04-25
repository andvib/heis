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


var IP string
var PORT = "30021"
var Connected []Connection
var MasterConn Connection
var Broadcast Connection
var Master = false
var Alive bool


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

	//Establish broadcast-connection
	conn := connect("")

	Alive = true
	go alive(conn)
	go updateMessages()
}


func alive(conn *net.UDPConn){
	//Sends appropriate alive-message
	for ; Alive ; {
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

	    if (Master){
			//Detected another master-unit on the network, switches to slave
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
		}

    }else if order == "as" {
		//Alive signal from slave
		found := false
        for i := 0 ; i < len(Connected) ; i++ {
            if m.From == Connected[i].IP {
                Connected[i].LastSignal = time.Now()
				found = true
            }
        }

		if (!found) && (m.From != MasterConn.IP) {
			//Connects to the unit if not already connected
			connect(m.From)
		}

	}else if order == "ac" {
		//Previous message acknowledged
		messageAcknowledged(m)

    }else{
		//Sends message to message-handling in ****
		NewMessage <- m
	}
}


func timeout(){
	for ; true ; {
		if (time.Since(MasterConn.LastSignal) > 900*time.Millisecond) && (Master == false){
			//No master on the network
			go WhosMaster()
			time.Sleep(500*time.Millisecond)
		}

		for i := 0 ; i < len(Connected) ; i++ {
			if (time.Since(Connected[i].LastSignal) > 1200*time.Millisecond) {
				//Slave timeout
                println("Slave timeout: ", Connected[i].IP)
				Connected[i].Conn.Close()
				RemoveConn(i)

				if (Master) {
					//Tells ****-module to distribute backup-orders
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
	//Function for deciding who's the new master
	println("Master timeout : WhosMaster()")
    me := true
    number := len(IP)
    own, _ := strconv.Atoi(string(IP[number-2:number]))

	for i := 0 ; i < len(Connected) ; i++ {
		println("Connected[",i,"] : ", Connected[i].IP)
	}

	//Compares IP, lowest IP is the new master
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
		
		//Distributes backup
		var temp Message
		temp.From = IP
		temp.Message = "nm"
		NewMessage <- &temp
    }

    MasterConn.IP = ""
	MasterConn.Conn = nil
}
