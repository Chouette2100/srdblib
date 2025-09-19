// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"time"
)

// アクセスログ accesslog 2024-11-27 〜
type Accesslog struct {
	Handler       string
	Remoteaddress string
	Useragent     string
	Referer       string
	Formvalues    string
	Eventid       string
	Roomid        int
	Ts            time.Time
}

// 参照回数の多いイベントを抽出する
func GetFeaturedEvents(
	mode string, // "scheduled", "current" or "closed"  現時点（2025-07-02）では "scheduled" は無意味
	hours int, // 現在からこの時間遡ったところまでを対象とする
	num int, //	抽出するイベント数の最大値
	lmct int, //	抽出するイベントアクセス数の最小値
) (
	eventmap map[string]int, //	指定された条件でアクセス数が多かったイベント（アクセス数の降順）
	err error, //	エラーがあれば返す
) {

	eventmap = make(map[string]int)

	// sqlst := "select eventid, count(*) ct from accesslog where ts > SUBDATE(now(),INTERVAL ? hour) AND is_bot = 0 "
	// sqlst += " group by eventid order by ct desc limit ? "

	sqlst := "select a.eventid, count(*) ct from accesslog a "
	if mode == "scheduled" {
		sqlst += " join wevent e "
	} else {
		sqlst += " join event e "
	}
	sqlst += "on e.eventid = BINARY a.eventid "
	sqlst += " where a.ts > SUBDATE(now(),INTERVAL ? hour) AND a.is_bot = 0 "
	switch mode {
	case "scheduled":
		sqlst += " and e.starttime > Now() "
	case "current":
		sqlst += " and e.starttime <= Now() and e.endtime > Now() "
	case "closed":
		sqlst += " and e.endtime <= Now() "
	default:
		err = fmt.Errorf("GetFeaturedEvents: invalid mode %s", mode)
		log.Printf("Error: %v", err)
		return
	}
	sqlst += " group by a.eventid order by ct desc limit ? "

	type event struct {
		Eventid string
		Ct      int
	}
	var events []event

	_, err = Dbmap.Select(&events, sqlst, hours, num)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	for _, v := range events {
		if v.Ct < lmct {
			continue
		}
		eventmap[v.Eventid] = v.Ct
	}

	return
}
