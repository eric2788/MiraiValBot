package twitter

import "github.com/Logiase/MiraiGo-Template/bot"

type LiveDataHandler func(*bot.Bot, *TweetStreamData) error

func (m *messageHandler) HandleLiveData(bot *bot.Bot, data interface{}, handle interface{}) error {
	liveData, handler := data.(*TweetStreamData), handle.(LiveDataHandler)
	return handler(bot, liveData)
}

func RegisterDataHandler(cmd string, handler LiveDataHandler) {
	MessageHandler.AddHandler(cmd, handler)
}
