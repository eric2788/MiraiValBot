package command

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"sync"
)

const Tag = "valbot.command"

type command struct {
}

var (
	instance = &command{}
	logger   = utils.GetModuleLogger(Tag)
)

func (c *command) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (c *command) Init() {
	//TODO implement me
	panic("implement me")
}

func (c *command) PostInit() {
	//TODO implement me
	panic("implement me")
}

func (c *command) Serve(bot *bot.Bot) {
	bot.OnGroupMessage(func(client *client.QQClient, message *message.GroupMessage) {

	})
}

func (c *command) Start(bot *bot.Bot) {
	//TODO implement me
	panic("implement me")
}

func (c *command) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	//TODO implement me
	panic("implement me")
}
