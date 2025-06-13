package srdblib

import (
	"fmt"
	//	"log"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi/v2"
)

// 指定した条件に該当するイベントのリストを作る。
func SelectEventinflistFromEvent(
	cond int, // 抽出条件	-1:終了したイベント、0: 開催中のイベント、1: 開催予定のイベント、 99: すべて
	keyword string, // イベント名検索キーワード
	offset int,
	limit int,
) (
	eventinflist []exsrapi.Event_Inf,
	err error,
) {

	tnow := time.Now().Truncate(time.Second)

	sqls := "select eventid, ieventid, event_name, period, starttime, endtime, noentry, intervalmin, modmin, modsec, "
	sqls += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, cmap, target, `rstatus`, maxpoint, achk "
	sqls += " from event"
	switch cond {
	case -1:
		sqls += " and endtime < ?"
	case 0:
		sqls += " and starttime < ? and endtime > ?"
	case 1:
		sqls += " and we.starttime > ?"
	case 99:
	default:
		err = fmt.Errorf("cond=%d is not valid", cond)
		return
	}
	if keyword != "" {
		sqls += " and we.event_name like ?"
	}

	switch cond {
	case -1: //	終了済みイベント
		//	sqls += " order by we.endtime desc, we.starttime desc "
		sqls += " order by we.starttime desc, we.endtime desc "
	case 0: //	開催中のイベント
		sqls += " order by we.endtime, we.starttime "
	case 1: //	開催予定のイベント
		sqls += " order by we.starttime, we.endtime "
	}

	//	log.Printf("sql=[%s]\n", sqls)
	var stmts *sql.Stmt
	stmts, err = Db.Prepare(sqls)
	if err != nil {
		err = fmt.Errorf("Prepare(sqls): %w", err)
		return
	}
	defer stmts.Close()

	var rows *sql.Rows

	switch {
	case cond == 99 && keyword == "":
		rows, err = stmts.Query()
	case cond == 99 && keyword != "":
		rows, err = stmts.Query("%" + keyword + "%")
	case cond == 0 && keyword == "":
		rows, err = stmts.Query(tnow, tnow)
	case cond == 0 && keyword != "":
		rows, err = stmts.Query(tnow, tnow, "%"+keyword+"%")
	case cond != 0 && keyword == "":
		rows, err = stmts.Query(tnow)
	case cond != 0 && keyword != "":
		rows, err = stmts.Query(tnow, "%"+keyword+"%")
	}
	if err != nil {
		err = fmt.Errorf("Query(tnow): %w", err)
		return
	}
	defer rows.Close()

	eventinflist = make([]exsrapi.Event_Inf, 0)
	eventinf := exsrapi.Event_Inf{}

	for rows.Next() {

		err = rows.Scan(
			&eventinf.Event_ID,
			&eventinf.I_Event_ID,
			&eventinf.Event_name,
			&eventinf.Period,
			&eventinf.Start_time,
			&eventinf.End_time,
			&eventinf.NoEntry,
			&eventinf.Intervalmin,
			&eventinf.Modmin,
			&eventinf.Modsec,
			&eventinf.Fromorder,
			&eventinf.Toorder,
			&eventinf.Resethh,
			&eventinf.Resetmm,
			&eventinf.Nobasis,
			&eventinf.Maxdsp,
			&eventinf.Cmap,
			&eventinf.Target,
			&eventinf.Rstatus,
			&eventinf.Maxpoint,
			&eventinf.Achk,
		)

		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return
			} else {
				err = fmt.Errorf("row.Exec(): %w", err)
				return
			}
		}

		//	log.Printf("eventno=%d\n", Event_inf.Event_no)

		start_date := eventinf.Start_time.Truncate(time.Hour).Add(-time.Duration(eventinf.Start_time.Hour()) * time.Hour)
		end_date := eventinf.End_time.Truncate(time.Hour).Add(-time.Duration(eventinf.End_time.Hour())*time.Hour).AddDate(0, 0, 1)

		//	log.Printf("start_t=%v\nstart_d=%v\nend_t=%v\nend_t=%v\n", Event_inf.Start_time, start_date, Event_inf.End_time, end_date)

		eventinf.Start_date = float64(start_date.Unix()) / 60.0 / 60.0 / 24.0
		eventinf.Dperiod = float64(end_date.Unix())/60.0/60.0/24.0 - eventinf.Start_date

		eventinf.Gscale = eventinf.Maxpoint % 1000
		eventinf.Maxpoint = eventinf.Maxpoint - eventinf.Gscale

		//	log.Printf("eventinf=[%v]\n", eventinf)

		eventinflist = append(eventinflist, eventinf)

		//	log.Printf("Start_data=%f Dperiod=%f\n", eventinf.Start_date, eventinf.Dperiod)
	}

	return
}
