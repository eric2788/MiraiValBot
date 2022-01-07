package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

var logger = utils.GetModuleLogger("sites.bilibili")

type messageHandler struct {
}

func (h *messageHandler) PubSubPrefix() string {
	return "blive:"
}

func (h *messageHandler) ToLiveData(message *redis.Message) (interface{}, error) {
	var liveData = &LiveData{}
	err := json.Unmarshal([]byte(message.Payload), liveData)
	return liveData, err
}

func (h *messageHandler) GetCommand(data interface{}) string {
	return data.(*LiveData).Command
}

func (h *messageHandler) GetOfflineListening() []string {
	listening := file.DataStorage.Listening.Bilibili
	topics := make([]string, len(listening))
	for i, v := range listening {
		topics[i] = fmt.Sprintf("blive:%d", v)
	}
	return topics
}

func (h *messageHandler) HandleError(bot *bot.Bot, error error) {
}

var MessageHandler = broadcaster.BuildHandle(logger, &messageHandler{})

func init() {
	broadcaster.RegisterHandler("bilibili", MessageHandler)
}
