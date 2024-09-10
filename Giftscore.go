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
		err = fmt.Errorf("Dbmap.Get error: %v", err)
		log.Printf("Dbmap.Get error: %v", err)
		return err
	}

	if intfc == nil {
		//	viewerid　が見つからない場合は新たに作成する
		viewer := &Viewer{
			Viewerid: cugr.UserID,
			Name:     cugr.User.Name,
			Ts:       tnow,
		}
		err = Dbmap.Insert(viewer)
		if err != nil {
			err = fmt.Errorf("Dbmap.Insert error: %v", err)
			log.Printf("error: %v", err)
			return err
		}
	}

	viewerGiftScore := &ViewerGiftScore{
		Giftid:   giftid,
		Orderno:  cugr.OrderNo,
		Viewerid: cugr.UserID,
		Score:    cugr.Score,
		Status:   "",
		Ts:       tnow,
	}

	err = Dbmap.Insert(viewerGiftScore)
	if err != nil {
		err = fmt.Errorf("ViewrGiftScore error: %v", err)
		log.Printf("error: %v", err)
		//	return err
	}

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
	err = UpinsUserSetProperty(client, tnow, user, 1440*7, 100)
	if err != nil {
		err = fmt.Errorf("UpinsUserSetProperty error: %v", err)
		log.Printf("UpinsUserSetProperty error: %v", err)
		//	return
	}

	giftScore := &GiftScore{
		Giftid:  giftid,
		Userno:  cgr.RoomID,
		Orderno:  cgr.OrderNo,
		Score:  cgr.Score,
		Status:  "",
		Ts:      tnow,
	}
	err = dbmap.Insert(giftScore)
	if err != nil {
		err = fmt.Errorf("dbmap.Insert error: %v", err)
		return
	}
	return
}
