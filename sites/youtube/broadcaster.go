package youtube

import (
	"github.com/Logiase/MiraiGo-Template/bot"
)

type LiveDataHandler func(*bot.Bot, *LiveInfo) error

func (m *messageHandler) HandleLiveData(bot *bot.Bot, data interface{}, handle interface{}) error {
	liveData, handler := data.(*LiveInfo), handle.(LiveDataHandler)
	return handler(bot, liveData)
}

func RegisterDataHandler(cmd string, handle LiveDataHandler) {
	MessageHandler.AddHandler(cmd, handle)
}
