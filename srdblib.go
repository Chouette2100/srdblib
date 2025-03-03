/*
!
Copyright © 2023 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/
package srdblib

import (
	"time"

	"github.com/Chouette2100/exsrapi"
)

/*
00AA00	srdblibを導入する（データベースアクセスを一本化する）
		Event_InfにAchkを追加する（wevent用）
01AA00	SSHConfigをDBConfigに統合し、DBConfigのファイル読み込みもこの関数内で行う
01AB00	CGIでイベントリストを操作の基本にすることとSRGCEの機能強化に関わる変更
01AC00	操作対象のテーブルをsrdblib.Teventで指定する方法から関数の引数とする方法に変える
01AD00	gorpを導入(OpenDB.go)し、User.goとPoints.goを追加する。
01AD00a	InsertNewOnes.go にコメントを追加する。
01AE00	user.goにUpinsUserSetPoperty()を追加する。
01AF00	InsertNewOnes()にgorpを適用する。
01AG00	InsertNewOnes()でuserのUPDATEは行わない（userのUPDATEは定期的に行う）
01AG01	UpinsUserSetProperty()でウェイトができるようにする（srapi.ApiRoomProfileAll()の実行頻度に制限がある模様）
01AH00	srapi.ApiRoomProfileAll()の関数名をsrapi.ApiRoomProfile()としたことに対応する。
01AJ00	Event.goを追加する。User.goにwuserと関連関数を追加する。
01AK00	Timetable.goを追加する。
01AL00	GetEventsRankingByApi.goを追加する。
01AL01	GetEventsRankingByApi()の引数にmode（1: イベント開催中、2: イベント終了後）を追加する。
01AL02	GetEventsRankingByApi()でイベントが存在しない場合のエラー処理を追加する。
		InsertIntoUser()のコメントを修正する
01AM00	ギフトランキング、視聴者ギフトランキングに関する機能を追加する。
		Giftscore.go, Giftscore_test.go, srdblib.go(変更), Env.yml, Viewer.go
01AM01	InsertIntoViewerGiftScore(), InsertIntoGiftScore()の引数を変更する。
01AN01	SRGGR対応
01AN02	campaign.goを作る、Viewer.go,ViewerにOrdernoを追加する、User.go、コメントを追加する。
01AN03	campaign.goにUrlを追加する
01AN04	giftrankingにCntrblsを追加する
01AP00	UpinsViewerSetProperty()をあらたに作成する、Giftscore.goにGiftScoreCntrbをあらたに作成する
01AP01	UpinsViewerSetProperty()のバグを修正する、ViererのOrdernoを削除する
01AP02	giftscorecntrbへのinsertでのusernoの抜けを修正する
01AP03	UpdateUserSetProperty()でデータの取得に失敗したときは処理を打ち切る、そうしないとデータがクリアされてしまう
01AQ00	UpinsEventuser() 新規作成（≒ InsertNewOnes.go() ）、GetEventsRankingByApi.go() でイベントが存在しない場合の処理を追加する
01AR00	UpinsEventuser() 引数にcmapを追加する（Event_infを使わないようにするため）
01AS00	eventにthinit、thdeltaを追加する。Accesslog.goを追加する。
01AS01	srdblib.GetFeaturedEvents()のインターフェース仕様を変更する。
01AT00	Wuser, Wuserhitory, Wevent, Weventuser を User, Userhistory, Event, Eventserから定義する。
01AT01	GetFeaturedEvents()でis_bot==1のデータを除外する。
01AU00	UpinsEventuser()をジェネリックで実装したUpdateEutableG()を作成する。
01AU01	UpinsEventuserG()で同じデータでの更新を行わないようにする。
		ポインターレシーバでないものが混在するときの処理を確認する（Print()、確認の上、コメントアウトする）
01AV00	UpinsUserG.goを新たに作成する。
01AV01	UpinsUserSetProperty()で最近のデータで更新が必要なくなった場合はウェイトしないようにする。
01AW01  ユーザテーブルの更新で無駄な更新やウェイトを行わないようにする。
01AW02  ジェネリック関数をfunc .....[T userT](xu T,...)の形に統一する。
01AW03  UserTのセッターの使い方を変更する。
01AW04  UpinsUser(), UpinsEventuser()のデータベース格納値の調整を行う。
01AX00  UpinsUser()関連の関数、とくにジェネリック関数の整備を行う。
*/

const Version = "01AX00"

type Environment struct {
	//	Intervalhour int	`yaml:"Intervalhour"`
	Lmin     int `yaml:"Lmin"`
	Waitmsec int `yaml:"Waitmsec"`
}

var Env Environment = Environment{
	//	Intervalhour: 6,     //	6時間以内にデータがあれば重複チェックを行う？
	Lmin: 14400, //	viewer, user で前回更新から10日間以上経っていれば更新する(UpdateUserSetPropertyのようにこの値を使わない場合もある)
	Waitmsec: 100, //	viewer, user で新しいデータをinsertしてから1秒間待つ(APIにアクセス制限があるように思えるため)
}

/*
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
	Achk		int	//	1: ブロック、2:ボックス、子ルーム未処理のあいだはそれぞれ +4
	//	aclr		int

	//	Event_no    int
	EventStatus string //	"Over", "BeingHeld", "NotHeldYet"
	Pntbasis    int
	Ordbasis    int
	League_ids  string
	//	Status		string		//	"Confirmed":	イベント終了日翌日に確定した獲得ポイントが反映されている。
}
*/

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

var Colorlist0 []Color = []Color{
			{"#00FFFF", "#00FFFF"},
			{"#FF00FF", "#FF00FF"},
			{"#FFFF00", "#FFFF00"},
			//      -----
			{"#7F7FFF", "#7F7FFF"},
			{"#FF7F7F", "#FF7F7F"},
			{"#7FFF7F", "#7FFF7F"},

			{"#7FBFFF", "#7FBFFF"},
			{"#FF7FBF", "#FF7FBF"},
			{"#BFFF7F", "#BFFF7F"},

			{"#7FFFFF", "#7FFFFF"},
			{"#FF7FFF", "#FF7FFF"},
			{"#FFFF7F", "#FFFF7F"},

			{"#7FFFBF", "#7FFFBF"},
			{"#BF7FFF", "#BF7FFF"},
			{"#FFBF7F", "#FFBF7F"},
			//      -----
			{"#ADADFF", "#ADADFF"},
			{"#FFADAD", "#FFADAD"},
			{"#ADFFAD", "#7FFFAD"},

			{"#ADD6FF", "#ADD6FF"},
			{"#FFADD6", "#FFADD6"},
			{"#D6FFAD", "#D6FFAD"},

			{"#ADFFFF", "#ADFFFF"},
			{"#FFADFF", "#FFADFF"},
			{"#FFFFAD", "#FFFFAD"},

			{"#ADFFD6", "#ADFFD6"},
			{"#D6ADFF", "#D6ADFF"},
			{"#FFD6AD", "#FFD6AD"},
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

var Event_inf exsrapi.Event_Inf
