/*"Main"-file of the network package. Contains global variables
and init function.*/

package network

import("net"
	"strings"
	"time"
	/*"runtime"*/)

type connection struct{
	ip string
	conn *net.UDPConn
	lastSignal time.Time
}


var connected []connection
var Master = false
var IP string
var PORT = "30020"
var masterConn connection
var broadcast connection


func NETWORK_init(){
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
	masterConn.lastSignal = time.Now()
	var temp time.Time
	temp = time.Now()

	go timeout()

	for ; (time.Since(temp) < 350*time.Millisecond) && (Master == false) ; {
	}

	conn := connect("")

	//Asks other units to connect
	sendMessage("n", broadcast.conn)

	go alive(conn)

	if Master == false{
		go slave()
	}
	println(Master)
}


func alive(conn *net.UDPConn){
	for ; true ; {
		if Master {
    		sendMessage("am", conn)
			println("Alive")
    		time.Sleep(100*time.Millisecond)
		}else{
			sendMessage("as", conn)
			time.Sleep(1000*time.Millisecond)
		}
	}
}

func slave() {
	for ; Master == false ; {
	    timeout()
		time.Sleep(100*time.Millisecond)
	}
}

func whatToDo(m *message){
	//Checks what to do with the new message
	if m.message == "am"{
		//Alive-signal from master
		masterConn.lastSignal = time.Now()
	    /*if MasterConn == nil{
			connect(m.from)
			MasterConn = connected[m.from]
		}*/
	}/*else if m.message == "as"{
		//Alive-signal from slave
		
	}else if m.message == "b"{
		println("Ny bestilling")
	}else if m.message == "c"{
		shouldConnect(m.from)
	}else if m.message[0] == 110{
		ip := strings.Split(m.message, ":")
		shouldConnect(ip[1])
  	}*/
}

func shouldConnect(addr string){
  // Iterate throught connected-map to check if the IP is allready connected
  for i := 0 ; i < len(connected) ; i++ {
    if connected[i].ip == addr{
      //Returns from function if the connection already exists
      return;
    }
  }
  connect(addr)
}

func timeout(){
	for ; true ; {		
		if (time.Since(masterConn.lastSignal) > 200*time.Millisecond) && (Master == false){
			//No master on the network
			Master = true
			println("Tar over som master")
		}

		for i := 0 ; i < len(connected) ; i++ {
			if (time.Since(connected[i].lastSignal) > 1200*time.Millisecond) {
				//Dead
			}
		}
	}
}
	


	/*//Checks for alive signal from master
	if time.Since(LastSignal) > 200*time.Millisecond {
		//No master on the network
		Master = true
		println("Tar over som master")
		//connected["Master"] = nil
	}
}*/
