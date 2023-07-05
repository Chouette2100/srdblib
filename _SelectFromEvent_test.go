package srdblib

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/Chouette2100/exsrapi"
)
func TestSelectFromEvent(t *testing.T) {

	logfile, err :=exsrapi.CreateLogfile("SelectFromEvent_testg")
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
	log.Printf("dbconfig = %v\n", dbconfig)


	eventinf, err := SelectFromEvent("wevent", "puzzle01")

	t.Errorf("eventinf = %v", eventinf)
	log.Printf("eventinf = %v", eventinf)
	t.Errorf("err = %v", err)
	log.Printf("err = %v", err)

}