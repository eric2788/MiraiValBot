package simulate

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func EnableDebug() {
	logrus.SetLevel(logrus.DebugLevel) // for test
}

func RunBasic() {

	file.GenerateConfig()
	file.GenerateDevice()

	config.Init()
	file.LoadApplicationYaml()
	file.LoadStorage()

	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()
}

func RunLogin() {

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

	qq.InitValGroupInfo(bot.Instance)

}

func SignalForStop() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()

}
