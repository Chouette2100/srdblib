/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/
package srdblib

import (


	"time"

)

/*
Ver.00AA00	srdblibを導入する（データベースアクセスを一本化する）
*/

const Version = "00AA00"

type Event_Inf struct {
	Event_ID    string
	I_Event_ID  int
	Event_name  string
	Period      string
	Dperiod     float64
	Start_time  time.Time
	Sstart_time string
	Start_date  float64
	End_time    time.Time
	Send_time   string
	NoEntry     int
	NoRoom      int //	ルーム数
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
	MaxPoint    int	//	DBには該当するものはない
	Gscale    int	//	DBのMaxpoint = 構造体の Maxpoint + Gscale
	//	achk		int
	//	aclr		int

	//	Event_no    int
	EventStatus string //	"Over", "BeingHeld", "NotHeldYet"
	Pntbasis    int
	Ordbasis    int
	League_ids  string
	//	Status		string		//	"Confirmed":	イベント終了日翌日に確定した獲得ポイントが反映されている。
}

type Color struct {
	Name  string
	Value string
}

// https://www.fukushihoken.metro.tokyo.lg.jp/kiban/machizukuri/kanren/color.files/colorudguideline.pdf
var Colorlist2 []Color = []Color{
	{"red", "#FF2800"},
	{"yellow", "#FAF500"},
	{"green", "#35A16B"},
	{"blue", "#0041FF"},
	{"skyblue", "#66CCFF"},
	{"lightpink", "#FFD1D1"},
	{"orange", "#FF9900"},
	{"purple", "#9A0079"},
	{"brown", "#663300"},
	{"lightgreen", "#87D7B0"},
	{"white", "#FFFFFF"},
	{"gray", "#77878F"},
}

var Colorlist1 []Color = []Color{
	{"cyan", "cyan"},
	{"magenta", "magenta"},
	{"yellow", "yellow"},
	{"royalblue", "royalblue"},
	{"coral", "coral"},
	{"khaki", "khaki"},
	{"deepskyblue", "deepskyblue"},
	{"crimson", "crimson"},
	{"orange", "orange"},
	{"lightsteelblue", "lightsteelblue"},
	{"pink", "pink"},
	{"sienna", "sienna"},
	{"springgreen", "springgreen"},
	{"blueviolet", "blueviolet"},
	{"salmon", "salmon"},
	{"lime", "lime"},
	{"red", "red"},
	{"darkorange", "darkorange"},
	{"skyblue", "skyblue"},
	{"lightpink", "lightpink"},
}

type ColorInf struct {
	Color      string
	Colorvalue string
	Selected   string
}

type ColorInfList []ColorInf

type RoomInfo struct {
	Name      string //	ルーム名のリスト
	Longname  string
	Shortname string
	Account   string //	アカウントのリスト、アカウントは配信のURLの最後の部分の英数字です。
	ID        string //	IDのリスト、IDはプロフィールのURLの最後の部分で5～6桁の数字です。
	Userno    int
	//	APIで取得できるデータ(1)
	Genre      string
	Rank       string
	Irank      int
	Nrank      string
	Prank      string
	Followers  int
	Sfollowers string
	Fans       int
	Fans_lst   int
	Level      int
	Slevel     string
	//	APIで取得できるデータ(2)
	Order        int
	Point        int //	イベント終了後12時間〜36時間はイベントページから取得できることもある
	Spoint       string
	Istarget     string
	Graph        string
	Iscntrbpoint string
	Color        string
	Colorvalue   string
	Colorinflist ColorInfList
	Formid       string
	Eventid      string
	Status       string
	Statuscolor  string
}

type RoomInfoList []RoomInfo

// sort.Sort()のための関数三つ
func (r RoomInfoList) Len() int {
	return len(r)
}

func (r RoomInfoList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RoomInfoList) Choose(from, to int) (s RoomInfoList) {
	s = r[from:to]
	return
}

var SortByFollowers bool

// 降順に並べる
func (r RoomInfoList) Less(i, j int) bool {
	//	return e[i].point < e[j].point
	if SortByFollowers {
		return r[i].Followers > r[j].Followers
	} else {
		return r[i].Point > r[j].Point
	}
}

type PerSlot struct {
	Timestart time.Time
	Dstart    string
	Tstart    string
	Tend      string
	Point     string
	Ipoint    int
	Tpoint    string
}

type PerSlotInf struct {
	Eventname   string
	Eventid     string
	Period      string
	Roomname    string
	Roomid      int
	Perslotlist []PerSlot
}

var Event_inf Event_Inf


