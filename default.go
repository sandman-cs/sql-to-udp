package main

import (
	"fmt"
	"log"
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
