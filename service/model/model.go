package model

var RouterMap = make(map[string][]*Router)
var PointsDetailMap = make(map[string][]*PointsDetail)
var TotalPointsMap = make(map[string]*TotalPointsDetail)

type Router struct {
	Mac        string `json:"mac"`
	FeedId     string `json:"feedId"`
	DeviceName string `json:"name"`
}

type RouterStatus struct {
	Rom        string `json:"rom"`
	Upload     string `json:"upload"`
	Download   string `json:"download"`
	Cpu        string `json:"cpu"`
	Mem        string `json:"mem"`
	Ip         string `json:"ip"`
	OnlineTime string `json:"onlineTime"`
}

type PointsDetail struct {
	Name         string `json:"name"`
	TodayIncome  int64  `json:"todayIncome" description:"单台今收入"`
	AllIncome    int64  `json:"allIncome" description:"单台总收入"`
	RemainIncome int64  `json:"remainIncome" description:"单台总剩余收入"`
	WaitFreeDay  int64  `json:"waitFreeDay" description:"坐享其成打卡天数"`
}

type TotalPointsDetail struct {
	TotalToday  int64 `json:"totalToday" description:"今日总收入"`
	TotalIncome int64 `json:"totalIncome" description:"总收入"`
	TotalRemain int64 `json:"totalRemain" description:"总剩余"`
}

func MacConvertName(pin, mac string) string {
	for _, v := range RouterMap[pin] {
		if v.Mac == mac {
			return v.DeviceName
		}
	}
	return "Unknown Device"
}
