// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/exsrapi"
)

func TestGetFeaturedEvents(t *testing.T) {
	type args struct {
		hours int
		num   int
	}
	tests := []struct {
		name       string
		args       args
		wantEvents map[string]bool
	}{
		{
			name: "SelectFromEvent_test",
			args: args{
				hours: 24,
				num: 20,
			},
			wantEvents: map[string]bool{
				"SelectFromEvent_test": true,
			},
		},
		// TODO: Add test cases.
	}

	logfile, err := exsrapi.CreateLogfile("SelectFromEvent_testg")
	if err != nil {
		t.Errorf("logfile error. err = %v\n", err)
		return
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	dbconfig, err := OpenDb("DBConfig.yml")
	if err != nil {
		t.Errorf("Database error. err = %v\n", err)
		log.Printf("Database error. err = %v\n", err)
		return
	}
	if dbconfig.UseSSH {
		defer Dialer.Close()
	}
	defer Db.Close()

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	Dbmap = &gorp.DbMap{Db: Db, Dialect: dial, ExpandSliceArgs: true}

	log.Printf("dbconfig = %+v\n", dbconfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEvents := GetFeaturedEvents(tt.args.hours, tt.args.num); !reflect.DeepEqual(gotEvents, tt.wantEvents) {
				t.Errorf("GetFeaturedEvents() = %v, want %v", gotEvents, tt.wantEvents)
			} else {
				t.Logf("GetFeaturedEvents() = %v, want %v", gotEvents, tt.wantEvents)
			}
		})
	}
}
