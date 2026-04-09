package srdblib

import (
	"time"
)

type Points struct {
	Ts      time.Time
	User_id int
	Eventid string
	Point   int
	Rank    int
	Gap     int
	Pstatus string
	Qstatus string
	Ptime   string
	Qtime   string
}
