package main

import (
	"fmt"
	"net"
)

// 10.100.23.11:47102
// Hello from UDP server at 10.100.23.11!

func main() {

	udpAddr, errResolve := net.ResolveUDPAddr("udp4", ":30000") // 30000
	if errResolve != nil {
		fmt.Println("Error resolving UDP Address:", errResolve)
	}

	recvConn, errListen := net.ListenUDP("udp4", udpAddr)
	if errListen != nil {
		fmt.Println("Error listening:", errListen)
	}
	defer recvConn.Close()

	buffer := make([]byte, 1024)

	for {
		// buffer = make([]byte, 1024)

		n, addr, readErr := recvConn.ReadFromUDP(buffer)
		if readErr != nil {
			fmt.Println("Error reading:", readErr)
			continue
		}
		fmt.Printf("Received message from %s: %s\n", addr.String(), string(buffer[:n]))
	}
}

// the address we are listening for messages on
// we have no choice in IP, so use 0.0.0.0, INADDR_ANY, or leave the IP field empty
// the port should be whatever the sender sends to
// alternate names: sockaddr, resolve(udp)addr,
// InternetAddress addr

// a socket that plugs our program to the network. This is the "portal" to the outside world
// alternate names: conn
// UDP is sometimes called SOCK_DGRAM. You will sometimes also find UDPSocket or UDPConn as separate types
// recvSock = new Socket(udp)

// bind the address we want to use to the socket
// recvSock.bind(addr)

// a buffer where the received network data is stored
// byte[1024] buffer

// an empty address that will be filled with info about who sent the data
// InternetAddress fromWho

// loop {
// 	// clear buffer (or just create a new one)

// 	// receive data on the socket
// 	// fromWho will be modified by ref here. Or it's a return value. Depends.
// 	// receive-like functions return the number of bytes received
// 	// alternate names: read, readFrom
// 	numBytesReceived = recvSock.receiveFrom(buffer, ref fromWho)

// 	// the buffer just contains a bunch of bytes, so you may have to explicitly convert it to a string

// 	// optional: filter out messages from ourselves
// 	if(fromWho.IP != localIP){
// 		// do stuff with buffer
// 	}
// }
