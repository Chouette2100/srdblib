// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"github.com/go-gorp/gorp"
	// "github.com/Chouette2100/srdblib/v3"
)

func AddTableWithName(dbmap *gorp.DbMap) {
	dbmap.AddTableWithName(User{}, "user").SetKeys(false, "Userno")
	dbmap.AddTableWithName(Wuser{}, "wuser").SetKeys(false, "Userno")
	dbmap.AddTableWithName(Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")
	dbmap.AddTableWithName(Event{}, "event").SetKeys(false, "Eventid")
	dbmap.AddTableWithName(Wevent{}, "wevent").SetKeys(false, "Eventid")
	dbmap.AddTableWithName(Eventuser{}, "eventuser").SetKeys(false, "Eventid", "Userno")
	dbmap.AddTableWithName(Weventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")

	dbmap.AddTableWithName(GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
	dbmap.AddTableWithName(ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	dbmap.AddTableWithName(Viewer{}, "viewer").SetKeys(false, "Viewerid")
	dbmap.AddTableWithName(ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")
	// srdblib.Dbmap.AddTableWithName(ShowroomCGIlib.Contribution{}, "contribution").SetKeys(false, "Ieventid", "Roomid", "Viewerid")

	dbmap.AddTableWithName(Campaign{}, "campaign").SetKeys(false, "Campaignid")
	dbmap.AddTableWithName(GiftRanking{}, "giftranking").SetKeys(false, "Campaignid", "Grid")
	dbmap.AddTableWithName(Accesslog{}, "accesslog").SetKeys(false, "Ts", "Eventid")

	dbmap.AddTableWithName(Todo{}, "todo").SetKeys(false, "ID")
}
