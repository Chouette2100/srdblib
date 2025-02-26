// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
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
	"github.com/jinzhu/copier"

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

type TWuser struct {
	Userno    int
	Userid    string
	User_name string
	Longname  string
	Shortname string
	Genre     string
	// GenreID      int
	Rank  string
	Nrank string
	Prank string
	// Irank        int
	// Inrank       int
	// Iprank       int
	// Itrank       int
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

// UserT is an interface for User
type UserT interface {
	Get() (*User, error)
	Set(*User) error
}

// Getter and Setter for User
func (u *User) Get() (
	result *User,
	err error,
) {
	result = u
	return
}

func (u *User) Set(nu *User) (err error) {
	*u = *nu
	return nil
}

// Getter and Setter for TWuser
func (twu *TWuser) Get() (
	result *User,
	err error,
) {
	result = new(User)
	err = copier.Copy(result, twu)
	if err != nil {
		err = fmt.Errorf("copier.Copy failed: %w", err)
		return
	}
	return
}

func (twu *TWuser) Set(nu *User) (err error) {
	err = copier.Copy(twu, nu)
	if err != nil {
		err = fmt.Errorf("copier.Copy failed: %w", err)
	}
	return
}

// userとtwuserのデータをデータベースから取得し、新しいデータを返す
func GetLastUserdata[T UserT](
	xuser T,
	lmin int,
) (
	estatus int, // 0: user,wuserテーブルにデータない 1: データがあるが古い、 2: userに新しいデータがある 3: wuserに新しいデータがある
	vdata *User,
	err error,
) {

	var cuser *User
	cuser, _ = xuser.Get()
	userno := cuser.Userno

	var user1, user2 T
	switch any(xuser).(type) {
	case User:
		user1 = any(new(User)).(T)
		user2 = any(new(TWuser)).(T)

	case TWuser:
		user1 = any(new(TWuser)).(T)
		user2 = any(new(User)).(T)
	}

	intf, err := Dbmap.Get(user1, userno)
	if err != nil {
		// userテーブルのデータが取得できない
		err = fmt.Errorf("Get(%d): database access error", userno)
		return
	} else if intf == nil {
		// userテーブルのデータが存在しない
		estatus = 0
	} else {
		// userテーブルのデータを仮の戻り値とする
		user1 = any(intf).(T)
		vdata, _ = user1.Get()
		if vdata.Ts.After(time.Now().Add(time.Duration(-lmin) * time.Minute)) {
			estatus = 2
		} else {
			estatus = 1
		}
	}

	intf, err = Dbmap.Get(user2, userno)
	if err != nil {
		// userテーブルのデータが取得できない
		err = fmt.Errorf("Get(%d): database access error", userno)
		return
	} else if intf != nil {
		user2 = any(intf).(T)
		tdata, _ := user2.Get()
		if estatus == 0 || tdata.Ts.After(vdata.Ts) {
			vdata = tdata
			if vdata.Ts.After(time.Now().Add(time.Duration(-lmin) * time.Minute)) {
				estatus = 3
			} else {
				estatus = 1
			}
		}
	}
	return
}

// ルーム番号 user.Userno が テーブル user に存在しないときはApiでユーザ情報を取得し新しいデータを挿入し、
// 存在するときは そのデータがlimin分より古ければAPIでユーザ情報取得して既存のデータを更新する。
// xuser は xuser.Userno のみが使用される。
func UpinsUser[T UserT](
	client *http.Client,
	tnow time.Time,
	xuser *T, // xuser.Userno が使用される。更新対象がuserであるかtwuserであるかを判定するためにxuserが使用される。
	lmin int, // DBのデータがlmin分より古いときはAPIでデータを取得して、内容が現在と異なるときは更新する。
	wait int, // APIでデータを取得したときは（アクセス制限を受けないように）waitミリ秒待つ。
) (
	rxuser *T,
	err error,
) {

	// xuser.Userno を取得する。
	var user *User
	user, err = (*xuser).Get()
	if err != nil {
		err = fmt.Errorf("failed to get user: %w", err)
		return
	}
	userno := user.Userno

	// テーブル xuser のデータを取得する。
	rxuser, err = SelectUserdata(xuser, userno)
	if err != nil {
		// データベースエラー
		err = fmt.Errorf("Get(xuser=%+v) returned error. %w", xuser, err)
		return
	} else {
		if rxuser == nil {
			// テーブル xuser にデータが存在しないのでAPIでデータを取得して挿入する。
			rxuser, err = InsertUsertable(client, tnow, wait, xuser)
			if err != nil {
				err = fmt.Errorf("InsertIntoUser(userno=%d) returned error. %w", userno, err)
			}
		} else {
			// テーブル xuser にデータが存在するので、データが古いかどうかを判定してAPIでデータを取得して更新する。
			ruser, _ := (*rxuser).Get() // rxuser.Get() とするとコンパイラが型推論できないらしい。
			//	lastrank := usert.Rank
			if ruser.Ts.After(tnow.Add(time.Duration(-lmin) * time.Minute)) {
				// データが古くないので更新しない。
				log.Printf("skipped. UpinsUser(userno=%d rank=%s %s)  %v",
					ruser.Userno, ruser.Rank, ruser.User_name, ruser.Ts)
				return
			}
			// APIでデータを取得し必要に応じて更新する。
			rxuser, err = UpdateUsertable(client, tnow, wait, xuser)
			if err != nil {
				err = fmt.Errorf("UpdateUserSetProperty(userno=%d) returned error. %w", ruser.Userno, err)
			}
			//	log.Printf("UpinsUserSetProperty(userno=%d %s) lastrank=%s -> %s", user.Userno, usert.User_name, lastrank, usert.Rank)
		}
	}

	return
}

// テーブル user からデータを取得する。
func SelectUserdata[T UserT](xu *T, userno int) (
	result *T,
	err error,
) {

	var intf interface{}

	intf, err = Dbmap.Get(xu, userno)
	if err != nil {
		err = fmt.Errorf("Dbmap.Get failed: %w", err)
		return
	} else if intf == nil {
		result = nil
		return
	} else {
		p := intf.(T) // result = intf.(*T) とするとコンパイラが型推論を間違うような。
		result = &p
	}
	return
}

/*
テーブル user を SHOWROOMのAPI api/roomprofile を使って得られる情報で更新する。
func UpdateUsertable[T UserT](xu *T) (err error) {

	var nr int64
	nr, err = Dbmap.Update(xu)
	if err != nil {
		err = fmt.Errorf("Dbmap.Update failed: %w", err)
	} else if nr != 1 {
		err = fmt.Errorf("Dbmap.Update failed: nr = %d", nr)
	}
	return
}
*/

// SHOWROOMのAPI api/roomprofile を使って得られる情報でユーザテーブルを更新する。
func UpdateUsertable[T UserT](client *http.Client, tnow time.Time, wait int, xuser *T) (
	rxuser *T,
	err error,
) {

	tuser, _ := (*xuser).Get()

	lastrank := tuser.Rank

	//	ユーザーのランク情報を取得する
	var ria *srapi.RoomInfAll
	ria, err = srapi.ApiRoomProfile(client, fmt.Sprintf("%d", tuser.Userno))
	if err != nil {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %w", tuser.Userno, err)
		return
	}
	if ria.Errors != nil {
		//	err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %v", userno, ria.Errors)
		//	return err
		ria.ShowRankSubdivided = "unknown"
		ria.NextScore = 0
		ria.PrevScore = 0
	}

	if ria.ShowRankSubdivided == "" {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned empty.ShowRankSubdivided", tuser.Userno)
		return
	}

	//	user.Userno =	キー
	tuser.Userid = ria.RoomURLKey
	tuser.User_name = ria.RoomName
	//	user.Longname =		新規作成時に設定され、変更はSRCGIで行う
	//	user.Shortname =　　〃
	tuser.Genre = ria.GenreName
	//	wuser.GenreID = ria.GenreID
	tuser.Rank = ria.ShowRankSubdivided
	tuser.Nrank = humanize.Comma(int64(ria.NextScore))
	tuser.Prank = humanize.Comma(int64(ria.PrevScore))
	//	wuser.Irank = MakeSortKeyOfRank(ria.ShowRankSubdivided, ria.NextScore)
	//	wuser.Inrank = ria.NextScore
	//	wuser.Iprank = ria.PrevScore
	//	if wuser.Itrank > user.Irank {
	//		user.Itrank = user.Irank
	//	}
	tuser.Level = ria.RoomLevel
	tuser.Followers = ria.FollowerNum

	var pafr *srapi.ActivefanRoom
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(tuser.Userno), tnow.Format("200601"))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", tuser.Userno, err)
		return
	}
	tuser.Fans = pafr.TotalUserCount
	//	wuser.FanPower = pafr.FanPower
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(tuser.Userno), fmt.Sprintf("%04d%02d", yy, mm))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", tuser.Userno, err)
		return
	}
	tuser.Fans_lst = pafr.TotalUserCount
	//	wuser.FanPower_lst = pafr.FanPower

	tuser.Ts = tnow
	//	user.Getp =		ここから三つのフィールドは現在使われていない。
	//	user.Graph =	default: ''
	//	user.Color =
	eurl := ria.Event.URL
	eurla := strings.Split(eurl, "/")
	tuser.Currentevent = eurla[len(eurla)-1]

	//	cnt, err := Dbmap.Update(user)
	// rxuser = new(T)
	// nuser.Set(tuser)    NG  コンパイラが型推論できない
	// (UserT(*nuser)).Set(tuser) OK?
	(*xuser).Set(tuser)
	rxuser = xuser
	_, err = Dbmap.Update(*xuser)
	if err != nil {
		log.Printf("error! %v", err)
		return
	}
	//	log.Printf("cnt = %d\n", cnt)

	time.Sleep(time.Duration(wait) * time.Millisecond)

	log.Printf("UPDATE userno=%d rank=%s -> %s nscore=%d pscore=%d longname=%s\n",
		tuser.Userno, lastrank, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)
	return
}

