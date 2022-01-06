package twitter

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

type messageHandler struct {
}

func (m *messageHandler) GetOfflineListening() []string {
	//TODO implement me
	panic("implement me")
}

func (m *messageHandler) HandleMessage(bot *bot.Bot, message *redis.Message) {
	//TODO implement me
	panic("implement me")
}

func (m *messageHandler) HandleError(bot *bot.Bot, error error) {
}

var MessageHandler = &messageHandler{}

func init() {
	broadcaster.RegisterHandler("twitter", MessageHandler)
}
