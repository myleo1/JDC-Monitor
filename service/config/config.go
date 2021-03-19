package config

import (
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
)

type Config struct {
	Pin string
	Tgt string
	//是否获取坐享其成天数
	WaitFree bool
	//是否采集路由器数据
	Collect bool
	User    string
}

var Conf []*Config

func Init() {
	configReader := configkit.Get("jd", "")
	if configReader == "" {
		logkit.Fatal("read jd config err")
	}
	configList := cast.ToSlice(configReader)
	if len(configList) == 0 {
		logkit.Fatal("jd config length 0")
	}
	for _, v := range configList {
		m := cast.ToStringMap(v)
		c := &Config{}
		if val := cast.ToString(m["user"]); val == "" {
			logkit.Fatal("read user err")
		}
		c.User = cast.ToString(m["user"])
		if val := cast.ToString(m["pin"]); val == "" {
			logkit.Fatal("read pin err")
		}
		c.Pin = cast.ToString(m["pin"])
		if val := cast.ToString(m["tgt"]); val == "" {
			logkit.Fatal("read tgt err")
		}
		c.Tgt = cast.ToString(m["tgt"])
		c.WaitFree = cast.ToBool(m["getZXQC"])
		c.Collect = cast.ToBool(m["collect"])
		Conf = append(Conf, c)
	}
}
