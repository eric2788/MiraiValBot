package eventhook

import (
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/sirupsen/logrus"
)

type eventHookModule struct {
	EventHooker
	tag    string
	name   string
	logger logrus.FieldLogger
}

func (e *eventHookModule) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.ModuleID(e.tag),
		Instance: e,
	}
}

func (e *eventHookModule) Init() {
}

func (e *eventHookModule) PostInit() {
}

func (e *eventHookModule) Serve(bot *bot.Bot) {
}

func (e *eventHookModule) Start(bot *bot.Bot) {
	e.logger.Infof("%s 模組已啓動。", e.name)
}

func (e *eventHookModule) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	e.logger.Infof("%s 模組已關閉。", e.name)
}

func RegisterAsModule(hooker EventHooker, name, tag string, logger logrus.FieldLogger) {
	module := &eventHookModule{
		EventHooker: hooker,
		tag:         tag,
		name:        name,
		logger:      logger,
	}
	bot.RegisterModule(module)
	HookLifeCycle(module)
}
