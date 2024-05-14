/*
!
Copyright © 2024 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/
package srdblib

import (
	"fmt"
	//	"io"
	"log"
	//	"os"
	"strconv"
	"strings"
	"time"

	"net/http"

	//	"github.com/go-gorp/gorp"
	//      "gopkg.in/gorp.v2"

	"github.com/dustin/go-humanize"

	//	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
)

/*

0.0.1 UserでGenreが空白の行のGenreをSHOWROOMのAPIで更新する。
0.0.2 Userでirankが-1の行のランクが空白の行のランク情報をSHOWROOMのAPIで更新する。
0.1.0 DBのアクセスにgorpを導入する。
0.1.1 database/sqlを使った部分（コメント）を削除する

*/

//	const Version = "0.1.1"

type User struct {
	Userno       int
	Userid       string
	User_name    string
	Longname     string
	Shortname    string
	Genre        string
	GenreID      int
	Rank         string
	Nrank        string
	Prank        string
	Irank        int
	Inrank       int
	Iprank       int
	Itrank       int
	Level        int
	Followers    int
	Fans         int
	FanPower     int
	Fans_lst     int
	FanPower_lst int
	Ts           time.Time
	Getp         string
	Graph        string
	Color        string
	Currentevent string
}

// Rank情報からランクのソートキーを作る
func MakeSortKeyOfRank(rank string, nextscore int) (
	irank int,
) {
	r2n := map[string]int{
		"SS-5":    50000000,
		"SS-4":    60000000,
		"SS-3":    70000000,
		"SS-2":    80000000,
		"SS-1":    90000000,
		"S-5":     150000000,
		"S-4":     160000000,
		"S-3":     170000000,
		"S-2":     180000000,
		"S-1":     190000000,
		"A-5":     250000000,
		"A-4":     260000000,
		"A-3":     270000000,
		"A-2":     280000000,
		"A-1":     290000000,
		"B-5":     350000000,
		"B-4":     360000000,
		"B-3":     370000000,
		"B-2":     380000000,
		"B-1":     390000000,
		"C-10":    400000000,
		"C-9":     410000000,
		"C-8":     420000000,
		"C-7":     430000000,
		"C-6":     440000000,
		"C-5":     450000000,
		"C-4":     460000000,
		"C-3":     470000000,
		"C-2":     480000000,
		"C-1":     490000000,
		"unknown": 1000000000, //	SHOWROOMのアカウントを削除した配信者さん
		//	888888888: irank 未算出
	}

	if sk, ok := r2n[rank]; ok {
		irank = sk + nextscore
	} else {
		irank = 999999999 //	(アイドルで)SHOWRANKの対象ではない配信者さん
	}

	return
}

/*
ルーム番号 user.Userno が テーブル user に存在しないときは新しいデータを挿入し、存在するときは 既存のデータを更新する。
*/
func UpinsUserSetProperty(client *http.Client, tnow time.Time, user *User) (
	err error,
) {

	row, err := Dbmap.Get(User{}, user.Userno)
	if err != nil {
		err = fmt.Errorf("Get(userno=%d) returned error. %w", user.Userno, err)
		return err
	} else {
		if row == nil {
			err = InsertIntoUser(client, tnow, user.Userno)
			if err != nil {
				err = fmt.Errorf("InsertIntoUser(userno=%d) returned error. %w", user.Userno, err)
			}
		} else {
			usert := row.(*User)
			//	if usert.Ts.After(tnow.Add(-12 * time.Hour)) {
			if usert.Ts.After(tnow.Add(-12 * time.Second)) {
				log.Printf("UpinsUserSetProperty(userno=%d %s) skipped. %v", user.Userno, usert.User_name, usert.Ts)
				return nil
			}
			err = UpdateUserSetProperty(client, tnow, usert)
			if err != nil {
				err = fmt.Errorf("UpdateUserSetProperty(userno=%d) returned error. %w", user.Userno, err)
			}
		}
	}
	return
}

