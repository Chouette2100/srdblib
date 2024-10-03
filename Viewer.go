//	Copyright © 2024 chouette.21.00@gmail.com
//	Released under the MIT license
//	https://opensource.org/licenses/mit-license.php
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

	//	"github.com/go-gorp/gorp"

	//	"github.com/dustin/go-humanize"

	//	"github.com/Chouette2100/exsrapi"
	//	"github.com/Chouette2100/srapi"
)

/*

0.0.1 新規作成

*/

//	const Version = "0.0.1"

type Viewer struct {
  Viewerid int
  Name string
  Sname string
  Ts time.Time
  Orderno int	//	GiftScoreのOrdernoを受けるために追加したメンバー、テーブルには存在しない
}

type ViewerHistory struct{
  Viewerid int
  Name string	//	ランキングデータにあるリスナー名
  Sname string	//	リスナーの表示名（User.Longnameと同様、User.Shortnameに相当するものはない）
  Ts time.Time
}

/*
テーブルviewerにリスナーの新規登録（あるいは更新登録）を行う
*/
func UpinsViewerSetProperty(client *http.Client, tnow time.Time, viewer *Viewer, lmin int) (
	err error,
) {

	if viewer.Sname == "" {
		viewer.Sname = viewer.Name
	}

	/*
		既存のvieweridの場合は更新、新規の場合は新規作成
	*/

	intfc, err := Dbmap.Get(Viewer{}, viewer.Viewerid)
	if err != nil {
		err = fmt.Errorf("Dbmap.Get(Viewer,UserID) error: %v", err)
		log.Printf("Dbmap.Get error: %v", err)
		return err
	}

	vwh := &ViewerHistory{}
	vw := &Viewer{}
	if intfc == nil {
		//	viewerにviewerid　のデータが見つからない場合は新たに作成する
		vw = viewer
		err = Dbmap.Insert(vw)
		if err != nil {
			err = fmt.Errorf("Dbmap.Insert error: %v", err)
			log.Printf("error: %v", err)
			return err
		}
		log.Printf(" ** INSERT(viewer) viewid=%10d name=%s\n", vw.Viewerid, vw.Name)

		//	viewerhistoryにデータを作る
		vwh = &ViewerHistory{
			Viewerid: viewer.Viewerid,
			Name:     viewer.Name,
			Sname:    viewer.Name,
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
		err = Dbmap.SelectOne(&vh, sqlst, vw.Viewerid)
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
			pintf, err = Dbmap.Get(ViewerHistory{}, vw.Viewerid, vh.Ts)
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
				Viewerid: viewer.Viewerid,
				Name:     viewer.Name,
				Sname:    viewer.Name,
				Ts:       tnow,
			}
			err = Dbmap.Insert(vwh)
			if err != nil {
				err = fmt.Errorf("Dbmap.Insert error: %v", err)
				log.Printf("error: %v", err)
				return err
			}
			log.Printf(" ** INSERT(viewerhistory) viewid=%10d name=%s\n", vwh.Viewerid, vwh.Name)
		} else if tnow.Sub(vwh.Ts) > time.Duration(lmin)*time.Minute && vh.Name != viewer.Name {
			vw.Name = viewer.Name
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
	return
}