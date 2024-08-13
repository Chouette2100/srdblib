// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"time"
)

//	Gorpのための構造体定義

//	0.0.0 新規に作成する

// 配信枠の構造体（配信枠別リスナー別貢献ポイント用）
//	PRIMARY KEY (eventid,userid,sampletm1)
type Timetable struct {
	Eventid string
	Userid int
	Sampletm1 time.Time
	Sampletm2 time.Time
	Stime time.Time
	Etime time.Time
	Target int
	Totalpoint int
	Earnedpoint int
	Status int
}

