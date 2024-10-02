// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	//	"fmt"
	"time"
	//	"io"
	//	"log"
	//	"os"
	//	"strconv"
	//	"strings"
	//	"net/http"
	//	"github.com/go-gorp/gorp"
	//      "gopkg.in/gorp.v2"
	//	"github.com/dustin/go-humanize"
	//	"github.com/Chouette2100/srapi"
	//	"github.com/Chouette2100/srapi"
)

/*

0.0.1 新規作成

*/

//	const Version = "0.0.1"

type Campaign struct {
	Campaignid   string
	Campaignname string
	Url          string
	Startedat    time.Time
	Endedat      time.Time
}

type GiftRanking struct {
	Campaignid string
	Grid       int
	Grname     string
	Grtype     string
	Norder     int
	Cntrblst   int
	Startedat  time.Time
	Endedat    time.Time
}
