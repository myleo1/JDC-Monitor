package service

import (
	"JDC-Monitor/service/cryptokit"
	"JDC-Monitor/service/model"
	"github.com/go-resty/resty/v2"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"time"
)

func InitRouterList(pin, tgt string) {
	model.RouterMap.Set(pin, ListRouter(tgt))
	if model.RouterMap.Len(pin) == 0 {
		logkit.Fatal("获取路由器列表失败")
	}
}

// GetPCDNStatus 获取路由器状态:ip,cpu,memory,rom,上传下载等信息
func GetPCDNStatus(feedId string, tgt string) *model.RouterStatus {
	accessKey := "b8f9c108c190a39760e1b4e373208af5cd75feb4"
	body := `{"feed_id":"` + feedId + `","command":[{"stream_id":"SetParams","current_value":"{\n  \"cmd\" : \"get_router_status_detail\"\n}"}]}`
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Authorization": cryptokit.EncodeAuthorization(body, accessKey),
			"timestamp":     time.Now().Format(timekit.TimeLayout),
			"accesskey":     accessKey,
			"tgt":           tgt,
			"User-Agent":    "ios",
			"appkey":        "996",
			//"pin":           pin,
		}).
		//SetBody(`{"feed_id":"457381614087937282","command":[{"stream_id":"SetParams","current_value":"{\n  \"cmd\" : \"jdcplugin_opt.get_pcdn_status\"\n}"}]}`).
		SetBody(body).
		Post("https://gw.smart.jd.com/f/service/controlDevice?plat=ios&hard_platform=iPhone11%2C8&app_version=6.5.5&plat_version=13.7&channel=jd")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "status").String() != "0" {
		panic(exception.New("获取路由器状态失败"))
	}
	m := jsonkit.ParseMap(gjson.Get(resp.String(), "result").String())
	data := cast.ToStringMap(cast.ToStringMap(cast.ToStringMap(cast.ToSlice(m["streams"])[0])["current_value"])["data"])
	s := &model.RouterStatus{}
	s.Ip = cast.ToString(data["wanip"])
	s.Cpu = cast.ToString(data["cpu"])
	s.Mem = cast.ToString(data["mem"])
	s.Rom = cast.ToString(data["rom"])
	s.OnlineTime = cast.ToString(data["onlineTime"])
	s.Upload = cast.ToString(data["upload"])
	s.Download = cast.ToString(data["download"])
	return s
}

// ListRouter 获取路由器列表,得到mac、name等信息
func ListRouter(tgt string) []*model.Router {
	accessKey := "b8f9c108c190a39760e1b4e373208af5cd75feb4"
	body := `{"appversion":"2.7.3","appplatform":"iPhone11,8","appplatformversion":"13.7"}`
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Authorization": cryptokit.EncodeAuthorization(body, accessKey),
			"timestamp":     time.Now().Format(timekit.TimeLayout),
			"accesskey":     accessKey,
			"tgt":           tgt,
			"User-Agent":    "ios",
			"appkey":        "996",
			//"pin":           pin,
		}).
		SetBody(body).
		Post("https://gw.smart.jd.com/f/service/listAllUserDevices?plat=ios&hard_platform=iPhone11,8&app_version=6.5.5&plat_version=13.7&channel=jd")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "status").String() != "0" {
		panic(exception.New("获取路由器列表失败"))
	}
	list := gjson.Get(resp.String(), "result").Array()[0].Get("list").Array()
	routerList := make([]*model.Router, 0, len(list))
	if len(list) > 0 {
		for _, v := range list {
			r := &model.Router{}
			r.Mac = v.Get("device_id").String()
			r.DeviceName = v.Get("device_name").String()
			r.FeedId = v.Get("feed_id").String()
			r.Status = v.Get("status").String()
			routerList = append(routerList, r)
		}
		return routerList
	}
	return nil
}

// GetTotalPoints 获取剩余总收益(扣除已兑换的) -未使用
func GetTotalPoints(wsKey string) int64 {
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeader("wskey", wsKey).
		Get("https://router-app-api.jdcloud.com/v1/regions/cn-north-1/pinTotalAvailPoint")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "code").String() != "200" {
		panic(exception.New("获取总收益失败"))
	}
	res := gjson.Get(resp.String(), "result").Get("totalAvailPoint").Int()
	return res
}

