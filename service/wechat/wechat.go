package wechat

import (
	"JDC-Monitor/service/config"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"time"
)

func Push2Wechat(to, msg string) {
	client := resty.New().SetRetryCount(5).SetRetryWaitTime(1 * time.Second)
	_, err := client.R().
		SetHeaders(map[string]string{
			"Cookie": fmt.Sprintf("session=%s", configkit.GetStringD(config.WechatToken)),
		}).SetFormData(map[string]string{
		"to":      to,
		"content": msg,
	}).Post(configkit.GetStringD(config.WechatApi))
	if err != nil {
		panic(exception.New(err.Error()))
	}
}
