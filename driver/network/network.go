/*"Main"-file of the network package. Contains global variables
and init function.*/

package network

import("net"
	"strings"
	"time"
	/*"runtime"*/)




var connected = make(map[string]*net.UDPConn)
var Master = false
var IP string
var PORT = "30020"
var LastSignal time.Time
var MasterConn *net.UDPConn


func NETWORK_init(){
	println("NETWORK_INIT")
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
	sendMessage("n", connected["bc"])

	go alive(conn)

	if Master == false{
		go slave()
	/*}else{
		go slave()
		//go receive(recConn)*/
	}
	
	println(Master)
}


func alive(conn *net.UDPConn){
	println("alive()")
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
		LastSignal = time.Now()
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
		//connected["Master"] = nil
	}
}
