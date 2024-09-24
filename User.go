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
0.2.0 Event.goを追加し、User.goにEvent, Eventuser, Wuserを追加する。

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

type Wuser struct {
	Userno       int
	Userid       string
	User_name    string
	Longname     string
	Shortname    string
	Genre        string
	Rank         string
	Nrank        string
	Prank        string
	Level        int
	Followers    int
	Fans         int
	Fans_lst     int
	Ts           time.Time
	Getp         string
	Graph        string
	Color        string
	Currentevent string
}

// userの履歴を保存する構造体
// PRIMARY KEY (`userno`,`ts`)
type Userhistory struct {
	Userno    int
	User_name string
	Genre     string
	Rank      string
	Nrank     string
	Prank     string
	Level     int
	Followers int
	Fans      int
	Fans_lst  int
	Ts        time.Time
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
func UpinsUserSetProperty(client *http.Client, tnow time.Time, user *User, lmin int, wait int) (
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
			//	lastrank := usert.Rank
			if usert.Ts.After(tnow.Add(time.Duration(-lmin) * time.Minute)) {
				log.Printf("skipped. UpinsUserSetProperty(userno=%d rank=%s %s)  %v", user.Userno, usert.Rank, usert.User_name, usert.Ts)
				return nil
			}
			err = UpdateUserSetProperty(client, tnow, usert)
			if err != nil {
				err = fmt.Errorf("UpdateUserSetProperty(userno=%d) returned error. %w", user.Userno, err)
			}
			//	log.Printf("UpinsUserSetProperty(userno=%d %s) lastrank=%s -> %s", user.Userno, usert.User_name, lastrank, usert.Rank)
		}
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}

	return
}

/*
テーブル user を SHOWROOMのAPI api/roomprofile を使って得られる情報で更新する。
*/
func UpdateUserSetProperty(client *http.Client, tnow time.Time, user *User) (
	err error,
) {

	lastrank := user.Rank

	//	ユーザーのランク情報を取得する
	ria, err := srapi.ApiRoomProfile(client, fmt.Sprintf("%d", user.Userno))
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
	if user.Itrank > user.Irank {
		user.Itrank = user.Irank
	}
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

	log.Printf("UPDATE(user) userno=%d rank=%s -> %s nscore=%d pscore=%d longname=%s\n", user.Userno, lastrank, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)

	//	check userhistory
	nodata := false
	uh := Userhistory{}
	sqlst := "select max(ts) ts from userhistory where userno = ? "
	err = Dbmap.SelectOne(&uh, sqlst, user.Userno)
	if err != nil {
		// log.Printf("<%s>\n", err.Error())
		if !strings.Contains(err.Error(), "sql: Scan error on column index 0, name \"ts\": unsupported Scan") {
			//	検索条件に該当するデータが一件もない
			err = fmt.Errorf("Dbmap.SelectOne error: %v", err)
			log.Printf("error: %v", err)
			return err
		}
		nodata = true
	}

	intf := interface{}(nil)
	if !nodata {
		//	userhisotry にルームのデータが存在するので、そのデータを取得する
		intf, err = Dbmap.Get(Userhistory{}, user.Userno, uh.Ts)
		if err != nil {
			err = fmt.Errorf("Get(userno=%d) returned error. %w", user.Userno, err)
			return err
		}
	}

	if intf == nil {
		//	userhisotryにデータが存在しなかった場合
		uh := &Userhistory{
			Userno:    user.Userno,
			User_name: user.User_name,
			Genre:     user.Genre,
			Rank:      user.Rank,
			Nrank:     user.Nrank,
			Prank:     user.Prank,
			Level:     user.Level,
			Followers: user.Followers,
			Fans:      user.Fans,
			Fans_lst:  user.Fans_lst,
			Ts:        tnow,
		}

		err = Dbmap.Insert(uh)
		if err != nil {
			err = fmt.Errorf("Insert(userhistory,userno=%d) returned error. %w", uh.Userno, err)
			return
		}
		log.Printf("INSERT(userhistory) userno=%d name =%s genre= %s rank=%s level=%d fikkiwers=%d\n",
			 uh.Userno, uh.User_name, uh.Genre, uh.Rank, uh.Level, uh.Followers)

	} else {
		//	userhisotryにデータがすでに存在するとき
		uh := intf.(*Userhistory)
		if tnow.Sub(uh.Ts) > time.Duration(Env.Lmin) * time.Minute {
			//	最後のデータから一定時間過ぎているときは新しいデータを挿入する
			uh.User_name = user.User_name
			uh.Genre = user.Genre
			uh.Rank = user.Rank
			uh.Nrank = user.Nrank
			uh.Prank = user.Prank
			uh.Level = user.Level
			uh.Followers = user.Followers
			uh.Fans = user.Fans
			uh.Fans_lst = user.Fans_lst
			uh.Ts = tnow

			err = Dbmap.Insert(uh)
			if err != nil {
				err = fmt.Errorf("Insert(userhistory, userno=%d) returned error. %w", uh.Userno, err)
				return
			}
			log.Printf("INSERT(userhistory) userno=%d name =%s genre= %s rank=%s level=%d fikkiwers=%d\n",
			uh.Userno, uh.User_name, uh.Genre, uh.Rank, uh.Level, uh.Followers)
		}

	}

	return
}

