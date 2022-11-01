package main

import (
	"flag"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/eventhook"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/simulate"
	"os"
	"os/signal"

	// 所有廣播訂閱平台
	_ "github.com/eric2788/MiraiValBot/sites/bilibili"
	_ "github.com/eric2788/MiraiValBot/sites/twitter"
	_ "github.com/eric2788/MiraiValBot/sites/valorant"
	_ "github.com/eric2788/MiraiValBot/sites/youtube"

	// 所有 redis 訂閱處理器
	_ "github.com/eric2788/MiraiValBot/handlers"

	// 所有指令
	_ "github.com/eric2788/MiraiValBot/cmd"

	// 所有定時器任務
	_ "github.com/eric2788/MiraiValBot/timer_tasks"

	// 註冊模組
	_ "github.com/eric2788/MiraiValBot/modules/broadcaster"
	_ "github.com/eric2788/MiraiValBot/modules/command"
	_ "github.com/eric2788/MiraiValBot/modules/response"
	_ "github.com/eric2788/MiraiValBot/modules/timer"
	_ "github.com/eric2788/MiraiValBot/modules/verbose"

	// 注册其他事件挂鈎
	_ "github.com/eric2788/MiraiValBot/chat_reply"
)

func init() {
	utils.WriteLogToFS()
}

var cliDebug = flag.Bool("debug", false, "enable debug logging level")

func main() {

	flag.Parse()

	if *cliDebug {
		simulate.EnableDebug()
	}

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
	bot.SaveToken()

	// 刷新好友列表，群列表
	bot.RefreshList()

	qq.InitValGroupInfo(bot.Instance)
	eventhook.HookBotEvents()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()

}
