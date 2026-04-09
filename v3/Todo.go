// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import (
	"time"
)

// Todo はTodoテーブルの構造を表します
type Todo struct {
	ID       int        `db:"id, primarykey, autoincrement"`
	Ts       time.Time  `db:"ts"`
	Itype    string     `db:"itype"`
	Target   string     `db:"target"`
	Issue    string     `db:"issue"`
	Solution string     `db:"solution"`
	Closed   *time.Time `db:"closed"` // nullable
}