/*
テーブル user に新しいデータを追加する
*/
func InsertIntoUser(client *http.Client, tnow time.Time, userno int) (
	err error,
) {

	//	ユーザーのランク情報を取得する
	ria, err := srapi.ApiRoomProfile(client, fmt.Sprintf("%d", userno))
	if err != nil {
		err = fmt.Errorf("ApiRoomProfile(%d) returned error. %w", userno, err)
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
	user.Itrank = user.Irank
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

	log.Printf("INSERT(user) userno=%d rank=%s nscore=%d pscore=%d longname=%s\n", user.Userno, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)

	// userhistory

	uh := &Userhistory{
		Userno:    user.Userno,
		User_name: user.User_name,
		Genre:     user.Genre,
		Rank:      user.Rank,
		Nrank:     user.Nrank,
		Prank:     user.Prank,
		Level:     user.Level,
		Followers: user.Followers,
		Fans:      user.Fans,
		Fans_lst:  user.Fans_lst,
		Ts:        tnow,
	}

	//	cnt, err = Dbmap.Update(uh)

	err = Dbmap.Insert(uh)
	if err != nil {
		log.Printf("error! %v", err)
		return
	}

	log.Printf("INSERT(userhistory) userno=%d rank=%s nscore=%d pscore=%d longname=%s\n", user.Userno, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)

	return
}

/*
ルーム番号 user.Userno が テーブル user に存在しないときは新しいデータを挿入し、存在するときは 既存のデータを更新する。
*/
func UpinsWuserSetProperty(client *http.Client, tnow time.Time, wuser *Wuser, lmin int, wait int) (
	err error,
) {

	row, err := Dbmap.Get(Wuser{}, wuser.Userno)
	if err != nil {
		err = fmt.Errorf("Get(userno=%d) returned error. %w", wuser.Userno, err)
		return err
	} else {
		if row == nil {
			err = InsertIntoWuser(client, tnow, wuser.Userno)
			if err != nil {
				err = fmt.Errorf("InsertIntoUser(userno=%d) returned error. %w", wuser.Userno, err)
			}
		} else {
			wusert := row.(*Wuser)
			//	lastrank := usert.Rank
			if wusert.Ts.After(tnow.Add(time.Duration(-lmin) * time.Minute)) {
				log.Printf("skipped. UpinsUserSetProperty(userno=%d rank=%s %s)  %v", wuser.Userno, wusert.Rank, wusert.User_name, wusert.Ts)
				return nil
			}
			err = UpdateWuserSetProperty(client, tnow, wusert)
			if err != nil {
				err = fmt.Errorf("UpdateUserSetProperty(userno=%d) returned error. %w", wuser.Userno, err)
			}
			//	log.Printf("UpinsUserSetProperty(userno=%d %s) lastrank=%s -> %s", user.Userno, usert.User_name, lastrank, usert.Rank)
		}
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}

	return
}

/*
テーブル user を SHOWROOMのAPI api/roomprofile を使って得られる情報で更新する。
*/
func UpdateWuserSetProperty(client *http.Client, tnow time.Time, wuser *Wuser) (
	err error,
) {

	lastrank := wuser.Rank

	//	ユーザーのランク情報を取得する
	ria, err := srapi.ApiRoomProfile(client, fmt.Sprintf("%d", wuser.Userno))
	if err != nil {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %w", wuser.Userno, err)
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
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned empty.ShowRankSubdivided", wuser.Userno)
		return err
	}

	//	user.Userno =	キー
	wuser.Userid = ria.RoomURLKey
	wuser.User_name = ria.RoomName
	//	user.Longname =		新規作成時に設定され、変更はSRCGIで行う
	//	user.Shortname =　　〃
	wuser.Genre = ria.GenreName
	//	wuser.GenreID = ria.GenreID
	wuser.Rank = ria.ShowRankSubdivided
	wuser.Nrank = humanize.Comma(int64(ria.NextScore))
	wuser.Prank = humanize.Comma(int64(ria.PrevScore))
	//	wuser.Irank = MakeSortKeyOfRank(ria.ShowRankSubdivided, ria.NextScore)
	//	wuser.Inrank = ria.NextScore
	//	wuser.Iprank = ria.PrevScore
	//	if wuser.Itrank > user.Irank {
	//		user.Itrank = user.Irank
	//	}
	wuser.Level = ria.RoomLevel
	wuser.Followers = ria.FollowerNum

	pafr, err := srapi.ApiActivefanRoom(client, strconv.Itoa(wuser.Userno), tnow.Format("200601"))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", wuser.Userno, err)
		return err
	}
	wuser.Fans = pafr.TotalUserCount
	//	wuser.FanPower = pafr.FanPower
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(wuser.Userno), fmt.Sprintf("%04d%02d", yy, mm))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", wuser.Userno, err)
		return err
	}
	wuser.Fans_lst = pafr.TotalUserCount
	//	wuser.FanPower_lst = pafr.FanPower

	wuser.Ts = tnow
	//	user.Getp =		ここから三つのフィールドは現在使われていない。
	//	user.Graph =	default: ''
	//	user.Color =
	eurl := ria.Event.URL
	eurla := strings.Split(eurl, "/")
	wuser.Currentevent = eurla[len(eurla)-1]

	//	cnt, err := Dbmap.Update(user)
	_, err = Dbmap.Update(wuser)
	if err != nil {
		log.Printf("error! %v", err)
		return
	}
	//	log.Printf("cnt = %d\n", cnt)

	log.Printf("UPDATE userno=%d rank=%s -> %s nscore=%d pscore=%d longname=%s\n", wuser.Userno, lastrank, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)
	return
}

