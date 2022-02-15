package cmd

import (
	"JDC-Monitor/service"
	"JDC-Monitor/service/config"
	"JDC-Monitor/service/task"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/service/cronkit"
	"github.com/spf13/cobra"
)

func init() {
	DefFlags(rootCmd)
}

var rootCmd = &cobra.Command{
	Use: "JDC-Monitor",
	Run: func(cmd *cobra.Command, args []string) {
		initkit.BindFlags(cmd)
		//init
		config.Init()
		//多用户支持,一个用户单独一个协程
		for _, c := range config.Conf {
			go run(c)
		}
		select {}
	},
}

func run(c *config.Config) {
	service.InitRouterList(c.Pin, c.Tgt)
	//task
	taskStart(c)
}

func taskStart(c *config.Config) {
	if c.Collect {
		task.CollectTask(c.Pin, c.Tgt)
	}
	if c.Reboot != config.DoNOTReboot {
		task.RebootTask(c.Pin, c.Tgt, c.User, c.WaitFree)
	}
	task.UpdateRouterList(c.Pin, c.Tgt)
	task.PushPointTask(c.Pin, c.Tgt, c.User, c.WaitFree)
	cronkit.Scheduler().Start()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

func DefFlags(cmd *cobra.Command) {
	cmd.Flags().String(config.WechatApi, "", "微信推送api")
	cmd.Flags().String(config.WechatToken, "", "微信推送服务token")
}
