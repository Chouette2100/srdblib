/*
!
Copyright © 2024 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/
package srdblib

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	//	"github.com/Chouette2100/exsrapi"
	//	"github.com/Chouette2100/srapi"
)
//	イベントに新しいユーザを追加する
//		InserNewOnes()をコピーし引数を変更したもの
//			TODO: こちらに統一すること
//		eventuser	新規作成または更新
//		user 新規作成または更新
//		points イベント開始時のデータを新規作成
func UpinsEventuser(
	client *http.Client,
	order int,
	point int,
	eventid string,
	starttime time.Time,
	userno int,
	tnow time.Time,
) (
	err error,
) {

	//	現在時
	//	tnow := time.Now().Truncate(time.Second)

	//
	//	eventid := eventinf.Event_ID

	//	userno := room.Room_id

	nrow := 0
	sqls := "select count(*) from eventuser where userno =? and eventid = ?"
	Dberr = Db.QueryRow(sqls, userno, eventid).Scan(&nrow)

	if Dberr != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", Dberr.Error())
		err = fmt.Errorf("Db.QueryRow().Scan(&nrow): %w", Dberr)
		return
	}

	Colorlist := Colorlist2
	if Event_inf.Cmap == 1 {
		Colorlist = Colorlist1
	}

	if nrow == 0 {
		//	eventuser に対象ルームが存在しないとき
		//	log.Printf("  =====Insert into eventuser userno=%d, eventid=%s\n", userno, eventid)
		var stmt *sql.Stmt
		sqli := "INSERT INTO eventuser(eventid, userno, istarget, graph, color, iscntrbpoints, point) VALUES(?,?,?,?,?,?,?)"
		stmt, Dberr = Db.Prepare(sqli)
		if Dberr != nil {
			err = fmt.Errorf("Db.Prepare(sqli): %w", Dberr)
			return
		}
		defer stmt.Close()

		//	if i < 10 {
		_, Dberr = stmt.Exec(
			eventid,
			userno,
			"Y",
			"Y",
			Colorlist[order%len(Colorlist)].Name,
			"N",
			point,
		)

		if Dberr != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", Dberr.Error())
			err = fmt.Errorf("stmt.Exec(stmt): %w", Dberr)
			return
		}
		log.Printf("   **** insert into eventuser.\n")

		sqlip := "insert into points (ts, user_id, eventid, point, `rank`, gap, pstatus) values(?,?,?,?,?,?,?)"
		_, Dberr = Db.Exec(
			sqlip,
			starttime.Truncate(time.Second),
			userno,
			eventid,
			0,
			1,
			0,
			"=",
		)
		if Dberr != nil {
			err = fmt.Errorf("Db.Exec(sqlip,...): %w", Dberr)
			return
		}
		log.Printf("   **** insert into points.\n")

		nrowu := 0
		/*
			sqlscu := "select count(*) from user where userno =?"
			Dberr = Db.QueryRow(sqlscu, userno).Scan(&nrowu)
			if Dberr != nil {
				err = fmt.Errorf("Db.QueryRow(sqlscu, userno).Scan(&nrow): %w", Dberr)
				return
			}
		*/
		var row interface{}
		row, err = Dbmap.Get(User{}, userno)
		if err != nil {
			err = fmt.Errorf("Get(userno=%d) returned error. %w", userno, err)
			return err
		}
		if row != nil {
			nrowu = 1
		}

		if nrowu == 0 {
			//	（eventuser に対象ルームが存在せず） userにも対象ルームが存在しないとき

			/*
				shortname := fmt.Sprintf("%d", userno)
				shortname = shortname[len(shortname)-2:]
				sqliu := "insert into user (userno, userid, user_name, longname, shortname, currentevent, ts) values(?,?,?,?,?,?,?)"
				_, Dberr = Db.Exec(
					sqliu,
					userno,
					room.Room_url_key,
					room.Room_name,
					room.Room_name,
					shortname,
					eventid,
					tnow,
				)
				if Dberr != nil {
					err = fmt.Errorf("Db.Exec(sqliu,...): %w", Dberr)
					return
				}
			*/
			//	user テーブルにusernoのデータを新たに作成する
			err = InsertIntoUser(client, tnow, userno)
			if err != nil {
				err = fmt.Errorf("InsertIntoUser(client, tnow, userno): %w", err)
				return
			}
			log.Printf("   **** insert into user.\n")

		} else {
			//	（eventuser に対象ルームが存在しないが） userには対象ルームが存在する
			log.Printf("   **** user already exists.\n")

			/*
			puser := row.(*User)
			err = UpdateUserSetProperty(client, tnow, puser)
			if err != nil {
				err = fmt.Errorf("UpdateUserSetProperty(client, tnow, puser): %w", err)
				return
			}

			log.Printf("   **** update user.\n")
			*/
		}

		/*
			//	ルーム情報を最新にする。
			var roominf srapi.RoomInf
			roominf, err = srapi.ApiRoomProfile(client, fmt.Sprintf("%d", room.Room_id))
			if err != nil {
				err = fmt.Errorf("srapi.ApiRoomProfile(): %w", err)
				return
			}

			roominf.Genre = ConverGenre2Abbr(roominf.Genre)

			var stmtuu *sql.Stmt
			sqluu := "UPDATE user SET "
			sqluu += "  userid=? "
			sqluu += ", user_name=? "
			sqluu += ", genre=? "
			sqluu += ", `rank`=? "
			sqluu += ", nrank=? "
			sqluu += ", prank=? "
			sqluu += ", level=? "
			sqluu += ", followers=? "
			sqluu += ", fans=? "
			sqluu += ", fans_lst=? "
			sqluu += ", currentevent=? "
			sqluu += ", ts=? "
			sqluu += " where userno=?"
			stmtuu, err = Db.Prepare(sqluu)
			if err != nil {
				log.Printf("error(UPDATE/Prepare) err=%s\n", err.Error())
				err = fmt.Errorf("Db.Prepare(sqluu): %w", Dberr)
				return
			}
			defer stmtuu.Close()
			_, Dberr = stmtuu.Exec(
				roominf.Account,
				roominf.Name,
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				eventid,
				tnow,
				userno,
			)
			if Dberr != nil {
				log.Printf("error(InsertIntoOrUpdateUser() UPDATE/Exec) err=%s\n", Dberr.Error())
				err = fmt.Errorf("stmtuu.Exec(): %w", Dberr)
				return
			}
			log.Printf("   **** update user.\n")
		*/

	} else {
		//	eventuser に対象ルームが存在するとき
		//	log.Printf("  =====Update eventuser userno=%d, eventid=%s\n", userno, eventid)
		var stmtu *sql.Stmt
		sqlu := "UPDATE eventuser SET istarget=? where eventid=? and userno=?"
		stmtu, err = Db.Prepare(sqlu)
		if err != nil {
			log.Printf("error(UPDATE/Prepare) err=%s\n", err.Error())
			err = fmt.Errorf("Db.Prepare(sqlu): %w", Dberr)
			return
		}
		defer stmtu.Close()

		_, Dberr = stmtu.Exec(
			"Y",
			eventid,
			userno,
		)

		if Dberr != nil {
			log.Printf("error(Update eventuser) err=%s\n", Dberr.Error())
			err = fmt.Errorf("stmtu.Exec(): %w", Dberr)
			return
		}
		log.Printf("   **** update eventuser.\n")
	}

	return

}
