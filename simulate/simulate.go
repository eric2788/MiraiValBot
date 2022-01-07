package simulate

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/sirupsen/logrus"
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
