package task

import (
	"JDC-Monitor/service"
	"JDC-Monitor/service/config"
	"JDC-Monitor/service/model"
	"JDC-Monitor/service/wechat"
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/commonkit"
	"github.com/mizuki1412/go-core-kit/service/cronkit"
	"github.com/mizuki1412/go-core-kit/service/influxkit"
	"github.com/spf13/cast"
	"time"
)

func UpdateRouterList(pin, tgt string) {
	//更新路由器列表
	cronkit.AddFunc("@every 10m", func() {
		model.RouterMap.Set(pin, service.ListRouter(tgt))
	})
}

func CollectTask(pin, tgt string) {
	//采集路由器数据
	cronkit.AddFunc("@every 1m", func() {
		for _, v := range model.RouterMap.Read(pin) {
			if v.Status != model.RouterStatusOffline {
				_ = commonkit.RecoverFuncWrapper(func() {
					s := service.GetPCDNStatus(v.FeedId, tgt)
					sql := fmt.Sprintf("%s ip=%s,online=%s,cpu=%s,mem=%s,upload=%s,download=%s,rom=%s %d", v.Mac, influxkit.Decorate(s.Ip), s.OnlineTime, s.Cpu, s.Mem, s.Upload, s.Download, influxkit.Decorate(s.Rom), time.Now().UnixNano())
					influxkit.WriteDefaultDB(sql)
					//log.Println(sql)
				})
			}
		}
	})
}

func RebootTask(pin, tgt, user string, waitFree bool) {
	//如果收益低于一定阈值，自动重启
	cronkit.AddFunc("0 45 6 */1 * ?", func() {
		//tgt与wsKey相等
		service.GetPointsDetail(pin, tgt, waitFree)
		for _, v := range model.PointsDetailMap.Read(pin) {
			_ = commonkit.RecoverFuncWrapper(func() {
				if threshold := config.Conf[pin].Reboot; v.TodayIncome < threshold {
					feedId := model.RouterMap.MacConvertFeedId(pin, v.Mac)
					service.RebootRouter(feedId, tgt)
					content := fmt.Sprintf("********京东云矿机********\n\n【%s】收益低于%s,已重启", v.Name, cast.ToString(threshold))
					wechat.Push2Wechat(user, content)
				}
			})
		}
	})
}

func PushPointTask(pin, tgt, user string, waitFree bool) {
	//获取路由器积分等信息、然后推送至企业微信
	cronkit.AddFunc("0 30 7 */1 * ?", func() {
		//tgt与wsKey相等
		service.GetPointsDetail(pin, tgt, waitFree)
		var content string
		totalPoints := model.TotalPointsMap.Read(pin)
		content += "********京东云矿机********\n\n"
		content += fmt.Sprintf("今日总收：%s\n总收：%s\n总剩：%s\n\n", cast.ToString(totalPoints.TotalToday), cast.ToString(totalPoints.TotalIncome), cast.ToString(totalPoints.TotalRemain))
		l := model.PointsDetailMap.Len(pin)
		for i, v := range model.PointsDetailMap.Read(pin) {
			if waitFree {
				if l != i+1 {
					content += fmt.Sprintf("设备名:%s\n坐享其成:%s\n单台今收:%s\n单台总收:%s\n单台剩余:%s\n\n", v.Name, cast.ToString(v.WaitFreeDay), cast.ToString(v.TodayIncome), cast.ToString(v.AllIncome), cast.ToString(v.RemainIncome))
					continue
				}
				//最后一台设备最后一个回车
				content += fmt.Sprintf("设备名:%s\n坐享其成:%s\n单台今收：%s\n单台总收:%s\n单台剩余:%s\n", v.Name, cast.ToString(v.WaitFreeDay), cast.ToString(v.TodayIncome), cast.ToString(v.AllIncome), cast.ToString(v.RemainIncome))
			} else {
				if l != i+1 {
					content += fmt.Sprintf("设备名:%s\n单台今收:%s\n单台总收:%s\n单台剩余:%s\n\n", v.Name, cast.ToString(v.TodayIncome), cast.ToString(v.AllIncome), cast.ToString(v.RemainIncome))
					continue
				}
				//最后一台设备最后一个回车
				content += fmt.Sprintf("设备名:%s\n单台今收：%s\n单台总收:%s\n单台剩余:%s\n", v.Name, cast.ToString(v.TodayIncome), cast.ToString(v.AllIncome), cast.ToString(v.RemainIncome))
			}
		}
		wechat.Push2Wechat(user, content)
	})
}
