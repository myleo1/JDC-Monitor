package task

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/service/cronkit"
	"github.com/mizuki1412/go-core-kit/service/influxkit"
	"github.com/spf13/cast"
	"jd-mining-server/service"
	"jd-mining-server/service/model"
	"jd-mining-server/service/wechat"
	"time"
)

func UpdateRouterList(pin, tgt string) {
	//更新路由器列表
	cronkit.AddFunc("@every 10m", func() {
		model.RouterMap[pin] = service.ListRouter(pin, tgt)
	})
}

func CollectTask(pin, tgt string) {
	//采集路由器数据
	cronkit.AddFunc("@every 1m", func() {
		for _, v := range model.RouterMap[pin] {
			s := service.GetPCDNStatus(v.FeedId, pin, tgt)
			sql := fmt.Sprintf("%s ip=%s,online=%s,cpu=%s,mem=%s,upload=%s,download=%s,rom=%s %d", v.Mac, influxkit.Decorate(s.Ip), s.OnlineTime, s.Cpu, s.Mem, s.Upload, s.Download, influxkit.Decorate(s.Rom), time.Now().UnixNano())
			influxkit.WriteDefaultDB(sql)
			//log.Println(sql)
		}
	})
}

func PushPointTask(pin, tgt, user string, waitFree bool) {
	//获取路由器积分等信息、然后推送至企业微信
	cronkit.AddFunc("0 30 7 */1 * ?", func() {
		//tgt与wsKey相等
		service.GetPointsDetail(pin, tgt, waitFree)
		var content string
		content += "********京东云矿机********\n\n"
		content += fmt.Sprintf("今日总收：%s\n总收：%s\n总剩：%s\n\n", cast.ToString(model.TotalPointsMap[pin].TotalToday), cast.ToString(model.TotalPointsMap[pin].TotalIncome), cast.ToString(model.TotalPointsMap[pin].TotalRemain))
		l := len(model.PointsDetailMap[pin])
		for i, v := range model.PointsDetailMap[pin] {
			if l != i+1 {
				content += fmt.Sprintf("设备名:%s\n单台今收：%s\n单台总收：%s\n单台剩余：%s\n\n", v.Name, cast.ToString(v.TodayIncome), cast.ToString(v.AllIncome), cast.ToString(v.RemainIncome))
				continue
			}
			//最后一台设备最后一个回车
			content += fmt.Sprintf("设备名:%s\n单台今收：%s\n单台总收：%s\n单台剩余：%s\n", v.Name, cast.ToString(v.TodayIncome), cast.ToString(v.AllIncome), cast.ToString(v.RemainIncome))
		}
		wechat.Push2Wechat(user, content)
	})
}
