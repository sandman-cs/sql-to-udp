package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	conf          configuration
	closeReceived bool
	alive         int
	dbSend        []chan string
)

func init() {

	//Load Default Configuration Values

	//Load Default Configuration Values
	conf.LocalEcho = true
	conf.SysLogSrv = "splunk"
	conf.SysLogPort = "514"
	conf.ServerName, _ = os.Hostname()
	conf.AppName = "Go - Sql to UDP"
	conf.AppVer = " 1.0"

	//Load Configuration Data
	dat, _ := ioutil.ReadFile("sql-to-udp.json")
	err := json.Unmarshal(dat, &conf)
	checkError(err, "Failed to load aql-to-udp.json")

	go healthCheck()

	for index, element := range conf.DbSrvList {
		// Create Channel and launch working thread.......
		fmt.Println("Creating Worker Thread #", index)
		dbSend = append(dbSend, make(chan string, 1024))
		// Launch Publisher thread......
		go func(index int, element dbSrv) {
			for !closeReceived {
				workLoop()
				if !closeReceived {
					time.Sleep(time.Second)
				}
			}
		}(index, element)

	}

}

func healthCheck() {
	sendMessage("Starting health check thread...")

	for {
		time.Sleep(300 * time.Second)
		if alive == 0 {
			log.Println("Health Check Failed, Workloop has not run for over 5 minute")
			if db == nil {
				log.Println("Database Connection Failed")
			}
			os.Exit((1))
		}

		alive = 0
	}

}
