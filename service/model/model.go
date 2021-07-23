package model

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"sync"
)

const (
	RouterStatusOnline  = "1"
	RouterStatusOffline = "0"
)

type Router struct {
	Mac        string `json:"mac"`
	FeedId     string `json:"feedId"`
	DeviceName string `json:"name"`
	Status     string `json:"status"`
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
	Mac          string `json:"mac"`
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

var RouterMap = new(SyncRouterMap)
var PointsDetailMap = new(SyncPointsDetailMap)
var TotalPointsMap = new(SyncTotalPointsMap)

type SyncRouterMap struct {
	Val map[string][]*Router
	sync.RWMutex
}

func (th *SyncRouterMap) Set(key string, val []*Router) {
	th.Lock()
	defer th.Unlock()
	if th.Val == nil {
		th.Val = map[string][]*Router{}
	}
	th.Val[key] = val
}

func (th *SyncRouterMap) Read(pin string) []*Router {
	th.RLock()
	defer th.RUnlock()
	return th.Val[pin]
}

// Len 某个key值对应数组的长度
func (th *SyncRouterMap) Len(pin string) int {
	th.RLock()
	defer th.RUnlock()
	return len(th.Val[pin])
}

func (th *SyncRouterMap) MacConvertName(pin, mac string) string {
	th.RLock()
	defer th.RUnlock()
	for _, v := range th.Val[pin] {
		if v.Mac == mac {
			return v.DeviceName
		}
	}
	return "Unknown Device"
}

func (th *SyncRouterMap) MacConvertFeedId(pin, mac string) string {
	th.RLock()
	defer th.RUnlock()
	for _, v := range th.Val[pin] {
		if v.Mac == mac {
			return v.FeedId
		}
	}
	panic(exception.New("mac convert feedId error"))
}

type SyncPointsDetailMap struct {
	Val map[string][]*PointsDetail
	sync.RWMutex
}

func (th *SyncPointsDetailMap) Set(key string, val []*PointsDetail) {
	th.Lock()
	defer th.Unlock()
	if th.Val == nil {
		th.Val = map[string][]*PointsDetail{}
	}
	th.Val[key] = val
}

func (th *SyncPointsDetailMap) Read(pin string) []*PointsDetail {
	th.RLock()
	defer th.RUnlock()
	return th.Val[pin]
}

func (th *SyncPointsDetailMap) Len(pin string) int {
	th.RLock()
	defer th.RUnlock()
	return len(th.Val[pin])
}

func (th *SyncPointsDetailMap) Append(pin string, detail *PointsDetail) {
	th.Lock()
	defer th.Unlock()
	if th.Val == nil {
		th.Val = map[string][]*PointsDetail{}
	}
	th.Val[pin] = append(th.Val[pin], detail)
}

func (th *SyncPointsDetailMap) Clear(pin string) {
	th.Lock()
	defer th.Unlock()
	if th.Val[pin] != nil {
		th.Val[pin] = nil
	}
}

type SyncTotalPointsMap struct {
	Val map[string]*TotalPointsDetail
	sync.RWMutex
}

func (th *SyncTotalPointsMap) Set(key string, val *TotalPointsDetail) {
	th.Lock()
	defer th.Unlock()
	if th.Val == nil {
		th.Val = map[string]*TotalPointsDetail{}
	}
	th.Val[key] = val
}

func (th *SyncTotalPointsMap) Read(pin string) *TotalPointsDetail {
	th.RLock()
	defer th.RUnlock()
	return th.Val[pin]
}
