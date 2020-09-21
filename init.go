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
	conf.LocalEcho = false
	conf.DefaultSrvList.SysLogSrv = "splunk"
	conf.DefaultSrvList.SysLogPort = "514"
	conf.DefaultSrvList.WorkDelay = 60
	conf.ServerName, _ = os.Hostname()
	conf.AppName = "Go - Sql to UDP"
	conf.AppVer = " 1.0"

	//Load Configuration Data
	dat, _ := ioutil.ReadFile("sql-to-udp.json")
	err := json.Unmarshal(dat, &conf)
	checkError(err, "Failed to load sql-to-udp.json")

	go healthCheck()

	for index, element := range conf.DbSrvList {
		// Create Channel and launch working thread.......
		fmt.Println("Creating Worker Thread #", index)
		dbSend = append(dbSend, make(chan string, 1024))
		// Launch Publisher thread......
		go func(index int, element dbSrv) {
			//Load Defaults if needed
			if len(element.DbDatabase) == 0 {
				element.DbDatabase = conf.DefaultSrvList.DbDatabase
			}
			if len(element.DbServer) == 0 {
				element.DbServer = conf.DefaultSrvList.DbServer
			}
			if len(element.DbUsr) == 0 {
				element.DbUsr = conf.DefaultSrvList.DbUsr
			}
			if len(element.DbPwd) == 0 {
				element.DbPwd = conf.DefaultSrvList.DbPwd
			}
			if len(element.SysLogSrv) == 0 {
				element.SysLogSrv = conf.DefaultSrvList.SysLogSrv
			}
			if len(element.SysLogPort) == 0 {
				element.SysLogPort = conf.DefaultSrvList.SysLogPort
			}
			if element.WorkDelay == 0 {
				element.WorkDelay = conf.DefaultSrvList.WorkDelay
			}

			for !closeReceived {
				sendDebugMessage(fmt.Sprintln("Work Loop for Index: ", index, " Starting.."))
				workLoopWithSender(index, element)
				if !closeReceived {
					sendWarnMessage(fmt.Sprintln("Work Loop for Index: ", index, " Exited.."))
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
			//if db == nil {
			//	log.Println("Database Connection Failed")
			//}
			os.Exit((1))
		}

		alive = 0
	}

}
