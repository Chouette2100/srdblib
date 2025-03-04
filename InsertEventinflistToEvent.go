/*!
Copyright © 2023 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/
package srdblib

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi/v2"
)

/*
イベント情報リストをイベントテーブルに格納する。
格納できたときは (*eventinflist)[i].Valid = false とする。
格納できないときは (*eventinflist)[i].Valid = true とする。
(*eventinflist)[i].Valid はイベント情報とは無関係で
処理の状況を示すために使われる。
*/
func InsertEventinflistToEvent(
	tevent	string,	//	insertするテーブル
	eventinflist *[]exsrapi.Event_Inf, //	イベント情報リスト
	bcheck bool, //	true:	キーが同一のレコードの存在チェックを行う
) (
	err error,
) {

	var stmt *sql.Stmt

	sql := "INSERT INTO " + tevent + " (eventid, ieventid, event_name, period, starttime, endtime, noentry,"
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

	for i, eventinf := range *eventinflist {
		//	存在確認
		nrow := 0
		if bcheck {
			if Dberr := Db.QueryRow("select count(*) from "+tevent+" where eventid = ?", eventinf.Event_ID).Scan(&nrow); Dberr != nil {
				err = fmt.Errorf("QeryRow().Scan(): %w", Dberr)
				return
			}
		}

		if bcheck && nrow == 0 || !bcheck && eventinf.Valid {
			//	同一のeventidのデータが存在しないのでデータを格納する。
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
			(*eventinflist)[i].Valid = false
		} else {
			(*eventinflist)[i].Valid = true
		}
	}

	return
}
