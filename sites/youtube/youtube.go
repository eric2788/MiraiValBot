package youtube

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

type messageHandler struct {
}

func (m *messageHandler) GetOfflineListening() []string {
	listening := file.DataStorage.Listening.Youtube
	topics := make([]string, len(listening))
	for i, v := range listening {
		topics[i] = fmt.Sprintf("ylive:%s", v)
	}
	return topics
}

func (m *messageHandler) HandleMessage(bot *bot.Bot, message *redis.Message) {
	//TODO implement me
	panic("implement me")
}

func (m *messageHandler) HandleError(bot *bot.Bot, error error) {
}

var MessageHandler = &messageHandler{}

func init() {
	broadcaster.RegisterHandler("youtube", MessageHandler)
}
