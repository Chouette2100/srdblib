//	Copyright © 2024 chouette.21.00@gmail.com
//	Released under the MIT license
//	https://opensource.org/licenses/mit-license.php
package srdblib

import (
	//	"fmt"
	//	"io"
	//	"log"
	//	"os"
	//	"strconv"
	//	"strings"
	"time"

	//	"net/http"

	//	"github.com/go-gorp/gorp"
	//      "gopkg.in/gorp.v2"

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
  Ts time.Time
}

type ViewerHistory struct{
  Viewerid int
  Name string
  Ts time.Time
}