/*
テーブル user を SHOWROOMのAPI api/roomprofile を使って得られる情報で更新する。
*/
func UpdateUserSetProperty(client *http.Client, tnow time.Time, user *User) (
	err error,
) {

	//	ユーザーのランク情報を取得する
	ria, err := srapi.ApiRoomProfileAll(client, fmt.Sprintf("%d", user.Userno))
	if err != nil {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %w", user.Userno, err)
		return err
	}
	if ria.Errors != nil {
		//	err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %v", userno, ria.Errors)
		//	return err
		ria.ShowRankSubdivided = "unknown"
		ria.NextScore = 0
		ria.PrevScore = 0
	}

	if ria.ShowRankSubdivided == "" {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned empty.ShowRankSubdivided", user.Userno)
		return err
	}

	//	user.Userno =	キー
	user.Userid = ria.RoomURLKey
	user.User_name = ria.RoomName
	//	user.Longname =		新規作成時に設定され、変更はSRCGIで行う
	//	user.Shortname =　　〃
	user.Genre = ria.GenreName
	user.GenreID = ria.GenreID
	user.Rank = ria.ShowRankSubdivided
	user.Nrank = humanize.Comma(int64(ria.NextScore))
	user.Prank = humanize.Comma(int64(ria.PrevScore))
	user.Irank = MakeSortKeyOfRank(ria.ShowRankSubdivided, ria.NextScore)
	user.Inrank = ria.NextScore
	user.Iprank = ria.PrevScore
	user.Itrank = 0
	user.Level = ria.RoomLevel
	user.Followers = ria.FollowerNum

	pafr, err := srapi.ApiActivefanRoom(client, strconv.Itoa(user.Userno), tnow.Format("200601"))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", user.Userno, err)
		return err
	}
	user.Fans = pafr.TotalUserCount
	user.FanPower = pafr.FanPower
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(user.Userno), fmt.Sprintf("%04d%02d", yy, mm))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", user.Userno, err)
		return err
	}
	user.Fans_lst = pafr.TotalUserCount
	user.FanPower_lst = pafr.FanPower

	user.Ts = tnow
	//	user.Getp =		ここから三つのフィールドは現在使われていない。
	//	user.Graph =	default: ''
	//	user.Color =
	eurl := ria.Event.URL
	eurla := strings.Split(eurl, "/")
	user.Currentevent = eurla[len(eurla)-1]

	//	cnt, err := Dbmap.Update(user)
	_, err = Dbmap.Update(user)
	if err != nil {
		log.Printf("error! %v", err)
		return
	}
	//	log.Printf("cnt = %d\n", cnt)

	log.Printf("UPDATE userno=%d rank=%s nscore=%d pscore=%d longname=%s\n", user.Userno, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)
	return
}

/*
テーブル user に新しいデータを追加する
*/
func InsertIntoUser(client *http.Client, tnow time.Time, userno int) (
	err error,
) {

	//	ユーザーのランク情報を取得する
	ria, err := srapi.ApiRoomProfileAll(client, fmt.Sprintf("%d", userno))
	if err != nil {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %w", userno, err)
		return err
	}
	if ria.Errors != nil {
		//	err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %v", userno, ria.Errors)
		//	return err
		ria.ShowRankSubdivided = "unknown"
		ria.NextScore = 0
		ria.PrevScore = 0
	}

	if ria.ShowRankSubdivided == "" {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned empty.ShowRankSubdivided", userno)
		return err
	}

	user := new(User)

	user.Userno = userno
	user.Userid = ria.RoomURLKey
	user.User_name = ria.RoomName
	user.Longname = ria.RoomName
	user.Shortname = strconv.Itoa(userno % 100)
	user.Genre = ria.GenreName
	user.GenreID = ria.GenreID
	user.Rank = ria.ShowRankSubdivided
	user.Nrank = humanize.Comma(int64(ria.NextScore))
	user.Prank = humanize.Comma(int64(ria.PrevScore))
	user.Irank = MakeSortKeyOfRank(ria.ShowRankSubdivided, ria.NextScore)
	user.Inrank = ria.NextScore
	user.Iprank = ria.PrevScore
	user.Itrank = 0
	user.Level = ria.RoomLevel
	user.Followers = ria.FollowerNum

	pafr, err := srapi.ApiActivefanRoom(client, strconv.Itoa(user.Userno), tnow.Format("200601"))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", user.Userno, err)
		return err
	}
	user.Fans = pafr.TotalUserCount
	user.FanPower = pafr.FanPower
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(user.Userno), fmt.Sprintf("%04d%02d", yy, mm))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", user.Userno, err)
		return err
	}
	user.Fans_lst = pafr.TotalUserCount
	user.FanPower_lst = pafr.FanPower

	user.Ts = tnow
	user.Getp = ""
	user.Graph = ""
	user.Color = ""
	eurl := ria.Event.URL
	eurla := strings.Split(eurl, "/")
	user.Currentevent = eurla[len(eurla)-1]

	//	cnt, err := Dbmap.Update(user)
	err = Dbmap.Insert(user)
	if err != nil {
		log.Printf("error! %v", err)
		return
	}
	//	log.Printf("cnt = %d\n", cnt)

	log.Printf("INSERTE userno=%d rank=%s nscore=%d pscore=%d longname=%s\n", user.Userno, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)
	return
}
