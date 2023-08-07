package srdblib

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
)

// DBから現在開催中のイベントを抜き出す
func SelectCurEventList() (
	eventlist []exsrapi.Event_Inf,
	err error,
) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	sqls := "select e.eventid, e.event_name, e.period, e.starttime, e.endtime, e.fromorder, e.toorder "
	sqls += "from event e join wevent we on e.eventid = we.eventid "
	sqls += " where e.endtime > now() and e.starttime < now()  and we.achk = 0"
	stmt, Dberr = Db.Prepare(sqls)
	if Dberr != nil {
		log.Printf("err=[%s]\n", Dberr.Error())
		err = fmt.Errorf("Db.Prepare(sqls): %w", Dberr)
		return
	}
	defer stmt.Close()

	rows, Dberr = stmt.Query()
	if Dberr != nil {
		log.Printf("err=[%s]\n", Dberr.Error())
		err = fmt.Errorf("stmt.Query(): %w", Dberr)
		return
	}
	defer rows.Close()

	var event exsrapi.Event_Inf
	for rows.Next() {
		Dberr = rows.Scan(&event.Event_ID, &event.Event_name, &event.Period, &event.Start_time, &event.End_time, &event.Fromorder, &event.Toorder)
		if Dberr != nil {
			log.Printf("err=[%s]\n", Dberr.Error())
			err = fmt.Errorf("rows.Next(): %w", Dberr)
			return
		}
		eventlist = append(eventlist, event)
	}
	if Dberr = rows.Err(); Dberr != nil {
		log.Printf("err=[%s]\n", Dberr.Error())
		err = fmt.Errorf("rows.Err(): %w", Dberr)
		return
	}

	return

}