// GetTodayPoints 获取今日总收益
func GetTodayPoints(wsKey string) int64 {
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeader("wskey", wsKey).
		Get("https://router-app-api.jdcloud.com/v1/regions/cn-north-1/todayPointIncome")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "code").String() != "200" {
		panic(exception.New("获取今日总收益失败"))
	}
	res := gjson.Get(resp.String(), "result").Get("todayTotalPoint").Int()
	return res
}

// GetRouterPoints 获取单个路由器总收益(扣除已兑换的)
func GetRouterPoints(mac, wsKey string) int64 {
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeader("wskey", wsKey).
		SetQueryParam("mac", mac).
		Get("https://router-app-api.jdcloud.com/v1/regions/cn-north-1/routerAccountInfo")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "code").String() != "200" {
		panic(exception.New("获取单个路由总收益失败"))
	}
	res := gjson.Get(resp.String(), "result").Get("accountInfo").Get("amount").Int()
	return res
}

// GetWaitFreeDay 获取单个总收益、坐享其成打卡天数
func GetWaitFreeDay(mac, wsKey string) int64 {
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeader("wskey", wsKey).
		SetQueryParam("mac", mac).
		Get("https://router-app-api.jdcloud.com/v1/regions/cn-north-1/router:activityInfo")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "code").String() != "200" {
		panic(exception.New("获取总收益、坐享其成天数失败"))
	}
	res := gjson.Get(resp.String(), "result").Get("routerUnderwayResult").Get("satisfiedTimes").Int()
	return res
}

// GetPointsDetail 获取单台总收益、今日收益汇总
func GetPointsDetail(pin, wsKey string, waitFree bool) {
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeader("wskey", wsKey).
		Get("https://router-app-api.jdcloud.com/v1/regions/cn-north-1/todayPointDetail")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "code").String() != "200" {
		panic(exception.New("获取收益汇总信息失败"))
	}
	res := gjson.Get(resp.String(), "result").Get("pointInfos").Array()
	model.PointsDetailMap.Clear(pin)
	for _, v := range res {
		m := v.Map()
		r := &model.PointsDetail{}
		r.Name = model.RouterMap.MacConvertName(pin, m["mac"].String())
		r.Mac = m["mac"].String()
		r.TodayIncome = m["todayPointIncome"].Int()
		r.AllIncome = m["allPointIncome"].Int()
		time.Sleep(time.Millisecond * 500)
		r.RemainIncome = GetRouterPoints(m["mac"].String(), wsKey)
		time.Sleep(time.Millisecond * 500)
		if waitFree {
			r.WaitFreeDay = GetWaitFreeDay(m["mac"].String(), wsKey)
		}
		model.PointsDetailMap.Append(pin, r)
		GetPointsDetailTotal(pin, wsKey)
	}
}

// GetPointsDetailTotal 获取总收益、剩余总收益
func GetPointsDetailTotal(pin, wsKey string) {
	m := &model.TotalPointsDetail{}
	m.TotalToday = GetTodayPoints(wsKey)
	var total1, total2 int64
	for _, v := range model.PointsDetailMap.Read(pin) {
		total1 += v.AllIncome
		total2 += v.RemainIncome
	}
	m.TotalIncome = total1
	m.TotalRemain = total2
	model.TotalPointsMap.Set(pin, m)
}

// RebootRouter 重启路由器
func RebootRouter(feedId string, tgt string) {
	accessKey := "b8f9c108c190a39760e1b4e373208af5cd75feb4"
	body := `{"feed_id":"` + feedId + `","command":[{"stream_id":"SetParams","current_value":"{\n  \"cmd\" : \"reboot_system\"\n}"}]}`
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Authorization": cryptokit.EncodeAuthorization(body, accessKey),
			"timestamp":     time.Now().Format(timekit.TimeLayout),
			"accesskey":     accessKey,
			"tgt":           tgt,
			"User-Agent":    "ios",
			"appkey":        "996",
			//"pin":           pin,
		}).
		SetBody(body).
		Post("https://gw.smart.jd.com/f/service/controlDevice?plat=ios&hard_platform=iPhone11%2C8&app_version=6.5.5&plat_version=13.7&channel=jd")
	if err != nil || resp.IsError() || gjson.Get(resp.String(), "status").String() != "0" {
		panic(exception.New("重启路由器失败"))
	}
}
