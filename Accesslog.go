//	Copyright © 2024 chouette2100@gmail.com
//	Released under the MIT license
//	https://opensource.org/licenses/mit-license.php
package srdblib

import (
	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Accesslog struct {
	Handler       string
	Remoteaddress string
	Useragent     string
	Formvalues    string
	Eventid       string
	Roomid        int
	Ts            time.Time
}

func GetFeaturedEvents(hours int, num int) (eventmap map[string]bool) {

	eventmap = make(map[string]bool)

	sqlst := "select eventid from accesslog where ts > SUBDATE(now(),INTERVAL ? hour) group by eventid order by count(*) desc limit ? "

	type Event struct {
		Eventid string
	}
	var events []Event
	
	_, err := Dbmap.Select(&events, sqlst, hours, num)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	for _, v := range events {
		eventmap[v.Eventid] = true
	}

	return
}
