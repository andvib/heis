/*"Main"-file of the network package. Contains global variables
and init function.*/

package network

import("net"
	   "strings"
	   "time"
       "strconv"
	   /*"runtime"*/)


type Connection struct{
	IP string
	Conn *net.UDPConn
	LastSignal time.Time
}


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
	SendMessage("nw", Broadcast.Conn)

	go alive(conn)
	//go updateMessages()
}


func alive(conn *net.UDPConn){
	for ; true ; {
		if Master {
    		SendMessage("am", conn)
    		time.Sleep(100*time.Millisecond)
		}else{
			SendMessage("as", conn)
			time.Sleep(1000*time.Millisecond)
		}
	}
}


func whatToDo(m *message){
	//Checks what to do with the new message
    //println("Whattodo")
	if m.message == "am"{
        //println("Alive signal from master")
		//Alive-signal from master
		MasterConn.LastSignal = time.Now()
	    if MasterConn.IP == ""{
			//masterConn.Conn = connect(m.from)
            MasterConn.IP = m.from
		    MasterConn.LastSignal = time.Now()
            for i := 0 ; i < len(Connected) ; i++ {
                if m.from == Connected[i].IP {
                    MasterConn.Conn = Connected[i].Conn
                }
            }
		}

    }else if m.message == "as" {
	//println("Alive signal from slave")
        for i := 0 ; i < len(Connected) ; i++ {
            if m.from == Connected[i].IP {
                Connected[i].LastSignal = time.Now()
            }
        }
    }else if m.message == "nw" {
		if (m.from == IP){
			return
		}
        println("New connection from: ", m.from)
        conn := connect(m.from)
        SendMessage("co", conn)
    }else if m.message == "co" {
        println("Connect to: ", m.from)
        connect(m.from)
    }        
}


func timeout(){
	for ; true ; {		
		if (time.Since(MasterConn.LastSignal) > 200*time.Millisecond) && (Master == false){
			//No master on the network
            println("Master timeout")
			WhosMaster()
			time.Sleep(500*time.Millisecond)
			
		}

		for i := 0 ; i < len(Connected) ; i++ {
			if (time.Since(Connected[i].LastSignal) > 1200*time.Millisecond) {
                println("Slave timeout: ", Connected[i].IP)
				RemoveConn(i)
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


func findConn(ip string) (*net.UDPConn){
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
    own_1, _ := strconv.Atoi(string(IP[number-1]))
    own_2, _ := strconv.Atoi(string(IP[number-2]))

	for i := 0 ; i < len(Connected) ; i++ {
		println("Connected[",i,"]",Connected[i].IP)
	}


    for i := 0 ; i < len(Connected) ; i++ {
        other_1, _ := strconv.Atoi(string(Connected[i].IP[number-1]))
        other_2, _ := strconv.Atoi(string(Connected[i].IP[number-2]))

        if (other_1 <= own_1) {
            if (other_2 <= own_2) {
                me = false
				println("I am not the new master")
            }
        }
    }
    
    if (me == true) {
        println("I am the new master")
        Master = true
    }

    MasterConn.IP = ""
    MasterConn.Conn = nil
}