/*
テーブル user に新しいデータを追加する
*/
func InsertIntoWuser(client *http.Client, tnow time.Time, userno int) (
	err error,
) {

	//	ユーザーのランク情報を取得する
	ria, err := srapi.ApiRoomProfile(client, fmt.Sprintf("%d", userno))
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

	wuser := new(Wuser)

	wuser.Userno = userno
	wuser.Userid = ria.RoomURLKey
	wuser.User_name = ria.RoomName
	wuser.Longname = ria.RoomName
	wuser.Shortname = strconv.Itoa(userno % 100)
	wuser.Genre = ria.GenreName
	//	wuser.GenreID = ria.GenreID
	wuser.Rank = ria.ShowRankSubdivided
	wuser.Nrank = humanize.Comma(int64(ria.NextScore))
	wuser.Prank = humanize.Comma(int64(ria.PrevScore))
	//	wuser.Irank = MakeSortKeyOfRank(ria.ShowRankSubdivided, ria.NextScore)
	//	wuser.Inrank = ria.NextScore
	//	wuser.Iprank = ria.PrevScore
	//	wuser.Itrank = user.Irank
	wuser.Level = ria.RoomLevel
	wuser.Followers = ria.FollowerNum

	pafr, err := srapi.ApiActivefanRoom(client, strconv.Itoa(wuser.Userno), tnow.Format("200601"))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", wuser.Userno, err)
		return err
	}
	wuser.Fans = pafr.TotalUserCount
	//	wuser.FanPower = pafr.FanPower
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(wuser.Userno), fmt.Sprintf("%04d%02d", yy, mm))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", wuser.Userno, err)
		return err
	}
	wuser.Fans_lst = pafr.TotalUserCount
	//	wuser.FanPower_lst = pafr.FanPower

	wuser.Ts = tnow
	wuser.Getp = ""
	wuser.Graph = ""
	wuser.Color = ""
	eurl := ria.Event.URL
	eurla := strings.Split(eurl, "/")
	wuser.Currentevent = eurla[len(eurla)-1]

	//	cnt, err := Dbmap.Update(user)
	err = Dbmap.Insert(wuser)
	if err != nil {
		log.Printf("error! %v", err)
		return
	}
	//	log.Printf("cnt = %d\n", cnt)

	log.Printf("INSERT userno=%d rank=%s nscore=%d pscore=%d longname=%s\n", wuser.Userno, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)
	return
}
