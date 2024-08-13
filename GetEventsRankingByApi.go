//	Copyright © 2024 chouette.21.00@gmail.com
//	Released under the MIT license
//	https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/Chouette2100/srapi"
)

func GetEventsRankingByApi(
	client *http.Client, //	HTTPクライアント
	eid string, //	イベントID
) (
	pranking *srapi.Eventranking,
	err error,
) {

	// イベントの詳細を得る、ここではIeventidが必要である
	var row interface{}
	row, err = Dbmap.Get(Event{}, eid)
	if err != nil {
		err = fmt.Errorf("Dbmap.Get(): %w", err)
		return
	}
	event := row.(*Event)

	// イベントに参加しているルームを取得する
	roomlistinf, err := srapi.GetRoominfFromEventByApi(client, event.Ieventid, 1, 1)
	if err != nil {
		err = fmt.Errorf("GetRoominfFromEventByApi(): %w", err)
		return
	}
	roomid := roomlistinf.RoomList[0].Room_id

	// イベント結果を取得する
	bid := 0
	if strings.Contains(event.Eventid, "block_id") {
		eida := strings.Split(event.Eventid, "=")
		bid, _ = strconv.Atoi(eida[1])
	}
	pranking, err = srapi.ApiEventsRanking(client, (event).Ieventid, roomid, bid)
	if err != nil {
		err = fmt.Errorf("ApiEventsRanking(): %w", err)
		return
	}
	return
}