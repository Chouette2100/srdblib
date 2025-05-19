// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"net/http"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
)

func TestGetEventsRankingByApi(t *testing.T) {
	type args struct {
		client *http.Client
		eid    string
		mode   int
	}
	logfile, err := exsrapi.CreateLogfile("TestGetRoominfAll", "log")
	if err != nil {
		log.Printf("exsrapi.CreateLogfile() error. err=%s.\n", err.Error())
		return
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	//      データベースとの接続をオープンする。
	dbconfig, err := OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("srdblib.OpenDb() error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer Dialer.Close()
	}
	defer Db.Close()

	log.Printf("dbconfig=%v\n", dbconfig)

	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	Dbmap = &gorp.DbMap{Db: Db, Dialect: dial, ExpandSliceArgs: true}
	//	srdblib.Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
	//	srdblib.Dbmap.AddTableWithName(srdblib.Userhistory{}, "wuserhistory").SetKeys(false, "Userno", "Ts")
	//	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
	//	srdblib.Dbmap.AddTableWithName(srdblib.Eventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")
	Dbmap.AddTableWithName(Event{}, "event").SetKeys(false, "Eventid")

	//      cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//      すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	tests := []struct {
		name         string
		args         args
		wantPranking *srapi.Eventranking
		wantErr      bool
	}{
		{
			name: "tifdedebut2025_y?block_id=37020",
			args: args{
				client: client,
				eid:    "tifdedebut2025_y?block_id=37020",
				mode:   1, //	1: イベント開催中, 2: イベント終了後
			},
			wantPranking: nil,
			wantErr:      false,
		},
		// TODO: Add test cases.
		/*
			{
				name: "tifdedebut2025_s?block_id=36610",
				args: args{
					client: client,
					eid:    "tifdedebut2025_s?block_id=36610",
					mode:   1, //	1: イベント開催中, 2: イベント終了後
				},
				wantPranking: nil,
				wantErr:      false,
			},
			{
				name: "mattari_fireworks189",
				args: args{
					client: client,
					eid:    "mattari_fireworks189",
					mode:   1, //	1: イベント開催中, 2: イベント終了後
				},
				wantPranking: nil,
				wantErr:      false,
			},
				{
					name: "test_20901",
					args: args{
						client: client,
						eid:    "safaripark_showroom?block_id=20901",
						mode:	1,	//	1: イベント開催中, 2: イベント終了後
					},
					wantPranking: nil,
					wantErr:      false,
				},
				{
					name: "test_0",
					args: args{
						client: client,
						eid:    "safaripark_showroom?block_id=0",
					},
					wantPranking: nil,
					wantErr:      false,
				},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPranking, err := GetEventsRankingByApi(tt.args.client, tt.args.eid, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsRankingByApi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPranking, tt.wantPranking) {
				t.Errorf("GetEventsRankingByApi() = %v, want %v", gotPranking, tt.wantPranking)
			}
		})
	}
}
