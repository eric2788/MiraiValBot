package main

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/file"
	"os"
	"os/signal"
)

func init() {
	utils.WriteLogToFS()
}

func main() {

	file.GenerateConfig()
	file.GenerateDevice()

	config.Init()
	file.LoadApplicationYaml()
	file.LoadStorage()

	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.AndroidPhone)

	// 登录
	err := bot.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 刷新好友列表，群列表
	bot.RefreshList()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()

}
