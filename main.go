package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/simulate"

	// 所有廣播訂閱平台
	_ "github.com/eric2788/MiraiValBot/hooks/sites/bilibili"
	_ "github.com/eric2788/MiraiValBot/hooks/sites/twitter"
	_ "github.com/eric2788/MiraiValBot/hooks/sites/valorant"
	_ "github.com/eric2788/MiraiValBot/hooks/sites/youtube"

	// 所有 redis 訂閱處理器
	_ "github.com/eric2788/MiraiValBot/hooks/handlers"

	// 所有指令
	_ "github.com/eric2788/MiraiValBot/hooks/cmd"

	// 所有 Discord 指令
	_ "github.com/eric2788/MiraiValBot/hooks/discord_cmd"

	// 所有定時器任務
	_ "github.com/eric2788/MiraiValBot/hooks/timer_tasks"

	// 所有游戏
	_ "github.com/eric2788/MiraiValBot/hooks/games"

	// 所有回应
	_ "github.com/eric2788/MiraiValBot/hooks/responses"

	// 註冊模組
	_ "github.com/eric2788/MiraiValBot/modules/broadcaster"
	_ "github.com/eric2788/MiraiValBot/modules/chat_reply"
	_ "github.com/eric2788/MiraiValBot/modules/command"
	_ "github.com/eric2788/MiraiValBot/modules/counting"
	_ "github.com/eric2788/MiraiValBot/modules/game"
	_ "github.com/eric2788/MiraiValBot/modules/privatechat"
	_ "github.com/eric2788/MiraiValBot/modules/repeatchat"
	_ "github.com/eric2788/MiraiValBot/modules/response"
	_ "github.com/eric2788/MiraiValBot/modules/timer"
	_ "github.com/eric2788/MiraiValBot/modules/urlparser"
	_ "github.com/eric2788/MiraiValBot/modules/verbose"
)

func init() {
	utils.WriteLogToFS()
}

var cliDebug = flag.Bool("debug", os.Getenv("DEBUG") == "true", "enable debug logging level")

func main() {

	flag.Parse()

	if *cliDebug {
		simulate.EnableDebug()
	}

	go debugServe()

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
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	bot.Stop()
	file.SaveStorage()
}
