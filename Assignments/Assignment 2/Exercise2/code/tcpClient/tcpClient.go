// DOESN'T WORK. EVERYTHING IS IMPLEMENTED IN tcpServer.go INSTEAD.

package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

// 10.100.23.11

func handleRecvConnection(recvConn net.Conn) {
	fmt.Println("Fiber to handle connections started.")
	defer recvConn.Close()

	buffer := make([]byte, 1024)

	for {
		// buffer = make([]byte, 1024)

		fmt.Println("Waiting to read.")
		n, readErr := recvConn.Read(buffer)
		fmt.Println("Reading from connection.")
		if readErr != nil {
			if readErr != io.EOF {
				fmt.Println("Error reading:", readErr)
			}
			break
		}
		fmt.Printf("Received message: %s\n", string(buffer[:n]))
		time.Sleep(time.Second)
	}

}

// Error listening: listen tcp4 10.100.23.11:34933: bind: cannot assign requested address

// HOW DOES TCP LISTENING TCP. WORK???

func main() {
	// 34933 is fixed-size
	// 33546 is 0-terminated
	tcpAddr, errResolve := net.ResolveTCPAddr("tcp4", ":34933")
	if errResolve != nil {
		fmt.Println("Error resolving UDP Address:", errResolve)
	}

	listener, errListen := net.ListenTCP("tcp4", tcpAddr)
	if errListen != nil {
		fmt.Println("Error listening:", errListen)
	}
	fmt.Println("Listening.")
	defer listener.Close()

	for {
		// buffer = make([]byte, 1024)

		fmt.Println("Waiting to accept.")
		recvConn, acceptErr := listener.Accept()
		fmt.Println("Accepted.")
		if acceptErr != nil {
			fmt.Println("Error accepting:", acceptErr)
			continue
		}
		// defer recvConn.close() ???????????
		fmt.Println("Creating fiber to handle connection.")
		go handleRecvConnection(recvConn)
		fmt.Println("Finished creating fiber to handle connection.")
		time.Sleep(time.Second)
	}
}

// addr = new InternetAddress(serverIP, serverPort)
// sock = new Socket(tcp) // TCP, aka SOCK_STREAM
// sock.connect(addr)
// // use sock.recv() and sock.send(), just like with UDP
