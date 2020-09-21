package main

import (
	"fmt"
	"log"
	"net"
)

//checkError function
func checkError(err error, txt string) {
	if err != nil {
		log.Println("checkError: ", txt, "\n: ", err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func sendUDPMessage(msg string) {
	ServerAddr, err := net.ResolveUDPAddr("udp", conf.SysLogSrv+":"+conf.SysLogPort)
	checkError(err, "Error resolving syslog server address...")
	if err == nil {

		LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
		checkError(err, "Error creating socket to send UDP message...")

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		checkError(err, "Error connecting too syslog destination...")

		defer Conn.Close()
		buf := []byte(msg)
		if _, err := Conn.Write(buf); err != nil {
			checkError(err, "Error sending data too syslog destination...")
		}
	}
}

//sendMessage to udp listener
func sendMessage(msg string) {

	log.Println(msg)
}

//sendMessage to udp listener
func sendDebugMessage(msg string) {

	log.Println("Debug: ", msg)
}

//sendMessage to udp listener
func sendWarnMessage(msg string) {

	log.Println("Warn: ", msg)
}
