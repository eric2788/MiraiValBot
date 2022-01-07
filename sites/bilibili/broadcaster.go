package bilibili

import (
	"github.com/Logiase/MiraiGo-Template/bot"
)

type LiveDataHandler = func(*bot.Bot, *LiveData) error

func (h *messageHandler) HandleLiveData(bot *bot.Bot, data interface{}, handle interface{}) error {
	handler, liveData := handle.(LiveDataHandler), data.(*LiveData)
	return handler(bot, liveData)
}

func RegisterDataHandler(command string, handle LiveDataHandler) {
	MessageHandler.AddHandler(command, handle)
}
