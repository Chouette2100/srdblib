// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"fmt"
	"log"
	"time"

	// "net/http"

	"github.com/jinzhu/copier"
	// "github.com/mitchellh/copystructure"
	// "github.com/Chouette2100/srapi"
)

type EventuserR interface {
	Get() (*Eventuser, error)
	Set(*Eventuser) error
	/*
	SelectEudata() (*Eventuser, error)
	InsertEutable() error
	UpdateEutable() error
	UpinsEventuser(time.Time) error
	*/
}

// イベントに参加しているユーザの構造体(ベース)
type EventuserBR struct {
	Eventid       string // PRIMARY KEY(1)
	Userno        int    // PRIMARY KEY(2)
	Istarget      string
	Iscntrbpoints string
	Graph         string
	Color         string
	Point         int
	Vld           int
}

type Eventuser struct {
	EventuserBR
	// 他のフィールド
	Status int //	1: ユーザーによって指定された＝無条件にデータ取得対象とする
}

func (e *Eventuser) Get() (
	result *Eventuser,
	err error,
) {
	result = e
	return
}

func (e *Eventuser) Set(ne *Eventuser) (err error) {
	*e = *ne
	return nil
}

type Weventuser struct {
	EventuserBR
	// 他のフィールド
	Status int //	1: ユーザーによって指定された＝無条件にデータ取得対象とする
}

func (e *Weventuser) Get() (
	result *Eventuser,
	err error,
) {
	result = new(Eventuser)
	err = copier.Copy(result, e)
	if err != nil {
		err = fmt.Errorf("copier.Copy failed: %w", err)
		return
	}
	return
}

func (e *Weventuser) Set(ne *Eventuser) (err error) {
	err = copier.Copy(e, ne)
	if err != nil {
		err = fmt.Errorf("copier.Copy failed: %w", err)
	}
	return
}

func SelectEudata[T EventuserR](xeu T, eventid string, userno int) (
	result *T,
	err error,
) {

	var intf interface{}

	intf, err = Dbmap.Get(xeu, eventid, userno)
	if err != nil {
		err = fmt.Errorf("Dbmap.Get failed: %w", err)
		return
	} else if intf == nil {
		result = nil
		return
	} else {
		fmt.Printf("intf type: %T\n", intf)
		p := intf.(T)
		fmt.Printf("intf type: %T\n", intf)
		result = &p
	}
	return
}

func UpdateEutable[T EventuserR](xeu T) (err error) {

	var nr int64
	nr, err = Dbmap.Update(xeu)
	if err != nil {
		err = fmt.Errorf("Dbmap.Update failed: %w", err)
	} else if nr == 0 {
		err = fmt.Errorf("Dbmap.Update failed: nr = %d", nr)
	}
	return
}

func InsertEutable[T EventuserR](xeu T) (err error) {
	err = Dbmap.Insert(xeu)
	return
}

// イベント最終結果（確定結果）をeventuserテーブルに格納する（既存の場合は更新する）
func UpinsEventuserG[T EventuserR](xeu T, tnow time.Time) (err error) {

	// イベントユーザー情報を取得する
	var teu *Eventuser
	var cxeu *T
	teu, err = xeu.Get()
	if err != nil {
		err = fmt.Errorf("xeu.Get failed: %w", err)
		return
	}
	eventid := teu.Eventid
	userno := teu.Userno
	cxeu, err = SelectEudata(xeu, eventid, userno)
	if err != nil {
		err = fmt.Errorf("Dbmap.Get failed: %w", err)
		return
	} else if cxeu == nil {
		// なければ新規作成
		teu.Istarget = "Y"
		teu.Iscntrbpoints = "N"
		if teu.Vld < 21 {
			teu.Graph = "Y"
		} else {
			teu.Graph = "N"
		}
		teu.Color = Colorlist2[(teu.Vld-1)%len(Colorlist2)].Name
		// TODO: Vldがマイナスのときの処理を追加（あるいはレベルイベントのときの処理を追加）

		xeu.Set(teu)
		err = InsertEutable(xeu)
		if err != nil {
			err = fmt.Errorf("InsertEutable failed: %w", err)
			return
		} else {
			log.Printf("InsertEutable() success: %s %d\n", teu.Eventid, teu.Userno)
		}
	} else {
		// あれば更新
		var ceu *Eventuser
		ceu, err = (*cxeu).Get()
		if err != nil {
			err = fmt.Errorf("cxeu.Get failed: %w", err)
			return
		}

		if teu.Vld > 0 {
			ceu.Vld = teu.Vld
		}
		if teu.Point > ceu.Point {
			ceu.Point = teu.Point
		}
		ceu.Color = Colorlist2[(ceu.Vld-1)%len(Colorlist2)].Name
		// TODO: Vldがマイナスのときの処理を追加（あるいはレベルイベントのときの処理を追加）

		xeu.Set(ceu)
		err = UpdateEutable(xeu)
		if err != nil {
			err = fmt.Errorf("UpdateEutable failed: %w", err)
			return
		} else {
			log.Printf("UpdateEutable() success: %s %d\n", teu.Eventid, teu.Userno)
		}
	}
	return
}
