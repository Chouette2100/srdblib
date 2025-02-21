// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"log"
	"testing"
	"time"

	"github.com/go-gorp/gorp"
	// "github.com/chouette2100/srdblib"
)

func TestUpinsEventuser(t *testing.T) {
	type args struct {
		tnow time.Time
		weu  Weventuser
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestOfUpinsEventuser",
			args: args{
				weu: Weventuser{
					EventuserBR: EventuserBR{
						Eventid: "mattari_fireworks215",
						Userno:  455133,
						Vld:     1,
						Point:   50000,
					},
					Status: 0,
				},
				tnow: time.Now().Truncate(time.Second),
			},
			wantErr: false,
		},
	}
	// データベース接続
	dbconfig, err := OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("Database error. err = %v\n", err)
		return
	}
	if dbconfig.UseSSH {
		defer Dialer.Close()
	}
	defer Db.Close()

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	Dbmap = &gorp.DbMap{Db: Db,
		Dialect:         dial,
		ExpandSliceArgs: true, //スライス引数展開オプションを有効化する
	}
	Dbmap.AddTableWithName(User{}, "user").SetKeys(false, "Userno")
	Dbmap.AddTableWithName(Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")
	// srdblib.Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
	// srdblib.Dbmap.AddTableWithName(TWuser{}, "wuser").SetKeys(false, "Userno")
	Dbmap.AddTableWithName(Wuserhistory{}, "wuserhistory").SetKeys(false, "Userno", "Ts")
	Dbmap.AddTableWithName(Event{}, "event").SetKeys(false, "Eventid")
	Dbmap.AddTableWithName(Eventuser{}, "eventuser").SetKeys(false, "Eventid", "Userno")
	Dbmap.AddTableWithName(Wevent{}, "wevent").SetKeys(false, "Eventid")
	Dbmap.AddTableWithName(Weventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpinsEventuserG(&tt.args.weu, tt.args.tnow); (err != nil) != tt.wantErr {
				t.Errorf("UpinsEventuser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
