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
	//	"strings"
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

// ギフト貢献ランキング
type GiftScoreCntrb struct {
	Giftid   int
	Userno   int
	Viewerid int
	Orderno  int
	Score    int
	Ts       time.Time
}

// リスナーギフトランキングデータを格納する
// 必要に応じてリスナー（viewer）の情報を新たに作る、あるいは更新する。
func InserIntoViewerGiftScore(
	client *http.Client,
	dbmap *gorp.DbMap,
	giftid int,
	cugr *srapi.UgrRanking,
	tnow time.Time,
) (
	err error,
) {

	vw := new(Viewer)
	vw.Viewerid = cugr.UserID
	vw.Name = cugr.User.Name

	err = UpinsViewerSetProperty(client, tnow, vw, Env.Lmin)
	if err != nil {
		err = fmt.Errorf("UpinsViewerSetProperty error: %v", err)
		return
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

	// user := new(User)
	user := new(User)
	user.Userno = cgr.RoomID
	//	err = UpinsUserSetProperty(client, tnow, user, 1440*7, 100)
	// err = UpinsUserSetProperty(client, tnow, user, Env.Lmin, Env.Waitmsec)
	_, err = UpinsUser(client, tnow, user, Env.Lmin, Env.Waitmsec)
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


// ギフトランキング・リスナー貢献ランキングデータを格納する
// 必要に応じてリスナー（viewer）の情報を新たに作る、あるいは更新する。
func InserIntoGiftScoreCntrb(
	client *http.Client,
	dbmap *gorp.DbMap,
	giftid int,
	userno int,
	grc *srapi.GrcRanking,
	tnow time.Time,
) (
	err error,
) {

	vw := new(Viewer)
	vw.Viewerid = grc.UserID
	vw.Name = grc.User.Name

	err = UpinsViewerSetProperty(client, tnow, vw, Env.Lmin)
	if err != nil {
		err = fmt.Errorf("UpinsViewerSetProperty error: %v", err)
		return
	}


	//	ユーザーギフトランキングを格納する
	vgs := &GiftScoreCntrb{
		Userno:   userno,
		Giftid:   giftid,
		Orderno:  grc.OrderNo,
		Viewerid: grc.UserID,
		Score:    grc.Score,
		Ts:       tnow,
	}

	err = Dbmap.Insert(vgs)
	if err != nil {
		err = fmt.Errorf("ViewrGiftScore error: %v", err)
		log.Printf("error: %v", err)
		//	return err
	}
	log.Printf("    INSERT(GiftScoreCntrb) viewid=%10d name=%s score=%d\n", vgs.Viewerid, vw.Name, vgs.Score)

	return
}

