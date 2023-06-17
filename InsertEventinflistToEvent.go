package srdblib

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InsertEventinflistToEvent(eventinf []Event_Inf) (
	isntins map[int]bool, //	insertできなかったデータのリスト
	err error,
) {

	var stmt *sql.Stmt

	sql := "INSERT INTO event(eventid, ieventid, event_name, period, starttime, endtime, noentry,"
	sql += " intervalmin, modmin, modsec, "
	sql += " fromorder, toorder, resethh, resetmm, nobasis, maxdsp, cmap, target, rstatus, maxpoint " //	, achk, aclr	未使用
	sql += ") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	log.Printf("db.Prepare(sql)\n")
	stmt, Dberr = Db.Prepare(sql)
	if Dberr != nil {
		err = fmt.Errorf("row.Exec(): %w", Dberr)
		return
	}
	defer stmt.Close()

	for i, event := range eventinf {
		//	存在確認
		nrow := 0
		if Dberr := Db.QueryRow("select * from event where eventid = ?", event.Event_ID).Scan(&nrow); Dberr != nil {
			err = fmt.Errorf("QeryRow().Scan(): %w", Dberr)
			return
		}

		if nrow == 0 {
			//	同一のeventidのデータが存在しない。
			log.Printf("row.Exec()\n")
			_, Dberr = stmt.Exec(
				event.Event_ID,
				event.I_Event_ID,
				event.Event_name,
				event.Period,
				event.Start_time,
				event.End_time,
				event.NoEntry,
				event.Intervalmin,
				event.Modmin,
				event.Modsec,
				event.Fromorder,
				event.Toorder,
				event.Resethh,
				event.Resetmm,
				event.Nobasis,
				event.Maxdsp,
				event.Cmap,
				event.Target,
				event.Rstatus,
				event.Maxpoint + event.Gscale,
			)

			if Dberr != nil {
				err = fmt.Errorf("row.Exec(): %w", Dberr)
				err = fmt.Errorf("row.Exec(): %w", Dberr)
				return
			}
		} else {
			//	同一のeventidのデータが存在する。
			isntins[i] = true
		}
	}

	return
}
