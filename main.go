package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"unicode/utf8"

	_ "github.com/denisenkom/go-mssqldb"
)

// Try to connect to the DB server as
// long as it takes to establish a connection
//
func connectToDB(server string, usr string, pwd string, database string) *sql.DB {
	for {
		log.Printf("Attempting to connect to DB...")
		sqlString := "server=" + server + ";user id=" + usr + "; password=" + pwd + ";database=" + database
		//log.Println("sqlString: ", sqlString)
		db, err := sql.Open("mssql", sqlString)
		//db, err := sql.Open("mssql", "server="+server+";user id="+usr+"; password="+pwd+";database="+database)
		if err == nil {
			err = db.Ping()
			if err != nil {
				log.Println("db Ping error: ", err)
			} else {
				log.Printf("DB Connected...")
				return db
			}
		}

		log.Println("sql.Open error: ", err)
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
		closeReceived = true
		//Put logic in here for wait groups for clean close
		done <- true
	}()

	sendMessage("Press Ctrl+C to exit: ")
	<-done
}

func workLoopWithSender(offset int, c dbSrv) {

	//Connect to Database
	db := connectToDB(c.DbServer, c.DbUsr, c.DbPwd, c.DbDatabase)

	//Create Send Channel
	sender := make(chan string, 1)

	go func() {
		for {
			sendUDPMessage(c.SysLogSrv, c.SysLogPort, sender)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		sendMessage("Checking for work...")
		alive++

		rows, err := db.Query("execute get_failed_logins_into_splunk")
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
					//fmt.Printf("%s\n", newData)
					sender <- fmt.Sprintf("%s\n", newData)
				}
			}
		}
		sendMessage("Batch Done.")
		err = rows.Err()
		if err != nil {
			log.Println(err)
		}
		rows.Close()

		time.Sleep(c.WorkDelay * time.Second)
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
