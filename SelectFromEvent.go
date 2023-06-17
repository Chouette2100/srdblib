package srdblib

import (
	"fmt"
	"log"
	"time"
)

func SelectFromEvent(eventid string) (
	peventinf *Event_Inf,
	err error,
) {

	eventinf := Event_Inf{}
	peventinf = &eventinf
	log.Printf("eventid=[%s]\n", eventid)
	log.Printf("eventinf=[%v]\n", eventinf)

	sql := "select eventid,ieventid,event_name, period, starttime, endtime, noentry, intervalmin, modmin, modsec, "
	sql += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, cmap, target, `rstatus`, maxpoint "
	sql += " from event where eventid = ?"
	log.Printf("sql=[%s]\n", sql)
	Dberr = Db.QueryRow(sql, eventid).Scan(
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
	)

	if Dberr != nil {
		if Dberr.Error() != "sql: no rows in result set" {
			peventinf = nil
			return
		} else {
			err = fmt.Errorf("row.Exec(): %w", Dberr)
			log.Printf("%s\n", sql)
			log.Printf("err=[%v]\n", err)
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

	log.Printf("eventinf=[%v]\n", eventinf)

	//	log.Printf("Start_data=%f Dperiod=%f\n", eventinf.Start_date, eventinf.Dperiod)

	return
}
