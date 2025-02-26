package srdblib
import (
	"fmt"
	"time"
	"github.com/jinzhu/copier"
)
// userの履歴を保存する構造体
// PRIMARY KEY (`userno`,`ts`)
type Userhistory struct {
	Userno    int
	User_name string
	Genre     string
	Rank      string
	Nrank     string
	Prank     string
	Level     int
	Followers int
	Fans      int
	Fans_lst  int
	Ts        time.Time
}

type Wuserhistory Userhistory

// UserhistoryT is an interface for Userhistory
type UserhistoryT interface {
	Get() (*User, error)
	Set(*User) error
}	

// Getter and Setter for Userhistory
func (uh *Userhistory) Get() (
	result *User,
	err error,
) {
	copier.Copy(result, uh)
	return
}

func (uh *Userhistory) Set(nuh *User) (err error) {
	copier.Copy(uh, nuh)
	return
}

// Getter and Setter for Wuserhistory
func (wh *Wuserhistory) Get() (
	result *User,
	err error,
) {
	copier.Copy(result, wh)
	return
}

func (wh *Wuserhistory) Set(nwh *User) (err error) {
	copier.Copy(wh, nwh)
	return
}

// userデータをuserhistory, wuserhistoryにinsertする。
func InsertUserhistory[T UserhistoryT](
	xuserhistory T, user *User,
) (
	err error,
) {
	xuserhistory.Set(user)
	err = Dbmap.Insert(xuserhistory)
	if err != nil {
		err = fmt.Errorf("Dbmap.Insert failed: %w", err)
	}
	return
}

