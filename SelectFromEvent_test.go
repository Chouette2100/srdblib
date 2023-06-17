package srdblib
import (
	"log"
	"os"
	"io"
	"testing"
	"github.com/Chouette2100/exsrapi"
)
func TestSelectFromEvent(t *testing.T) {

	log.SetOutput(io.MultiWriter(os.Stdout))

	var dbconfig DBConfig
	err := exsrapi.LoadConfig("DBConfig.yml", &dbconfig)
	if err != nil {
		t.Errorf("err=%s.\n", err.Error())
		log.Printf("err=%s.\n", err.Error())
		os.Exit(1)
	}
	t.Errorf("%+v\n", dbconfig)

	err = OpenDb(&dbconfig)
	if err != nil {
		t.Errorf("Database error. err = %v\n", err)
		log.Printf("Database error. err = %v\n", err)
		return
	}
	if dbconfig.UseSSH {
		defer Dialer.Close()
	}
	defer Db.Close()


	eventinf, err := SelectFromEvent("bestofhawaiianwedding2023_3_2a")

	t.Errorf("eventinf = %v", eventinf)
	log.Printf("eventinf = %v", eventinf)
	t.Errorf("err = %v", err)
	log.Printf("err = %v", err)

}