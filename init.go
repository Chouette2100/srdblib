package srdblib

var Clmlist map[string]string = map[string]string{}

func init() {
	Clmlist["user"] = ExtractStructColumns(&User{})
	Clmlist["event"] = ExtractStructColumns(&Event{})
	Clmlist["eventuser"] = ExtractStructColumns(&Eventuser{})

	Clmlist["wevent"] = ExtractStructColumns(&Wevent{})
	Clmlist["weventuser"] = ExtractStructColumns(&Weventuser{})
}
