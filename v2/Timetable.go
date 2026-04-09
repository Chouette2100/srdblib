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
	Eventid string // イベントID、=event_url_keyを設定する
	Userid int	// ルームIDを設定する
	Sampletm1 time.Time	// 貢献ポイント取得予定（希望）時刻（貢献ポイントが落ち着く頃を指定する）
	Sampletm2 time.Time	// 貢献ポイント取得時刻、実際にデータを取得した時刻が格納される
	Stime time.Time		// 配信枠の開始時刻
	Etime time.Time		// 配信枠の終了推定時刻
	Target int		// -1を設定する。配信枠での無料ギフトの最大獲得ポイントを格納される（最近はあんまり意味ないし、一概にいくらと決めにくい）
	Totalpoint int		// 貢献ポイントの合計値が格納される
	Earnedpoint int		// 獲得ポイントの増分を設定する
	Status int		// 0を設定する、貢献ポイントの増分が格納されたら1が格納される
}

