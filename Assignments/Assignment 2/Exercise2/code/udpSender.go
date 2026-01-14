package main

import (
	"fmt"
	"net"
	"strconv"
)

func send(remoteIPStr string, remotePort int) {
	ipAndPort := remoteIPStr + ":" + strconv.Itoa(remotePort)
	// fmt.Printf("\"%v\"\n", ipAndPort)
	udpAddr, errResolve := net.ResolveUDPAddr("udp4", ipAndPort)
	if errResolve != nil {
		fmt.Println("Error resolving UDP Address:", errResolve)
	}

	sendConn, errDial := net.DialUDP("udp4", nil, udpAddr)
	if errDial != nil {
		fmt.Println("Error dialing:", errDial)
	}
	defer sendConn.Close()

	messageBuffer := []byte("Kritt Er KuL")

	_, errSending := sendConn.Write(messageBuffer)
	if errSending != nil {
		fmt.Println("Error sending:", errSending)
	}
}

func recv() {

	udpAddr, errResolve := net.ResolveUDPAddr("udp4", ":20023") // 30000
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

func main() {
	const workspaceNum int = 23
	const remotePort int = 20000 + workspaceNum
	const remoteIPStr string = "10.100.23.11"
	// const remoteIPStr string = "10.100.23.255"
	// const remoteIPStr string = "255.255.255.255"

	send(remoteIPStr, remotePort)
	recv()
}

// // if sending directly to a single remote machine:
//     addr = new Address(remoteIP, remotePort)
//     sock = new Socket(udp)

//     // either: set up the socket to use a single remote address
//         sock.connect(addr)
//         sock.send(message)
//     // or: set up the remote address when sending
//         sock.sendTo(message, addr)

// // if sending on broadcast:
// // you have to set up the BROADCAST socket option before calling connect / sendTo
//     broadcastIP = #.#.#.255 // First three bytes are from the local IP, or just use 255.255.255.255
//     addr = new InternetAddress(broadcastIP, port)
//     sendSock = new Socket(udp) // UDP, aka SOCK_DGRAM
//     sendSock.setOption(broadcast, true)
//     sendSock.sendTo(message, addr)
