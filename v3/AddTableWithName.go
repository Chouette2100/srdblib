// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	// "github.com/Chouette2100/srdblib/v3"
)

func AddTableWithName() {
	Dbmap.AddTableWithName(User{}, "user").SetKeys(false, "Userno")
	Dbmap.AddTableWithName(Wuser{}, "wuser").SetKeys(false, "Userno")
	Dbmap.AddTableWithName(Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")
	Dbmap.AddTableWithName(Event{}, "event").SetKeys(false, "Eventid")
	Dbmap.AddTableWithName(Wevent{}, "wevent").SetKeys(false, "Eventid")
	Dbmap.AddTableWithName(Eventuser{}, "eventuser").SetKeys(false, "Eventid", "Userno")
	Dbmap.AddTableWithName(Weventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")

	Dbmap.AddTableWithName(GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
	Dbmap.AddTableWithName(ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	Dbmap.AddTableWithName(Viewer{}, "viewer").SetKeys(false, "Viewerid")
	Dbmap.AddTableWithName(ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")
	// srdblib.Dbmap.AddTableWithName(ShowroomCGIlib.Contribution{}, "contribution").SetKeys(false, "Ieventid", "Roomid", "Viewerid")

	Dbmap.AddTableWithName(Campaign{}, "campaign").SetKeys(false, "Campaignid")
	Dbmap.AddTableWithName(GiftRanking{}, "giftranking").SetKeys(false, "Campaignid", "Grid")
	Dbmap.AddTableWithName(Accesslog{}, "accesslog").SetKeys(false, "Ts", "Eventid")

	Dbmap.AddTableWithName(Todo{}, "todo").SetKeys(false, "ID")
}
