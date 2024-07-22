// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"time"
)

//	Gorpのための構造体定義

//	0.0.0 新規に作成する

// イベント構造体
// PRIMARY KEY (eventid)
type Event struct {
	Eventid     string
	Ieventid    int
	Event_name  string
	Period      string
	Starttime   time.Time
	Endtime     time.Time
	Noentry     int
	Intervalmin int
	Modmin      int
	Modsec      int
	Fromorder   int
	Toorder     int
	Resethh     int
	Resetmm     int
	Nobasis     int
	Maxdsp      int
	Cmap        int
	Target      int
	Rstatus     string
	Maxpoint    int
	Achk        int
	Aclr        int
}

//	event := Event{
//		Intervalmin: 5,
//		Modmin:      4,
//		Modsec:      10,
//		Fromorder:   1,
//		Toorder:     10,
//		Resethh:     4,
//		Resetmm:     0,
//		Nobasis:     164614,
//		Maxdsp:      10,
//		Cmap:        1,
//	}

//	イベントに参加しているユーザの構造体
//  PRIMARY KEY (`eventid`,`userno`)
type Eventuser struct {
  Eventid string
  Userno int
  Istarget string
  Iscntrbpoints string
  Graph string
  Color string
  Point int
  Vld int
}
//	eventuser := Eventuser{
//		vld: 1,
//	}
