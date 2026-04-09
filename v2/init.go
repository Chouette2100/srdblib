package srdblib

var Clmlist map[string]string = map[string]string{}

func init() {
	/*
		Clmlist["user"] = ExtractStructColumns(&User{})
		Clmlist["event"] = ExtractStructColumns(&Event{})
		Clmlist["eventuser"] = ExtractStructColumns(&Eventuser{})

		Clmlist["wevent"] = ExtractStructColumns(&Wevent{})
		Clmlist["weventuser"] = ExtractStructColumns(&Weventuser{})
	*/
	Clmlist = map[string]string{
		"user": "`userno`, `userid`, `user_name`, `longname`, `shortname`, `genre`, `genreid`, `rank`, `nrank`, `prank`, `irank`, `inrank`, `iprank`, `itrank`, `level`, `followers`, `fans`, `fanpower`, `fans_lst`, `fanpower_lst`, `ts`, `getp`, `graph`, `color`, `currentevent`",

		"event": "`eventid`, `ieventid`, `event_name`, `period`, `starttime`, `endtime`, `noentry`, `intervalmin`, `modmin`, `modsec`, `fromorder`, `toorder`, `resethh`, `resetmm`, `nobasis`, `maxdsp`, `cmap`, `target`, `rstatus`, `maxpoint`, `thinit`, `thdelta`, `achk`, `aclr`",

		"eventuser": "`eventid`, `userno`, `istarget`, `iscntrbpoints`, `graph`, `color`, `point`, `vld`, `status`",

		"wevent": "`eventid`, `ieventid`, `event_name`, `period`, `starttime`, `endtime`, `noentry`, `intervalmin`, `modmin`, `modsec`, `fromorder`, `toorder`, `resethh`, `resetmm`, `nobasis`, `maxdsp`, `cmap`, `target`, `rstatus`, `maxpoint`, `thinit`, `thdelta`, `achk`, `aclr`",

		"weventuser": "`eventid`, `userno`, `istarget`, `iscntrbpoints`, `graph`, `color`, `point`, `vld`, `status`",
	}

}
