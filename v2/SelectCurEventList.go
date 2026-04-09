package srdblib

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi/v2"
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
	stmt, err = Db.Prepare(sqls)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		err = fmt.Errorf("Db.Prepare(sqls): %w", err)
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query()
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		err = fmt.Errorf("stmt.Query(): %w", err)
		return
	}
	defer rows.Close()

	var event exsrapi.Event_Inf
	for rows.Next() {
		err = rows.Scan(&event.Event_ID, &event.Event_name, &event.Period, &event.Start_time, &event.End_time, &event.Fromorder, &event.Toorder)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			err = fmt.Errorf("rows.Next(): %w", err)
			return
		}
		eventlist = append(eventlist, event)
	}
	if err = rows.Err(); err != nil {
		log.Printf("err=[%s]\n", err.Error())
		err = fmt.Errorf("rows.Err(): %w", err)
		return
	}

	return

}
