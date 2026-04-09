// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"time"
)

type Eventrank struct {
	Eventid   string
	Userid    int
	Ts        time.Time
	Listner   string
	Lastname  string
	Lsnid     int
	T_lsnid   int
	Norder    int
	Nrank     int
	Point     int
	Increment int
	Status    int
}
