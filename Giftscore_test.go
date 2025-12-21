// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func TestInserIntoGiftScore(t *testing.T) {
	type args struct {
		client *http.Client
		dbmap  *gorp.DbMap
		giftid int
		cgr    *srapi.CdnGiftRanking
		idx    int
		tnow   time.Time
	}

	fileenv := "Env.yml"
	err := exsrapi.LoadConfig(fileenv, &srdblib.Env)
	if err != nil {
		err = fmt.Errorf("exsrapi.Loadconfig(): %w", err)
		log.Printf("%s\n", err.Error())
		return
	}

	logfile, err := exsrapi.CreateLogfile("TestGetRoominfAll", "log")
	if err != nil {
		log.Printf("exsrapi.CreateLogfile() error. err=%s.\n", err.Error())
		return
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	//      データベースとの接続をオープンする。
	dbconfig, err := srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("srdblib.OpenDb() error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}
	defer srdblib.Db.Close()
	log.Printf("dbconfig=%v\n", dbconfig)

	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db, Dialect: dial, ExpandSliceArgs: true}
	/*
		//	srdblib.Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
		//	srdblib.Dbmap.AddTableWithName(srdblib.Userhistory{}, "wuserhistory").SetKeys(false, "Userno", "Ts")
		//	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
		//	srdblib.Dbmap.AddTableWithName(srdblib.Eventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")
		//	Dbmap.AddTableWithName(Event{}, "event").SetKeys(false, "Eventid")
		Dbmap.AddTableWithName(User{}, "user").SetKeys(false, "Userno")
		Dbmap.AddTableWithName(GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
		Dbmap.AddTableWithName(Viewer{}, "viewer").SetKeys(false, "Viewerid")
		Dbmap.AddTableWithName(ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	*/

	srdblib.AddTableWithName()

	//      cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//      すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	cgr := &srapi.CdnGiftRanking{
		TotalScore: 1,
		RankingList: []srapi.GrRanking{
			{
				RoomID:  307073,
				OrderNo: 1,
				Score:   3823,
				Room: srapi.GrRoom{
					Name:   "xxxxxx",
					URLKey: "yyyyy",
				},
			},
		},
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				client: client,
				dbmap:  srdblib.Dbmap,
				giftid: 1,
				cgr:    cgr,
				idx:    1,
				tnow:   time.Now(),
			},
			wantErr: false,
			// TODO: Add test cases.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("%+v\n", tt.args)
			if err := srdblib.InserIntoGiftScore(tt.args.client, tt.args.dbmap, tt.args.giftid, &tt.args.cgr.RankingList[0], tt.args.tnow); (err != nil) != tt.wantErr {
				t.Errorf("InserIntoGiftScore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInserIntoViewerGiftScore(t *testing.T) {
	type args struct {
		client *http.Client
		dbmap  *gorp.DbMap
		giftid int
		cugr   *srapi.CdnUserGiftRanking
		idx    int
		tnow   time.Time
	}
	logfile, err := exsrapi.CreateLogfile("TestGetRoominfAll", "log")
	if err != nil {
		log.Printf("exsrapi.CreateLogfile() error. err=%s.\n", err.Error())
		return
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	//      データベースとの接続をオープンする。
	dbconfig, err := srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("srdblib.OpenDb() error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}
	defer srdblib.Db.Close()
	log.Printf("dbconfig=%v\n", dbconfig)

	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db, Dialect: dial, ExpandSliceArgs: true}
	/*
		//	Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
		//	Dbmap.AddTableWithName(srdblib.Userhistory{}, "wuserhistory").SetKeys(false, "Userno", "Ts")
		//	Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
		//	Dbmap.AddTableWithName(srdblib.Eventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")
		//	Dbmap.AddTableWithName(Event{}, "event").SetKeys(false, "Eventid")
		Dbmap.AddTableWithName(User{}, "user").SetKeys(false, "Userno")
		Dbmap.AddTableWithName(GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
		Dbmap.AddTableWithName(Viewer{}, "viewer").SetKeys(false, "Viewerid")
		Dbmap.AddTableWithName(ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")
		Dbmap.AddTableWithName(ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	*/

	srdblib.AddTableWithName()

	//      cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//      すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	cugr := &srapi.CdnUserGiftRanking{
		RankingList: []srapi.UgrRanking{
			{
				UserID:  785616,
				OrderNo: 1,
				Score:   141,
				User: srapi.UgrUser{
					Name: "g͙a͙k͙u͙🏖Aϖϖϖϖϖa🌈ラ王、ア王ギフトお願い",
				},
			},
		},
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				client: client,
				dbmap:  srdblib.Dbmap,
				giftid: 206,
				cugr:   cugr,
				idx:    0,
				tnow:   time.Now(),
			},
			wantErr: false,
			// TODO: Add test cases.
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("%+v\n", tt.args)
			if err := srdblib.InserIntoViewerGiftScore(tt.args.client, tt.args.dbmap, tt.args.giftid, &tt.args.cugr.RankingList[tt.args.idx], tt.args.tnow); (err != nil) != tt.wantErr {
				t.Errorf("InserIntoViewerGiftScore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
