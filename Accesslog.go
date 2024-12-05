// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

//	アクセスログ accesslog 2024-11-27 〜
type Accesslog struct {
	Handler       string
	Remoteaddress string
	Useragent     string
	Formvalues    string
	Eventid       string
	Roomid        int
	Ts            time.Time
}

func GetFeaturedEvents(
	hours int, // 現在からこの時間遡ったところまでを対象とする
	num int, //	抽出するイベント数の最大値
	lmct int, //	抽出するイベントアクセス数の最小値
) (
	eventmap map[string]int, //	指定された条件でアクセス数が多かったイベント（アクセス数の降順）
) {

	eventmap = make(map[string]int)

	sqlst := "select eventid, count(*) ct from accesslog where ts > SUBDATE(now(),INTERVAL ? hour) "
	sqlst += " group by eventid order by ct desc limit ? "

	type event struct {
		Eventid string
		Ct      int
	}
	var events []event

	_, err := Dbmap.Select(&events, sqlst, hours, num)
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
