/*Handles the UDP. Functions making new connections, receiving and
sending messages.*/

package network

import("net"
	/*"time"*/
	/*"encoding/json"*/)

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
		broadcast.Conn = conn
	}else{
		//connected[ip] = conn
	}

	//Sends a message to the new connection that they should connect as well
	text := "c"
	sendMessage(text, conn)
	return conn
}


func receive(conn *net.UDPConn){
	//Receive message from network	
	received := make([]byte,500)
	//var s string
	for ; true ; {
		_, _, _ = conn.ReadFromUDP(received)
		//_ = json.Unmarshal(received, &s)		
        //println("Motatt: ", s)
		println("Motatt: ", string(received))
		go receiveMessage(string(received))
	}
	conn.Close()
	println("Stenger mottak")
}


func send(s string, conn *net.UDPConn){
    b := []byte(s)
	//b, _ := json.Marshal(s)
	_, _ = conn.Write(b)
}
