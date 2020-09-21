package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"unicode/utf8"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	db *sql.DB
)

// Try to connect to the DB server as
// long as it takes to establish a connection
//
func connectToDB() {
	for {
		var err error
		log.Printf("Attempting to connect to DB...")
		db, err = sql.Open("mssql", "server="+conf.DbServer+";user id="+conf.DbUsr+"; password="+conf.DbPwd+";database="+conf.DbDatabase)
		if err == nil {
			err = db.Ping()
			if err != nil {
				checkError(err, "connectToDB Error")
			} else {
				log.Printf("DB Connected...")
				return
			}
		}

		log.Println(err)
		log.Printf("Trying to reconnect to DB Server")
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		sendMessage("Exiting Program, recieved: ")
		sendMessage(fmt.Sprintln(sig))
		if db != nil {
			err := db.Close()
			checkError(err, "Closing DB Error")
		}
		done <- true
	}()

	//Establish DB Connection

	connectToDB()

	go func() {
		workLoop()
		log.Println("Application Error", "Work loop exited, restarting in 5 seconds...", "critical")
		time.Sleep(5 * time.Second)
	}()

	//slackSendMessageTest("Hello from Test")

	sendMessage("Press Ctrl+C to exit: ")
	<-done
}

func syslogSend(msg string) {

	ServerAddr, err := net.ResolveUDPAddr("udp", conf.SysLogSrv+":"+conf.SysLogPort)
	checkError(err, "Setting Syslog Send DST error")
	if err == nil {

		LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
		checkError(err, "LocalAddr error in syslogSend")

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		checkError(err, "Dial error in syslogSend")

		defer Conn.Close()
		buf := []byte(msg)
		if _, err := Conn.Write(buf); err != nil {
			//fmt.Println(msg, err)
			checkError(err, "Connection write error in syslogSend")
		}
	}
}

func workLoop() {

	for {
		sendMessage("Checking for work...")
		alive++

		if db == nil {
			connectToDB()
		}

		rows, err := db.Query("execute get_failed_logins_into_splunk '2020-09-13 00:00:00.000'")
		if err != nil {
			log.Println(err)
		} else {
			defer rows.Close()

			//Test Code
			cols, _ := rows.Columns() // Remember to check err afterwards
			vals := make([]interface{}, len(cols))
			for i := range cols {
				vals[i] = new(sql.RawBytes)
			}
			for rows.Next() {
				err = rows.Scan(vals...)
				if err == nil {
					m := make(map[string]interface{})
					for y, value := range vals {
						//tmpString := fmt.Sprintf("%s", value)
						tmpString2 := trimFirstRune(fmt.Sprintf("%s", value))
						if number, err := strconv.Atoi(tmpString2); err == nil {
							m[cols[y]] = number
						} else {
							m[cols[y]] = tmpString2
						}
					}
					newData, _ := json.Marshal(m)
					fmt.Printf("%s\n", newData)
				}
			}
		}
		sendMessage("Batch Done.")
		err = rows.Err()
		if err != nil {
			log.Println(err)
		}
		rows.Close()

		time.Sleep(conf.WorkDelay * time.Second)
	}

}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func isNum(s string) (bool, int) {
	if number, err := strconv.Atoi(s); err == nil {
		return true, number
	}
	return false, 0

}
