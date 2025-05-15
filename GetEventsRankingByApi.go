// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/Chouette2100/srapi/v2"
)

func GetEventsRankingByApi(
	client *http.Client, //	HTTPクライアント
	eid string, //	イベントID
	mode int, // 1: イベント開催中 2: イベント終了後
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
	if row == nil {
		err = fmt.Errorf("Dbmap.Get(Event{}, eid), %s not found", eid)
		return
	}
	event := row.(*Event)

	//	mode==1のときイベント終了後ならエラーとする。
	if mode == 1 && time.Now().After(event.Endtime.Add(1*time.Minute)) {
		//	err = fmt.Errorf("%s has ended", event.Eventid)
		pranking = &srapi.Eventranking{}
		r := make([]struct {
			Point int `json:"point"`
			Room  struct {
				Name        string `json:"name"`
				ImageSquare string `json:"image_square"`
				RoomID      int    `json:"room_id"`
				Image       string `json:"image"`
			} `json:"room"`
			Rank int `json:"rank"`
		}, 0)
		pranking.Ranking = r
		return
	}

	// イベントに参加しているルームを取得する
	//	ApiEventsRanking()にはイベントにエントリーしているルームのルームIDとが一つ必要だから
	//	REVIEW:  ブロックイベントの場合はこのルームがランキングを取得するブロックが違う可能性がある。この方法でいいのか？
	/*
		roomlistinf, err := srapi.GetRoominfFromEventByApi(client, event.Ieventid, 1, 1)
		if err != nil {
			err = fmt.Errorf("GetRoominfFromEventByApi(): %w", err)
			return
		}
	*/
	roomlistinf, err := srapi.GetEventRankingByApi(client, eid, 1, 1)
	if err != nil {
		err = fmt.Errorf("GetRoominfFromEventByApi(): %w", err)
		return
	}
	// if len(roomlistinf.RoomList) == 0 {
	if len(roomlistinf.Ranking) == 0 {
		//  エントリーしているルームが一つもない。
		var intf []interface{}
		intf, err = Dbmap.Select(Eventuser{}, "select userno from eventuser where eventid=?", eid)
		if err != nil {
			err = fmt.Errorf("Dbmap.Select(): %w", err)
			return
		}
		if len(intf) == 0 {
			err = fmt.Errorf("GetRoominfFromEventByApi(): %s has no room", event.Eventid)
			return
		}
		roomlistinf.Ranking = make([]srapi.Ranking, 1)
		roomlistinf.Ranking[0].RoomID = intf[0].(*Eventuser).Userno
	}

	roomid := roomlistinf.Ranking[0].RoomID

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
