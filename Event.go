// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"time"
	"fmt"
	"github.com/jinzhu/copier"
)

//	Gorpのための構造体定義

//	0.0.0 新規に作成する

// イベント構造体
// PRIMARY KEY (eventid)
type Event struct {
	Eventid     string `yaml:"eventid"` // イベントID　（event_url_key）
	Ieventid    int `yamal:"ieventid"`    //	イベントID　（整数）
	Event_name  string `yaml:"event_name"`
	Period      string `yaml:"period"`
	Starttime   time.Time `yaml:"starttime"` // イベント開始時刻
	Endtime     time.Time `yaml:"endtime"`   // イベント終了時刻
	Noentry     int `yaml:"noentry"`
	Intervalmin int `yaml:"intervalmin"`
	Modmin      int `yaml:"modmin"`
	Modsec      int `yaml:"modsec"`
	Fromorder   int `yaml:"fromorder"`
	Toorder     int `yaml:"toorder"`
	Resethh     int `yaml:"resethh"`
	Resetmm     int `yaml:"resetmm"`
	Nobasis     int `yaml:"nobasis"`
	Maxdsp      int `yaml:"maxdsp"`
	Cmap        int `yaml:"cmap"`
	Target      int `yaml:"target"`
	Rstatus     string `yaml:"rstatus"`
	Maxpoint    int `yaml:"maxpoint"`
	Thinit      int `yaml:"thinit"` //	獲得ポイントがThinit + Thdelta * int(time.Since(Starttime).Hours())を超えるルームのみデータ取得対象とする。
	Thdelta     int `yaml:"thdelta"`
	Achk        int `yaml:"achk"`
	Aclr        int `yaml:"aclr"`
}

type Wevent struct {
	Event
}

type EventC struct  {
	Event
}

func (e *Event) Get() (result *EventC, err error) {
	result = new(EventC)
	err = copier.Copy(result, e)
	if err != nil {
		err = fmt.Errorf("copier.Copy() error: %w", err)
		return
	}
	return
}
func (e *Wevent) Get() (result *EventC, err error) {
	result = new(EventC)
	err = copier.Copy(result, e)
	if err != nil {
		err = fmt.Errorf("copier.Copy() error: %w", err)
		return
	}
	return
}

func (e *Event) Set(event *EventC) (err error) {
	err = copier.Copy(e, event)
	if err != nil {
		err = fmt.Errorf("copier.Copy() error: %w", err)
		return
	}
	return
}

func (e *Wevent) Set(event *EventC) (err error) {
	err = copier.Copy(e, event)
	if err != nil {
		err = fmt.Errorf("copier.Copy() error: %w", err)
		return
	}
	return
}

//	event := Event{
//		Intervalmin: 5,
//		Modmin:      4,
//		Modsec:      10,
//		Fromorder:   1,
//		Toorder:     10,
//		Resethh:     4,
//		Resetmm:     0,
//		Nobasis:     164614,
//		Maxdsp:      10,
//		Cmap:        1,
//	}

type EventR interface {
	//	イベント情報を取得する
	Get() (*EventC, error)
	//	イベント情報を更新する
	Set(*EventC) error
}
