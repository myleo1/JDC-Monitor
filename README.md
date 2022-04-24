# 京东云无线宝监工

> JDC-Monitor  支持收益汇总信息定时推送、IP、CPU占用、内存使用率、上传下载速度等信息实时监控

## 简介

本项目使用Golang开发，用到时序数据库influxdb、可视化工具grafana对京东云无线宝矿机信息进行可视化，并支持每天定时推送积分信息至微信。

目前支持采集的数据有：IP地址、CPU占用、内存使用率、ROM信息、上传下载速度、在线时长。数据采集为1分钟1次，可自行在代码中调整，数据采集走的是京东云无线宝app的接口，建议采集频率不要太高。

## 更新日志

2022.04.24

1、新增Grafana模板(需要替换MACADDR-1、MACADDR-2、MACADDR-3为自己的京东云MAC地址)

2022.02.15

1、升级依赖库版本

2、修改项目名,添加logo

2021.07.23

1、修复了多用户多协程读写map导致程序panic的bug

2、新增获取路由器离线信息，并且离线之后不再采集数据的功能

3、新增自动重启功能，每日定时检测每台设备收益，当收益低于设定的阈值即重启路由器（使用该功能需要使用京东云无线宝作为一级路由器拨号）

2021.04.07

1、修复了一台路由器离线导致其他设备无法采集数据的bug

## 部署

#### 准备工作

1、部署企业微信推送，参考https://github.com/myleo1/wechat-work-pusher

> 如果不使用企业微信推送也可使用server酱等第三方推送服务，需自己改动/service/wechat/wechat.go中的推送代码

2、部署influxdb及grafana，建议使用docker部署，具体部署方法请使用搜索引擎搜索

3、部署好influx后建一个库，并记住库名，后面要用

#### config.json配置

1、修改project.dir为项目目录

2、填写influxdb相关信息

3、使用charles或其他抓包工具如HttpCanary(手机端) 抓取数据包https://gw.smart.jd.com/f/service/xxxxxxxxx  开头的，请求头(request header)中的tgt、pin参数

> 如果不会抓包，安卓手机有root权限的可以用ES文件管理等进入打开"/data/data/com.jdcloud.mt.smartrouter/shared_prefs/jdc_mt_secured_store.xml"文件，找到"loginpin和wskey"的值即可。

4、config中的user为接收推送信息的微信号，如使用其他推送方式请自行修改

5、getZXQC为是否获取坐享其成打卡天数信息，由于有些设备不是坐享其成机型，可自行关闭

6、collect为是否采集路由器信息，如果只需要积分推送，那么不需要部署influx和grafana，只需要部署好企业微信推送服务并且关闭collect

7、reboot为自动重启路由器的收益阈值，当收益低于这个阈值时，自动重启路由器，不填写或者填写0为关闭自动重启功能

8、wechat.api为微信推送的api地址，wechat.token请参考https://github.com/myleo1/wechat-work-pusher

## Grafana配置

TODO

## 效果展示

### grafana可视化展示

![image](https://user-images.githubusercontent.com/66349676/111759475-993a9980-88d8-11eb-874a-26500c1d5398.png)



### 积分推送展示

![image](https://user-images.githubusercontent.com/66349676/111759767-ecace780-88d8-11eb-9bfa-05df471bf51c.png)



## TODO

1、一键部署

2、开发可视化前端（云监工web应用）

3、修改抓包登陆，直接手机验证码登陆