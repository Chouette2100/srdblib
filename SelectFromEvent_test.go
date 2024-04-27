package srdblib

import (
	"log"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Chouette2100/exsrapi"
)

func TestSelectFromEvent(t *testing.T) {
	type args struct {
		tevent string
		eventid string
	}

	ts, _ := time.Parse("2006-01-02 15:04:05 MST", "2023-06-12 18:00:00 JST")
	te, _ := time.Parse("2006-01-02 15:04:05 MST", "2023-06-21 21:59:59 JST")
	tests := []struct {
		name          string
		args          args
		wantPeventinf *exsrapi.Event_Inf
		wantErr       bool
	}{
		{
			name: "test1",
			args: args{
				tevent: "wevent",
				eventid: "puzzle01",
			},
			wantPeventinf: &exsrapi.Event_Inf{
				Event_ID:    "puzzle01",
				I_Event_ID:  32653,
				Event_name:  "あなたの顔がパズルに！？オリジナルジグソーパズル争奪戦！！",
				Period:      "",
				Dperiod:     10,
				Start_time:  ts,
				Sstart_time: "",
				Start_date:  19519.625,
				End_time:    te,
				Send_time:   "",
				NoEntry:     0,
				NoRoom:      0,
				Intervalmin: 0,
				Modmin:      0,
				Modsec:      0,
				Fromorder:   0,
				Toorder:     0,
				Resethh:     0,
				Resetmm:     0,
				Nobasis:     0,
				Maxdsp:      0,
				Cmap:        0,
				Target:      0,
				Rstatus:     "NowSaved",
				Maxpoint:    0,
				MaxPoint:    0,
				Gscale:      0,
				Achk:        0,
				Aclr:        0,
				EventStatus: "",
				Pntbasis:    0,
				Ordbasis:    0,
				League_ids:  "",
			},

			wantErr: false,
		},
		// TODO: Add test cases.
	}

	//	Tevent = "wevent"

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

	log.Printf("dbconfig = %+v\n", dbconfig)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPeventinf, err := SelectFromEvent(tt.args.tevent, tt.args.eventid)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectFromEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPeventinf, tt.wantPeventinf) {
				t.Errorf("SelectFromEvent() = %v, want %v", gotPeventinf, tt.wantPeventinf)
			}
		})
	}
}
