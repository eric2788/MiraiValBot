package youtube

import (
	"github.com/Logiase/MiraiGo-Template/bot"
)

type LiveDataHandler func(*bot.Bot, *LiveInfo) error

var lastStatusMap = make(map[string]string)

func (m *messageHandler) HandleLiveData(bot *bot.Bot, data interface{}, handle interface{}) error {
	liveData, handler := data.(*LiveInfo), handle.(LiveDataHandler)
	status, ok := lastStatusMap[liveData.ChannelId]

	if !ok {
		status = Idle
	}

	if status == liveData.Status {
		logger.Infof("%s 的油管狀態與上一次相同，已略過。", liveData.ChannelId)
		return nil
	}

	lastStatusMap[liveData.ChannelId] = liveData.Status
	return handler(bot, liveData)
}

func RegisterDataHandler(cmd string, handle LiveDataHandler) {
	MessageHandler.AddHandler(cmd, handle)
}
