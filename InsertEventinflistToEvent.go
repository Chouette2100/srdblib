package srdblib

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
)

func InsertEventinflistToEvent(eventinflist *[]exsrapi.Event_Inf, bcheck bool) (
	err error,
) {

	var stmt *sql.Stmt

	sql := "INSERT INTO " + Tevent + " (eventid, ieventid, event_name, period, starttime, endtime, noentry,"
	sql += " intervalmin, modmin, modsec, "
	sql += " fromorder, toorder, resethh, resetmm, nobasis, maxdsp, cmap, target, rstatus, maxpoint, achk" //, aclr	未使用
	sql += ") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	//	log.Printf("db.Prepare(sql)\n")
	stmt, Dberr = Db.Prepare(sql)
	if Dberr != nil {
		err = fmt.Errorf("row.Exec(): %w", Dberr)
		return
	}
	defer stmt.Close()

	for _, eventinf := range *eventinflist {
		//	存在確認
		nrow := 0
		if bcheck {

			if Dberr := Db.QueryRow("select count(*) from "+Tevent+" where eventid = ?", eventinf.Event_ID).Scan(&nrow); Dberr != nil {
				err = fmt.Errorf("QeryRow().Scan(): %w", Dberr)
				return
			}
		}

		if bcheck && nrow == 0 || !bcheck && eventinf.Valid {
			//	同一のeventidのデータが存在しない。
			_, Dberr = stmt.Exec(
				eventinf.Event_ID,
				eventinf.I_Event_ID,
				eventinf.Event_name,
				eventinf.Period,
				eventinf.Start_time,
				eventinf.End_time,
				eventinf.NoEntry,
				eventinf.Intervalmin,
				eventinf.Modmin,
				eventinf.Modsec,
				eventinf.Fromorder,
				eventinf.Toorder,
				eventinf.Resethh,
				eventinf.Resetmm,
				eventinf.Nobasis,
				eventinf.Maxdsp,
				eventinf.Cmap,
				eventinf.Target,
				eventinf.Rstatus,
				eventinf.Maxpoint+eventinf.Gscale,
				eventinf.Achk,
			)

			if Dberr != nil {
				err = fmt.Errorf("row.Exec(): %w", Dberr)
				return
			}
			log.Printf("  **Inserted: %-30s %s\n", eventinf.Event_ID, eventinf.Event_name)
		}
	}

	return
}
