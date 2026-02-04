package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const port int = 20023

func sendHeartbeat(counter int) {
	remoteIPStr := "localhost"
	remotePort := port
	ipAndPort := remoteIPStr + ":" + strconv.Itoa(remotePort)

	udpAddr, errResolve := net.ResolveUDPAddr("udp4", ipAndPort)
	if errResolve != nil {
		fmt.Println("Error resolving UDP Address:", errResolve)
	}

	sendConn, errDial := net.DialUDP("udp4", nil, udpAddr)
	if errDial != nil {
		fmt.Println("Error dialing:", errDial)
	}
	defer sendConn.Close()

	for {
		counter += 1
		messageBuffer := []byte(strconv.Itoa(counter))
		time.Sleep(time.Second)
		_, errSending := sendConn.Write(messageBuffer)
		if errSending != nil {
			fmt.Println("Error sending:", errSending)
		}
	}
}

func receiveHearbeat() {
	udpAddr, errResolve := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%v", port))
	if errResolve != nil {
		fmt.Println("Error resolving UDP Address:", errResolve)
	}

	recvConn, errListen := net.ListenUDP("udp4", udpAddr)
	if errListen != nil {
		fmt.Println("Error listening:", errListen)
	}
	defer recvConn.Close()

	lastReceivedCount := -1
	buffer := make([]byte, 1024)
	for {
		recvConn.SetReadDeadline(time.Now().Add(time.Second * 3))
		n, addr, readErr := recvConn.ReadFromUDP(buffer)
		if readErr != nil {
			// Check if the error is a timeout
			if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Read timeout:", readErr)
				recvConn.Close()
				startBackup()
				sendHeartbeat(lastReceivedCount)
				return
			}
			panic(readErr)
		}
		fmt.Printf("Received message from %s: %s\n", addr.String(), string(buffer[:n]))
		num, _ := strconv.Atoi(string(buffer[:n]))
		lastReceivedCount = num
	}
}

func startBackup() {
	cmd := exec.Command(
		"gnome-terminal",
		"--",
		"bash", "-lc", "go run main.go",
	)
	fmt.Println("Terminal started.")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := os.Args
	isFirstLaunch := len(args) > 1

	fmt.Println("Is first launch:", isFirstLaunch)
	if isFirstLaunch {
		startBackup()
		sendHeartbeat(0)
	} else {
		receiveHearbeat()
	}

}