// APIで UserT.Userno の ユーザ情報を取得し、ユーザテーブルに新しいデータを挿入する。
func InsertUsertable[T UserT](
	client *http.Client,
	tnow time.Time,
	wait int, // srapi.RoomInfAll()実行後の待ち時間(ms)
	xuser *T,
) (
	rxuser *T,
	err error,
) {

	var user *User
	user, err = (*xuser).Get()
	if err != nil {
		err = fmt.Errorf("failed to get user: %w", err)
		return
	}
	userno := user.Userno

	//	ユーザーのランク情報を取得する
	var ria *srapi.RoomInfAll
	ria, err = srapi.ApiRoomProfile(client, fmt.Sprintf("%d", userno))
	if err != nil {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %w", userno, err)
		return
	}
	time.Sleep(time.Duration(wait) * time.Millisecond)
	if ria.Errors != nil {
		//	err = fmt.Errorf("ApiRoomProfile_All(%d) returned error. %v", userno, ria.Errors)
		//	return err
		ria.ShowRankSubdivided = "unknown"
		ria.NextScore = 0
		ria.PrevScore = 0
	}

	if ria.ShowRankSubdivided == "" {
		err = fmt.Errorf("ApiRoomProfile_All(%d) returned empty.ShowRankSubdivided", userno)
		return
	}

	user.Userid = ria.RoomURLKey
	user.User_name = ria.RoomName
	user.Longname = ria.RoomName
	user.Shortname = strconv.Itoa(userno % 100)
	user.Genre = ria.GenreName
	//	wuser.GenreID = ria.GenreID
	user.Rank = ria.ShowRankSubdivided
	user.Nrank = humanize.Comma(int64(ria.NextScore))
	user.Prank = humanize.Comma(int64(ria.PrevScore))
	//	wuser.Irank = MakeSortKeyOfRank(ria.ShowRankSubdivided, ria.NextScore)
	//	wuser.Inrank = ria.NextScore
	//	wuser.Iprank = ria.PrevScore
	//	wuser.Itrank = user.Irank
	user.Level = ria.RoomLevel
	user.Followers = ria.FollowerNum

	var pafr *srapi.ActivefanRoom
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(user.Userno), tnow.Format("200601"))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", user.Userno, err)
		return
	}
	user.Fans = pafr.TotalUserCount
	//	wuser.FanPower = pafr.FanPower
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	pafr, err = srapi.ApiActivefanRoom(client, strconv.Itoa(user.Userno), fmt.Sprintf("%04d%02d", yy, mm))
	if err != nil {
		err = fmt.Errorf("ApiActivefanRoom(%d) returned error. %w", user.Userno, err)
		return
	}
	user.Fans_lst = pafr.TotalUserCount
	//	wuser.FanPower_lst = pafr.FanPower

	user.Ts = tnow
	user.Getp = ""
	user.Graph = ""
	user.Color = ""
	eurl := ria.Event.URL
	eurla := strings.Split(eurl, "/")
	user.Currentevent = eurla[len(eurla)-1]

	//	cnt, err := Dbmap.Update(user)
	rxuser = new(T)
	(*rxuser).Set(user)
	err = Dbmap.Insert(rxuser)
	if err != nil {
		log.Printf("error! %v", err)
		return
	}
	//	log.Printf("cnt = %d\n", cnt)

	log.Printf("INSERT userno=%d rank=%s nscore=%d pscore=%d longname=%s\n",
		userno, ria.ShowRankSubdivided, ria.NextScore, ria.PrevScore, ria.RoomName)
	return
}
