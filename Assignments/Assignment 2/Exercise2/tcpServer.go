package main

import (
	"fmt"
	"net"
	"time"
	// "io"
)

func main() {
	const workspaceNum int = 23
	const localPort int = 20000 + workspaceNum
	const localIPStr string = "10.100.23.33"
	// localIPAndPort := localIPStr + ":" + strconv.Itoa(localPort)

	// 34933 is fixed-size
	// 33546 is 0-terminated
	tcpAddr, errResolve := net.ResolveTCPAddr("tcp4", "10.100.23.11:33546")
	if errResolve != nil {
		fmt.Println("Error resolving UDP Address:", errResolve)
	}

	conn, errDial := net.DialTCP("tcp4", nil, tcpAddr)
	if errDial != nil {
		fmt.Println("Error dialing:", errDial)
	}
	defer conn.Close()
	fmt.Println("Connection created.")

	message := []byte("Connect to: 10.100.23.33")
	message = append(message, 0)
	_, errWrite := conn.Write(message)
	if errWrite != nil {
		fmt.Println("Error writing:", errWrite)
	}
	fmt.Println("Messsage sent.")

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buffer := make([]byte, 1024)
	n, errRead := conn.Read(buffer)
	if errRead != nil {
		fmt.Println("Error reading:", errRead)
	}
	fmt.Println("Received message:", string(buffer[:n]))

	message = []byte("KrittErKuL")
	message = append(message, 0)
	_, errWrite = conn.Write(message)
	if errWrite != nil {
		fmt.Println("Error writing:", errWrite)
	}
	fmt.Println("Messsage sent.")

	message = []byte("kilopascal")
	message = append(message, 0)
	_, errWrite = conn.Write(message)
	if errWrite != nil {
		fmt.Println("Error writing:", errWrite)
	}
	fmt.Println("Messsage sent.")

}

// fixed-size

// message := "FIXED-SIZE"
// buffer := make([]byte, 1024)
// copy(buffer, message)
// _, errWrite := conn.Write(buffer)

// 0-terminated

// message := []byte("0-TERMINATED")
// message = append(message, 0)
// _, errWrite := conn.Write(message)

// // Send a message to the server:  "Connect to: " <your IP> ":" <your port> "\0"

// // do not need IP, because we will set it to listening state
// addr = new InternetAddress(localPort)
// acceptSock = new Socket(tcp)

// // You may not be able to use the same port twice when you restart the program, unless you set this option
// acceptSock.setOption(REUSEADDR, true)
// acceptSock.bind(addr)

// loop {
//     // backlog = Max number of pending connections waiting to connect()
//     newSock = acceptSock.listen(backlog)

//     // Spawn new thread to handle recv()/send() on newSock
// }
