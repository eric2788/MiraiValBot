package valorant

import (
	"github.com/Logiase/MiraiGo-Template/bot"
)

type LiveDataHandler = func(*bot.Bot, *MatchMetaDataSub) error

func (m *messageHandler) HandleLiveData(bot *bot.Bot, data interface{}, handle interface{}) error {
	liveData, handler := data.(*MatchMetaDataSub), handle.(LiveDataHandler)
	return handler(bot, liveData)
}

func RegisterDataHandler(command string, handler LiveDataHandler) {
	MessageHandler.AddHandler(command, handler)
}
