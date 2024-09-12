// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"fmt"
	//	"io"
	"log"
	//	"os"
	//	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/go-gorp/gorp"
	//      "gopkg.in/gorp.v2"
	//	"github.com/dustin/go-humanize"

	//	"github.com/Chouette2100/srapi"
	"github.com/Chouette2100/srapi"
)

/*

0.0.1 新規作成

*/

//	const Version = "0.0.1"

// ギフトランキング（イベントギフトランキングではない）
type GiftScore struct {
	Giftid  int
	Userno  int
	Orderno int
	Score   int
	Status  string
	Ts      time.Time
}

// ユーザーギフトランキング（userは配信者の意味で使用しているのでviewerとしてある）
type ViewerGiftScore struct {
	Giftid   int
	Viewerid int
	Orderno  int
	Score    int
	Status   string
	Ts       time.Time
}

// ユーザーギフトランキングデータを格納する
// 必要に応じて視聴者（viewer）の情報を新たに作る、あるいは更新する。
func InserIntoViewerGiftScore(
	client *http.Client,
	dbmap *gorp.DbMap,
	giftid int,
	cugr *srapi.UgrRanking,
	tnow time.Time,
) (
	err error,
) {

	intfc, err := Dbmap.Get(Viewer{}, cugr.UserID)
	if err != nil {
		err = fmt.Errorf("Dbmap.Get(Viewer,UserID) error: %v", err)
		log.Printf("Dbmap.Get error: %v", err)
		return err
	}

	vwh := &ViewerHistory{}
	vw := &Viewer{}
	if intfc == nil {
		//	viewerにviewerid　のデータが見つからない場合は新たに作成する
		vw = &Viewer{
			Viewerid: cugr.UserID,
			Name:     cugr.User.Name,
			Sname:    cugr.User.Name,
			Ts:       tnow,
		}
		err = Dbmap.Insert(vw)
		if err != nil {
			err = fmt.Errorf("Dbmap.Insert error: %v", err)
			log.Printf("error: %v", err)
			return err
		}
		log.Printf(" ** INSERT(viewer) viewid=%10d name=%s\n", vw.Viewerid, vw.Name)

		//	viewerhistoryにデータを作る
		vwh = &ViewerHistory{
			Viewerid: cugr.UserID,
			Name:     cugr.User.Name,
			Sname:    cugr.User.Name,
			Ts:       tnow,
		}
		err = Dbmap.Insert(vwh)
		if err != nil {
			err = fmt.Errorf("Dbmap.Insert error: %v", err)
			log.Printf("error: %v", err)
			return err
		}
		log.Printf(" ** INSERT(viewerhistory) viewid=%10d name=%s\n", vwh.Viewerid, vwh.Name)

	} else {
		//	viewerにvieweridのデータがすでに存在する
		vw = intfc.(*Viewer)

		nodata := false
		vh := ViewerHistory{}
		sqlst := "select max(ts) ts from viewerhistory where viewerid = ? "
		err = dbmap.SelectOne(&vh, sqlst, vw.Viewerid)
		if err != nil {
			//	log.Printf("<%s>\n", err.Error())
			if !strings.Contains(err.Error(), "sql: Scan error on column index 0, name \"ts\": unsupported Scan") {
				err = fmt.Errorf("Dbmap.SelectOne error: %v", err)
				log.Printf("error: %v", err)
				return err
			}
			nodata = true
		}

		pintf := interface{}(nil)
		if !nodata {
			pintf, err = dbmap.Get(ViewerHistory{}, vw.Viewerid, vh.Ts)
			if err != nil {
				err = fmt.Errorf("Dbmap.Get error: %v", err)
				log.Printf("error: %v", err)
				return err
			}
			vwh = pintf.(*ViewerHistory)
		}

		//	if tnow.Sub(pvh.Ts) > 7*24*time.Hour && viewer.Name != cugr.User.Name {
		if nodata {
			//	viewhistoryにデータが存在しない
			vwh = &ViewerHistory{
				Viewerid: cugr.UserID,
				Name:     cugr.User.Name,
				Sname:    cugr.User.Name,
				Ts:       tnow,
			}
			err = Dbmap.Insert(vwh)
			if err != nil {
				err = fmt.Errorf("Dbmap.Insert error: %v", err)
				log.Printf("error: %v", err)
				return err
			}
			log.Printf(" ** INSERT(viewerhistory) viewid=%10d name=%s\n", vwh.Viewerid, vwh.Name)
		} else if tnow.Sub(vwh.Ts) > time.Duration(Env.Lmin) * time.Minute && vh.Name != cugr.User.Name {
			vw.Name = cugr.User.Name
			vw.Ts = tnow
			_, err = Dbmap.Update(vw)
			if err != nil {
				err = fmt.Errorf("Dbmap.Update error: %v", err)
				log.Printf("error: %v", err)
				return err
			}
			log.Printf(" ** UPDATE(viewer) viewid=%10d name=%s from %s\n", vw.Viewerid, vw.Name, vwh.Name)

			//	vh := intfc.(*ViewerHistory)
			err = Dbmap.Insert(vwh)
			if err != nil {
				err = fmt.Errorf("Dbmap.Insert(vh) error: %v", err)
				log.Printf("error: %v", err)
				return err
			}
			log.Printf("    INSERT(viewerhistory) viewid=%10d name=%s from k\n", vh.Viewerid, vh.Name)

		} else {
			log.Printf(" ** SKIP(viewer/viewerhistory)   viewid=%10d name=%s\n", vw.Viewerid, vw.Name)
		}
	}

	//	ユーザーギフトランキングを格納する
	vgs := &ViewerGiftScore{
		Giftid:   giftid,
		Orderno:  cugr.OrderNo,
		Viewerid: cugr.UserID,
		Score:    cugr.Score,
		Status:   "",
		Ts:       tnow,
	}

	err = Dbmap.Insert(vgs)
	if err != nil {
		err = fmt.Errorf("ViewrGiftScore error: %v", err)
		log.Printf("error: %v", err)
		//	return err
	}
	log.Printf("    INSERT(viewerGiftScore) viewid=%10d name=%s score=%d\n", vgs.Viewerid, vw.Name, vgs.Score)

	return
}

// ギフトランキングデータを格納する
// 必要に応じて配信者（user）の情報を新たに作る、あるいは更新する。
func InserIntoGiftScore(
	client *http.Client,
	dbmap *gorp.DbMap,
	giftid int,
	cgr *srapi.GrRanking,
	tnow time.Time,
) (
	err error,
) {

	user := new(User)
	user.Userno = cgr.RoomID
	//	err = UpinsUserSetProperty(client, tnow, user, 1440*7, 100)
	err = UpinsUserSetProperty(client, tnow, user, Env.Lmin, Env.Waitmsec )
	if err != nil {
		err = fmt.Errorf("UpinsUserSetProperty error: %v", err)
		log.Printf("UpinsUserSetProperty error: %v", err)
		//	return
	}

	giftScore := &GiftScore{
		Giftid:  giftid,
		Userno:  cgr.RoomID,
		Orderno: cgr.OrderNo,
		Score:   cgr.Score,
		Status:  "",
		Ts:      tnow,
	}
	err = dbmap.Insert(giftScore)
	if err != nil {
		err = fmt.Errorf("dbmap.Insert error: %v", err)
		return
	}
	log.Printf("    INSERT(GiftScore) userno=%10d\n", cgr.RoomID)
	return
}
