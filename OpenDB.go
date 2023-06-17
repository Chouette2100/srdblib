/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/

package srdblib

import (
	"fmt"

	"github.com/goark/sshql"
	"github.com/goark/sshql/mysqldrv"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	//	"github.com/PuerkitoBio/goquery"
	//	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi"
	//	ghsrapi "github.com/Chouette2100/srapi"
)

/*

	00AA00	新規作成

*/


type DBConfig struct {
	UseSSH bool   `yaml:"UseSSH"`
	Dbhost string `yaml:"DBhost"`
	Dbport string `yaml:"DBport"`
	Dbname string `yaml:"DBname"`
	Dbuser string `yaml:"DBuser"`
	Dbpw   string `yaml:"DBpw"`
	Sracct string `yaml:"SRacct"`
	Srpswd string `yaml:"SRpswd"`
}

type SSHConfig struct {
	Hostname   string `yaml:"Hostname"`
	Port       int    `yaml:"Port"`
	Username   string `yaml:"Username"`
	Password   string `yaml:"Password"`
	PrivateKey string `yaml:"PrivateKey"`
}

var sshconfig *SSHConfig

var Dialer sshql.Dialer
var Db *sql.DB  //  プログラム中では一貫してこの変数を使うこと
var Dberr error //  プログラム中では一貫してこの変数を使うこと

func OpenDb(dbconfig *DBConfig) (err error) {

	//	https://leben.mobi/go/mysql-connect/practice/
	//	OS := runtime.GOOS

	//	https://ssabcire.hatenablog.com/entry/2019/02/13/000722
	//	https://konboi.hatenablog.com/entry/2016/04/12/100903

	if dbconfig.Dbhost == "" {
		dbconfig.Dbhost = "localhost"
	}
	if dbconfig.Dbport == "" {
		dbconfig.Dbport = "3306"
	}
	cnc := "@tcp"
	if dbconfig.UseSSH {
		err = exsrapi.LoadConfig("SSHConfig.yml", &sshconfig)
		if err != nil {
			err = fmt.Errorf("exsrapi.Loadconfig: %w", err)
			return err
		}
		//	log.Printf("%+v\n", *sshconfig)

		Dialer.Hostname = sshconfig.Hostname
		Dialer.Port = sshconfig.Port
		Dialer.Username = sshconfig.Username
		Dialer.Password = sshconfig.Password
		Dialer.PrivateKey = sshconfig.PrivateKey

		mysqldrv.New(&Dialer).RegisterDial("ssh+tcp")
		cnc = "@ssh+tcp"
	}
	cnc += "(" + dbconfig.Dbhost + ":" + dbconfig.Dbport + ")"
	Db, Dberr = sql.Open("mysql", dbconfig.Dbuser+":"+dbconfig.Dbpw+cnc+"/"+dbconfig.Dbname+"?parseTime=true&loc=Asia%2FTokyo")
	if Dberr != nil {
		err = fmt.Errorf("sql.Open(): %w", Dberr)
	}

	return err
}
